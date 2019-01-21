package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// WhSvrParameters are configuration parameters for Webhook Server
type WhSvrParameters struct {
	port        int    // webhook server port
	certFile    string // path to the x509 certificate for https
	keyFile     string // path to the x509 private key matching `CertFile`
	clusterName string // name of the cluster
}

func main() {
	var parameters WhSvrParameters

	flag.IntVar(&parameters.port, "port", 443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/tls-key-cert-pair/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/tls-key-cert-pair/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&parameters.clusterName, "clusterName", "cluster", "The name of the Kubernetes cluster")
	flag.Parse()

	logger := setupLogger()
	defer func() { _ = logger.Sync() }()

	pair, err := tls.LoadX509KeyPair(parameters.certFile, parameters.keyFile)
	if err != nil {
		logger.Errorw("failed to load key pair", "err", err)
	}

	watcher, _ := fsnotify.NewWatcher()
	defer func() { _ = watcher.Close() }()
	// Watch the parent directory of the target files so we can catch
	// symlink updates of k8s secrets volumes and reload the certificates whenever they change.
	for _, file := range []string{parameters.certFile, parameters.keyFile} {
		watchDir, _ := filepath.Split(file)
		if err := watcher.Add(watchDir); err != nil {
			logger.Errorw("could not watch file", "file", file, "err", err)
		}
	}

	whsvr := &WebhookServer{
		keyFile:     parameters.keyFile,
		certFile:    parameters.certFile,
		cert:        &pair,
		clusterName: parameters.clusterName,
		certWatcher: watcher,
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", parameters.port),
		},
		logger: logger,
	}
	whsvr.server.TLSConfig = &tls.Config{GetCertificate: whsvr.getCert}

	mux := http.NewServeMux()
	mux.Handle("/mutate", whsvr)
	whsvr.server.Handler = mux

	go func() {
		logger.Info("starting the webhook server")
		if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
			logger.Errorw("failed to start webhook server", "err", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case event := <-whsvr.certWatcher.Events:
			// TODO: use a timer to debounce configuration updates
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				pair, err := tls.LoadX509KeyPair(whsvr.certFile, whsvr.keyFile)
				if err != nil {
					logger.Errorw("reload cert error", "err", err)
					break
				}
				whsvr.mu.Lock()
				whsvr.cert = &pair
				whsvr.mu.Unlock()
				logger.Info("Cert/key pair reloaded!")
			}
		case <-signalChan:
			logger.Info("Got OS shutdown signal, shutting down webhook server gracefully...")
			_ = watcher.Close()
			_ = whsvr.server.Shutdown(context.Background())
		}
	}
}

func setupLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	return zapLogger.Sugar()
}
