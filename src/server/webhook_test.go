package server

import (
	"bytes"
	json_encoding "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/api/admission/v1beta1"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

func TestServeHTTP(t *testing.T) {
	patchForValidBody, err := ioutil.ReadFile("testdata/expectedAdmissionReviewPatch.json")
	if err != nil {
		t.Fatalf("cannot read testdata file: %v", err)
	}
	var expectedPatchForValidBody bytes.Buffer
	if len(patchForValidBody) > 0 {
		if err := json_encoding.Compact(&expectedPatchForValidBody, patchForValidBody); err != nil {
			t.Fatalf(err.Error())
		}
	}

	missingObjectRequestBody := bytes.Replace(makeTestData(t, "default"), []byte("\"object\""), []byte("\"foo\""), -1)

	patchTypeForValidBody := v1beta1.PatchTypeJSONPatch
	cases := []struct {
		name                      string
		requestBody               []byte
		contentType               string
		expectedStatusCode        int
		expectedBodyWhenHTTPError string
		expectedAdmissionReview   v1beta1.AdmissionReview
	}{
		{
			name:               "mutation applied - valid body",
			requestBody:        makeTestData(t, "default"),
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectedAdmissionReview: v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					UID:       types.UID(1),
					Allowed:   true,
					Result:    nil,
					Patch:     expectedPatchForValidBody.Bytes(),
					PatchType: &patchTypeForValidBody,
				},
			},
		},
		{
			name:               "mutation not applied - valid body for ignored namespaces",
			requestBody:        makeTestData(t, "kube-system"),
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectedAdmissionReview: v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					UID:       types.UID(1),
					Allowed:   true,
					Result:    nil,
					Patch:     nil,
					PatchType: nil,
				},
			},
		},
		{
			name:                      "empty body",
			contentType:               "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedBodyWhenHTTPError: "empty body" + "\n",
		},
		{
			name:                      "wrong content-type",
			requestBody:               makeTestData(t, "default"),
			contentType:               "application/yaml",
			expectedStatusCode:        http.StatusUnsupportedMediaType,
			expectedBodyWhenHTTPError: "invalid Content-Type, expect `application/json`" + "\n",
		},
		{
			name:                      "invalid body",
			requestBody:               []byte{0, 1, 2},
			contentType:               "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedBodyWhenHTTPError: "could not decode request body: \"yaml: control characters are not allowed\"\n",
		},
		{
			name:                      "mutation fails - object not present in request body",
			requestBody:               missingObjectRequestBody,
			contentType:               "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedBodyWhenHTTPError: fmt.Sprintf("object not present in request body: %q\n", missingObjectRequestBody),
		},
	}

	whsvr := &Webhook{
		ClusterName: "foobar",
		Server:      &http.Server{},
	}

	server := httptest.NewServer(whsvr)
	defer server.Close()

	for i, c := range cases {
		t.Run(fmt.Sprintf("[%d] %s", i, c.name), func(t *testing.T) {

			resp, err := http.Post(server.URL, c.contentType, bytes.NewReader(c.requestBody))
			assert.NoError(t, err)
			assert.Equal(t, c.expectedStatusCode, resp.StatusCode)

			gotBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read body: %v", err)
			}
			var gotReview v1beta1.AdmissionReview
			if err := json.Unmarshal(gotBody, &gotReview); err != nil {
				assert.Equal(t, c.expectedBodyWhenHTTPError, string(gotBody))
				return
			}

			assert.Equal(t, c.expectedAdmissionReview, gotReview)
		})
	}

}

func Benchmark_WebhookPerformance(b *testing.B) {
	body := makeTestData(b, "default")

	whsvr := &Webhook{
		ClusterName: "foobar",
		Server: &http.Server{
			Addr: ":8080",
		},
	}

	server := httptest.NewServer(whsvr)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		http.Post(server.URL, "application/json", bytes.NewReader(body)) //nolint: errcheck
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
			Containers:       []corev1.Container{{Name: "c1", Image: "newrelic/image:latest"}, {Name: "c2", Image: "newrelic/image2:1.0.0"}},
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
