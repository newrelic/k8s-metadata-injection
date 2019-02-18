package server

import "net/http"

// TLSReadyReadinessProbe defines a readiness check for a Webhook struct based on the presence of its TLS certificate and key.
// It requires the whole webhook as parameter to be able to RLock on the certificate for the presence confirmation.
func TLSReadyReadinessProbe(webhook *Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webhook.RLock()
		defer webhook.RUnlock()

		if webhook.Cert == nil {
			response := "Certificate not present"
			w.WriteHeader(503)
			if _, err := w.Write([]byte(response)); err != nil {
				webhook.Logger.Errorw("can't write response", "err", err, "response", response)
			}
			return
		}

		okResponse := "OK"
		if _, err := w.Write([]byte(okResponse)); err != nil {
			webhook.Logger.Errorw("can't write response", "err", err, "response", okResponse)
		}
	}
}
