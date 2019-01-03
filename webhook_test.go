package main

import (
	"bytes"
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Benchmark_WebhookPerformance(b *testing.B) {
	body := makeTestData(b, false)

	whsvr := &WebhookServer {
		clusterName: "foobar",
		server: &http.Server {
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

func makeTestData(t testing.TB, skip bool) []byte {
	t.Helper()

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test",
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			Volumes:          []corev1.Volume{{Name: "v0"}},
			InitContainers:   []corev1.Container{{Name: "c0"}},
			Containers:       []corev1.Container{{Name: "c1"}},
			ImagePullSecrets: []corev1.LocalObjectReference{{Name: "p0"}},
		},
	}

	if skip {
		fmt.Println("skipped!")
		return []byte{}
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
		},
	}
	reviewJSON, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("Failed to create AdmissionReview: %v", err)
	}
	return reviewJSON
}

