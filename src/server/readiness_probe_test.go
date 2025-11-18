package server

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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

// failingResponseWriter is a mock ResponseWriter that fails on Write().
type failingResponseWriter struct {
	statusCode int
	header     http.Header
}

func (f *failingResponseWriter) Header() http.Header {
	if f.header == nil {
		f.header = http.Header{}
	}
	return f.header
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("mock write error") //nolint:err113
}

func (f *failingResponseWriter) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}

func TestTLSReadyReadinessProbe_WriteErrorWithoutCert(t *testing.T) {
	t.Parallel()
	// Create an observer to capture log entries
	observedZapCore, observedLogs := observer.New(zap.ErrorLevel)
	observedLogger := zap.New(observedZapCore).Sugar()

	webhook := &Webhook{
		Cert:   nil,
		Logger: observedLogger,
	}

	handler := TLSReadyReadinessProbe(webhook)
	req := httptest.NewRequest("GET", "/ready", nil)
	w := &failingResponseWriter{}

	handler.ServeHTTP(w, req)

	// Verify that WriteHeader was called with 503
	assert.Equal(t, 503, w.statusCode, "Should return 503 when certificate is not present")

	// Verify that the error was logged
	logEntries := observedLogs.All()
	assert.Equal(t, 1, len(logEntries), "Should have logged one error")
	assert.Equal(t, "can't write response", logEntries[0].Message)
	assert.Equal(t, "Certificate not present", logEntries[0].ContextMap()["response"])
	assert.Contains(t, logEntries[0].ContextMap()["err"], "mock write error")
}

func TestTLSReadyReadinessProbe_WriteErrorWithCert(t *testing.T) {
	t.Parallel()
	// Create an observer to capture log entries
	observedZapCore, observedLogs := observer.New(zap.ErrorLevel)
	observedLogger := zap.New(observedZapCore).Sugar()

	webhook := &Webhook{
		Cert:   &tls.Certificate{},
		Logger: observedLogger,
	}

	handler := TLSReadyReadinessProbe(webhook)
	req := httptest.NewRequest("GET", "/ready", nil)
	w := &failingResponseWriter{}

	handler.ServeHTTP(w, req)

	// Verify that WriteHeader was not called (status 200 is implicit)
	// and that the handler completed without panicking despite Write() failing
	assert.Equal(t, 0, w.statusCode, "Should not explicitly set status code when cert is present (defaults to 200)")

	// Verify that the error was logged
	logEntries := observedLogs.All()
	assert.Equal(t, 1, len(logEntries), "Should have logged one error")
	assert.Equal(t, "can't write response", logEntries[0].Message)
	assert.Equal(t, "OK", logEntries[0].ContextMap()["response"])
	assert.Contains(t, logEntries[0].ContextMap()["err"], "mock write error")
}
