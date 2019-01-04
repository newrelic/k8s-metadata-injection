package main

import (
	"bytes"
	json_encoding "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

func TestServeHTTP(t *testing.T) {
	patch, err := ioutil.ReadFile("testdata/expectedAdmissionReviewPatch.json")
	if err != nil {
		t.Fatalf("cannot read testdata file: %v", err)
	}

	cases := []struct {
		name               string
		body               []byte
		contentType        string
		expectedAllowed    bool
		expectedStatusCode int
		expectedPatch      []byte
		expectedUID        types.UID
		expectedHTTPError  string
	}{
		{
			name:               "mutation applied - valid body",
			body:               makeTestData(t, "default"),
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectedAllowed:    true,
			expectedPatch:      patch,
			expectedUID:        types.UID(1),
		},
		{
			name:               "mutation not applied - valid body for ignored namespaces",
			body:               makeTestData(t, "kube-system"),
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectedAllowed:    true,
			expectedPatch:      nil,
			expectedUID:        types.UID(1),
		},
		{
			name:               "empty body",
			contentType:        "application/json",
			expectedStatusCode: http.StatusBadRequest,
			expectedHTTPError:  fmt.Sprintln("empty body"),
			expectedAllowed:    false,
			expectedPatch:      nil,
			expectedUID:        "",
		},
		{
			name:               "wrong content-type",
			body:               makeTestData(t, "default"),
			contentType:        "application/yaml",
			expectedStatusCode: http.StatusUnsupportedMediaType,
			expectedHTTPError:  fmt.Sprintln("invalid Content-Type, expect `application/json`"),
			expectedAllowed:    false,
			expectedPatch:      nil,
			expectedUID:        "",
		},
		{
			name:               "invalid body",
			body:               []byte{0, 1, 2},
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectedHTTPError:  "",
			expectedAllowed:    false,
			expectedPatch:      nil,
			expectedUID:        "",
		},
	}

	whsvr := &WebhookServer{
		clusterName: "foobar",
		server: &http.Server{
			Addr: ":8080",
		},
	}

	server := httptest.NewServer(whsvr)
	defer server.Close()

	for i, c := range cases {
		t.Run(fmt.Sprintf("[%d] %s", i, c.name), func(t *testing.T) {

			resp, err := http.Post(server.URL, c.contentType, bytes.NewReader(c.body))
			assert.NoError(t, err)
			assert.Equal(t, c.expectedStatusCode, resp.StatusCode)

			gotBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read body: %v", err)
			}
			var gotReview v1beta1.AdmissionReview
			if err := json.Unmarshal(gotBody, &gotReview); err != nil {
				assert.Equal(t, c.expectedHTTPError, string(gotBody))
				return
			}
			assert.Equal(t, c.expectedAllowed, gotReview.Response.Allowed)

			var gotPatch bytes.Buffer
			if len(gotReview.Response.Patch) > 0 {
				if err := json_encoding.Compact(&gotPatch, gotReview.Response.Patch); err != nil {
					t.Fatalf(err.Error())
				}
			}
			var expectedPatch bytes.Buffer
			if len(c.expectedPatch) > 0 {
				if err := json_encoding.Compact(&expectedPatch, c.expectedPatch); err != nil {
					t.Fatalf(err.Error())
				}
			}
			assert.Equal(t, expectedPatch.Bytes(), gotPatch.Bytes())
		})
	}

}

func Benchmark_WebhookPerformance(b *testing.B) {
	body := makeTestData(b, "default")

	whsvr := &WebhookServer{
		clusterName: "foobar",
		server: &http.Server{
			Addr: ":8080",
		},
	}

	server := httptest.NewServer(whsvr)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		http.Post(server.URL, "application/json", bytes.NewReader(body))
	}
}

func makeTestData(t testing.TB, namespace string) []byte {
	t.Helper()

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "test-123-123",
			GenerateName:    "test-123-123", // required for creating metadata for deployment
			Annotations:     map[string]string{},
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{{Kind: "ReplicaSet"}}, // required for populating metadata for deployment
		},
		Spec: corev1.PodSpec{
			Volumes:          []corev1.Volume{{Name: "v0"}},
			InitContainers:   []corev1.Container{{Name: "c0"}},
			Containers:       []corev1.Container{{Name: "c1"}, {Name: "c2"}},
			ImagePullSecrets: []corev1.LocalObjectReference{{Name: "p0"}},
		},
	}

	raw, err := json.Marshal(&pod)
	if err != nil {
		t.Fatalf("Could not create test pod: %v", err)
	}

	review := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{},
			Object: runtime.RawExtension{
				Raw: raw,
			},
			Operation: v1beta1.Create,
			UID:       types.UID(1),
		},
	}
	reviewJSON, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("Failed to create AdmissionReview: %v", err)
	}
	return reviewJSON
}
