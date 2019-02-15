package server

import "net/http"

// HealthCheck defines a readiness check for a Webhook struct based on the presence of its TLS certificate and key.
type HealthCheck struct {
	webhookServer *Webhook
}

// ServeHTTP checks if the Webhook server has a TLS certificate/key pair.
func (h *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.webhookServer.RLockCert()
	defer h.webhookServer.RUnlockCert()

	if h.webhookServer.Cert == nil {
		w.WriteHeader(503)
		var response = "Certificate not present."
		w.WriteHeader(503)
		if _, err := w.Write([]byte(response)); err != nil {
			h.webhookServer.Logger.Errorw("can't write response", "err", err, "response", response)
		}
		return
	}

	var okResponse = "OK"
	w.WriteHeader(200)
	if _, err := w.Write([]byte(okResponse)); err != nil {
		h.webhookServer.Logger.Errorw("can't write response", "err", err, "response", okResponse)
	}
}

// NewHealthCheck is a constructor for HealthCheck.
func NewHealthCheck(w *Webhook) *HealthCheck {
	return &HealthCheck{
		webhookServer: w,
	}
}
