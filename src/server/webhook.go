// Based on https://github.com/morvencao/kube-mutating-webhook-tutorial/

package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

var ignoredNamespaces = []string{
	metav1.NamespaceSystem,
	metav1.NamespacePublic,
}

func createEnvVarFromFieldPath(envVarName, fieldPath string) corev1.EnvVar {
	return corev1.EnvVar{Name: envVarName, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: fieldPath}}}
}

func createEnvVarFromString(envVarName, envVarValue string) corev1.EnvVar {
	return corev1.EnvVar{Name: envVarName, Value: envVarValue}
}

// getEnvVarsToInject returns the environment variables to inject in the given container
func (whsvr *Webhook) getEnvVarsToInject(pod *corev1.Pod, container *corev1.Container) []corev1.EnvVar {
	vars := []corev1.EnvVar{
		createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME", whsvr.ClusterName),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_NODE_NAME", "spec.nodeName"),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_NAMESPACE_NAME", "metadata.namespace"),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_POD_NAME", "metadata.name"),
		createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_CONTAINER_NAME", container.Name),
		createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_CONTAINER_IMAGE_NAME", container.Image),
	}

	// Guess the name of the deployment. We check whether the Pod is Owned by a ReplicaSet and confirms with the
	// naming convention for a Deployment. This can give a false positive if the user uses ReplicaSets directly.
	if len(pod.OwnerReferences) == 1 && pod.OwnerReferences[0].Kind == "ReplicaSet" {
		podParts := strings.Split(pod.GenerateName, "-")
		if len(podParts) >= 3 {
			deployment := strings.Join(podParts[:len(podParts)-2], "-")
			vars = append(vars, createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_DEPLOYMENT_NAME", deployment))
		}
	}

	return vars
}

// Webhook is a webhook server that can accept requests from the Apiserver
type Webhook struct {
	CertFile    string
	KeyFile     string
	Cert        *tls.Certificate
	ClusterName string
	Logger      *zap.SugaredLogger
	Mu          sync.RWMutex
	Server      *http.Server
	CertWatcher *fsnotify.Watcher
}

// GetCert returns the certificate that should be used by the server in the TLS handshake.
func (whsvr *Webhook) GetCert(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	whsvr.Mu.Lock()
	defer whsvr.Mu.Unlock()
	return whsvr.Cert, nil
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	_ = corev1.AddToScheme(runtimeScheme)
}

// Check whether the target resoured need to be mutated
func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			return false
		}
	}
	return true
}

func (whsvr *Webhook) updateContainer(pod *corev1.Pod, index int, container *corev1.Container) (patch []patchOperation) {
	// Create map with all environment variable names
	envVarMap := map[string]bool{}
	for _, envVar := range container.Env {
		envVarMap[envVar.Name] = true
	}

	// Create a patch for each EnvVar in toInject (if they are not yet defined on the container)
	first := len(envVarMap) == 0
	var value interface{}
	basePath := fmt.Sprintf("/spec/containers/%d/env", index)

	for _, inject := range whsvr.getEnvVarsToInject(pod, container) {
		if _, present := envVarMap[inject.Name]; !present {
			value = inject
			path := basePath

			if first {
				// For the first element we have to create the list
				value = []corev1.EnvVar{inject}
				first = false
			} else {
				// For the other elements we can append to the list
				path = path + "/-"
			}

			patch = append(patch, patchOperation{
				Op:    "add",
				Path:  path,
				Value: value,
			})
		}
	}
	return patch
}

// create mutation patch for resources
func (whsvr *Webhook) createPatch(pod *corev1.Pod) ([]byte, error) {
	var patch []patchOperation

	for i, container := range pod.Spec.Containers {
		patch = append(patch, whsvr.updateContainer(pod, i, &container)...)
	}

	return json.Marshal(patch)
}

// main mutation process
func (whsvr *Webhook) mutate(ar *v1beta1.AdmissionReview) ([]byte, error) {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		whsvr.Logger.Errorw("could not unmarshal raw object", "err", err, "object", string(req.Object.Raw))
		return nil, err
	}

	whsvr.Logger.Infow("received admission review", "kind", req.Kind, "namespace", req.Namespace, "name",
		req.Name, "pod", pod.Name, "UID", req.UID, "operation", req.Operation, "userinfo", req.UserInfo)

	// determine whether to perform mutation
	if !mutationRequired(ignoredNamespaces, &pod.ObjectMeta) {
		whsvr.Logger.Infow("skipped mutation", "namespace", pod.Namespace, "pod", pod.Name, "reason", "policy check (special namespaces)")
		return nil, nil
	}

	patchBytes, err := whsvr.createPatch(&pod)
	if err != nil {
		return nil, err
	}

	whsvr.Logger.Infow("admission response created", "response", string(patchBytes))
	return patchBytes, nil
}

// Serve method for webhook server
func (whsvr *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte

	if whsvr.Logger == nil {
		whsvr.Logger = zap.NewNop().Sugar()
	}

	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		whsvr.Logger.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		whsvr.Logger.Errorw("invalid content type", "expected", "application/json", "context type", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Allowed: true, // Always allow the creation of the pod since this webhook does not act as Validating Webhook.
		},
	}

	admissionReviewRequest := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &admissionReviewRequest); err != nil {
		whsvr.Logger.Errorw("can't decode body", "err", err, "body", body)
		http.Error(w, fmt.Sprintf("could not decode request body: %q", err.Error()), http.StatusBadRequest)
		return
	}

	if len(admissionReviewRequest.Request.Object.Raw) == 0 {
		whsvr.Logger.Errorw("object not present in request body", "body", body)
		http.Error(w, fmt.Sprintf("object not present in request body: %q", body), http.StatusBadRequest)
		return
	}

	patch, err := whsvr.mutate(&admissionReviewRequest)
	if err != nil {
		whsvr.Logger.Errorw("error during mutation", "err", err)
		http.Error(w, fmt.Sprintf("error during mutation: %q", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(patch) > 0 {
		admissionReviewResponse.Response.Patch = patch
		admissionReviewResponse.Response.PatchType = func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch // Only PatchTypeJSONPatch is allowed by now.
			return &pt
		}()
	}

	if admissionReviewRequest.Request != nil {
		admissionReviewResponse.Response.UID = admissionReviewRequest.Request.UID
	}

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		whsvr.Logger.Errorw("can't decode response", "err", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}
	whsvr.Logger.Info("writing response")
	if _, err := w.Write(resp); err != nil {
		whsvr.Logger.Errorw("can't write response", "err", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
		return
	}
}
