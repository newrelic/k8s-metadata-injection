package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

const (
	appName        = "new-relic-k8s-metadata-injection"
	defaultTimeout = time.Second * 30
)

// Specification contains the specs for this app.
type Specification struct {
	Port        int           `default:"443"`                                                      // Webhook server port.
	TLSCertFile string        `default:"/etc/tls-key-cert-pair/tls.crt" envconfig:"tls_cert_file"` // File containing the x509 Certificate for HTTPS.
	TLSKeyFile  string        `default:"/etc/tls-key-cert-pair/tls.key" envconfig:"tls_key_file"`  // File containing the x509 private key for TLSCERTFILE.
	ClusterName string        `default:"cluster" split_words:"true"`                               // The name of the Kubernetes cluster.
	Timeout     time.Duration // server timeout. Defaults to the timeout passed by K8s API via query param. If not present, to the defaultTimeout const value.
}

func main() {
	var s Specification
	err := envconfig.Process(strings.Replace(appName, "-", "_", -1), &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	zapLogger, err := zap.NewDevelopment() // We want a human readable log that includes as much as detail as possible.
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger := zapLogger.Sugar()

	defer func() { _ = logger.Sync() }()

	pair, err := tls.LoadX509KeyPair(s.TLSCertFile, s.TLSKeyFile)
	if err != nil {
		logger.Errorw("failed to load key pair", "err", err)
	}

	watcher, _ := fsnotify.NewWatcher()
	defer func() { _ = watcher.Close() }()
	// Watch the parent directory of the key/cert files so we can catch
	// symlink updates of k8s secrets volumes and reload the certificates whenever they change.
	watchDir, _ := filepath.Split(s.TLSCertFile)
	if err := watcher.Add(watchDir); err != nil {
		logger.Errorw("could not watch folder", "folder", watchDir, "err", err)
	}

	whsvr := &WebhookServer{
		keyFile:     s.TLSKeyFile,
		certFile:    s.TLSCertFile,
		cert:        &pair,
		clusterName: s.ClusterName,
		certWatcher: watcher,
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", s.Port),
		},
		logger: logger,
	}
	whsvr.server.TLSConfig = &tls.Config{GetCertificate: whsvr.getCert}

	mux := http.NewServeMux()
	mux.Handle("/mutate", withLoggingMiddleware(logger)(withTimeoutMiddleware(s.Timeout)(whsvr)))
	whsvr.server.Handler = mux

	go func() {
		logger.Info("starting the webhook server")
		if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
			logger.Errorw("failed to start webhook server", "err", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var debounceTimer <-chan time.Time
	for {
		select {
		case <-debounceTimer:
			pair, err := tls.LoadX509KeyPair(whsvr.certFile, whsvr.keyFile)
			if err != nil {
				logger.Errorw("reload cert error", "err", err)
				break
			}
			whsvr.mu.Lock()
			whsvr.cert = &pair
			whsvr.mu.Unlock()
			logger.Info("cert/key pair reloaded!")
		case event := <-whsvr.certWatcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				debounceTimer = time.After(500 * time.Millisecond)
			}
		case <-signalChan:
			logger.Info("got OS shutdown signal, shutting down webhook server gracefully...")
			_ = watcher.Close()
			_ = whsvr.server.Shutdown(context.Background())
			return
		}
	}
}

func withTimeoutMiddleware(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In case the user does not set a timeout, we use the timeout passed by K8s API via query param.
			// If the latest timeout is not present in the form of URL query param, we use the defaultTimeout const value.
			if timeout.Nanoseconds() == 0 {
				if qt := r.URL.Query().Get("timeout"); qt != "" {
					timeout, _ = time.ParseDuration(qt)
				} else {
					timeout = defaultTimeout
				}
			}

			http.TimeoutHandler(next, timeout, "server timeout").ServeHTTP(w, r)
		})
	}
}

func withLoggingMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logger.Infof("%s %s://%s%s %s\" from %s", r.Method, scheme, r.Host, r.RequestURI, r.Proto, r.RemoteAddr)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
