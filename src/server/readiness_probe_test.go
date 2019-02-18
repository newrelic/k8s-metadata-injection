package server

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLSReadyReadinessProbe(t *testing.T) {
	cases := []struct {
		desc         string
		certificate  *tls.Certificate
		responseCode int
	}{
		{
			desc:         "certificate not present (bad health)",
			certificate:  nil,
			responseCode: 503,
		},
		{
			desc:         "certificate present (good health)",
			certificate:  &tls.Certificate{},
			responseCode: 200,
		},
	}

	webhook := Webhook{}
	healthCheck := http.HandlerFunc(TLSReadyReadinessProbe(&webhook))
	server := httptest.NewServer(healthCheck)

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			webhook.Cert = c.certificate

			resp, err := http.Get(server.URL)

			assert.NoError(t, err)
			assert.Equal(t, c.responseCode, resp.StatusCode)
		})
	}
}
