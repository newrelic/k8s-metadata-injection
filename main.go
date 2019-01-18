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

	"github.com/howeyc/fsnotify"
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

	// get command line parameters
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
	// watch the parent directory of the target files so we can catch
	// symlink updates of k8s ConfigMaps volumes.
	for _, file := range []string{parameters.certFile, parameters.keyFile} {
		watchDir, _ := filepath.Split(file)
		if err := watcher.Watch(watchDir); err != nil {
			glog.Errorf("could not watch %v: %v", file, err)
		}
	}
	defer func() { _ = watcher.Close() }()

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

	// define http server and server handler
	logger.Info("starting the webhook server")

	mux := http.NewServeMux()
	mux.Handle("/mutate", whsvr)
	whsvr.server.Handler = mux

	// start webhook server in new rountine
	go func() {
		if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
			logger.Errorw("failed to start webhook server", "err", err)
		}
	}()

	for {
		select {
		case event := <-whsvr.certWatcher.Event:
			// TODO: use a timer to debounce configuration updates
			if event.IsModify() || event.IsCreate() {
				pair, err := tls.LoadX509KeyPair(whsvr.certFile, whsvr.keyFile)
				if err != nil {
					glog.Errorf("reload cert error: %v", err)
					break
				}
				whsvr.mu.Lock()
				whsvr.cert = &pair
				whsvr.mu.Unlock()
				glog.Info("Cert/key pair reloaded!")
			}
		case <-signalChan:
			glog.Infof("Got OS shutdown signal, shutting down wenhook server gracefully...")
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
