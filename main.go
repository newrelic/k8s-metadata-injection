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
	"syscall"

	"go.uber.org/zap"
)

// WhSvrParameters are configuration parameters for Webhook Server
type WhSvrParameters struct {
	port              int    // webhook server port
	certFile          string // path to the x509 certificate for https
	keyFile           string // path to the x509 private key matching `CertFile`
	clusterName       string // name of the cluster
	webhookConfigName string // name of the webhook config
	webhookName       string // name of the webhook
	caBundle          string // caBundle
}

func main() {
	var parameters WhSvrParameters

	// get command line parameters
	flag.IntVar(&parameters.port, "port", 443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&parameters.clusterName, "clusterName", "cluster", "The name of the Kubernetes cluster")
	flag.StringVar(&parameters.webhookConfigName, "webhookConfigName", "newrelic-metadata-injection-cfg", "Optional name of the MutatingAdmissionWebhook to push webhook caBundle")
	flag.StringVar(&parameters.webhookName, "webhookName", "metadata-injection.newrelic.com", "Optional name of the webhook to push to webhook caBundle")
	flag.StringVar(&parameters.caBundle, "caBundle", "", "Optional caBundle to push to the Kubernetes API")
	flag.Parse()

	logger := setupLogger()
	defer func() { _ = logger.Sync() }()

	pair, err := tls.LoadX509KeyPair(parameters.certFile, parameters.keyFile)
	if err != nil {
		logger.Errorw("failed to load key pair", "err", err)
	}

	whsvr := &WebhookServer{
		clusterName: parameters.clusterName,
		server: &http.Server{
			Addr:      fmt.Sprintf(":%v", parameters.port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
		logger: logger,
	}

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

	// push the caBundle to the Kubernetes API if provided
	if parameters.caBundle != "" {
		go func() {
			if err := UpdateCaBundle(parameters.webhookConfigName, parameters.webhookName, parameters.caBundle, logger); err != nil {
				logger.Errorw("failed to update caBundle on the MutatingAdmissionWebhook", "name", parameters.webhookConfigName, "err", err)
			} else {
				logger.Infof("successfully updated caBundle on MutatingAdmissionWebhook %s", parameters.webhookConfigName)
			}
		}()
	}

	// listening OS shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("got OS shutdown signal, shutting down webhook server gracefully...")
	whsvr.server.Shutdown(context.Background()) //nolint: errcheck
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
