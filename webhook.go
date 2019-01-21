// Based on https://github.com/morvencao/kube-mutating-webhook-tutorial/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
func (whsvr *WebhookServer) getEnvVarsToInject(pod *corev1.Pod, container *corev1.Container) []corev1.EnvVar {
	vars := []corev1.EnvVar{
		createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME", whsvr.clusterName),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_NODE_NAME", "spec.nodeName"),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_NAMESPACE_NAME", "metadata.namespace"),
		createEnvVarFromFieldPath("NEW_RELIC_METADATA_KUBERNETES_POD_NAME", "metadata.name"),
		createEnvVarFromString("NEW_RELIC_METADATA_KUBERNETES_CONTAINER_NAME", container.Name),
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

// WebhookServer is a webhook server that can accept requests from the Apiserver
type WebhookServer struct {
	clusterName string
	logger      *zap.SugaredLogger
	server      *http.Server
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
	// skip special kubernete system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			return false
		}
	}
	return true
}

func (whsvr *WebhookServer) updateContainer(pod *corev1.Pod, index int, container *corev1.Container) (patch []patchOperation) {
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
func (whsvr *WebhookServer) createPatch(pod *corev1.Pod) ([]byte, error) {
	var patch []patchOperation

	for i, container := range pod.Spec.Containers {
		patch = append(patch, whsvr.updateContainer(pod, i, &container)...)
	}

	return json.Marshal(patch)
}

// main mutation process
func (whsvr *WebhookServer) mutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		whsvr.logger.Errorw("could not unmarshal raw object", "err", err, "object", string(req.Object.Raw))
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	whsvr.logger.Infow("received admission review", "kind", req.Kind, "namespace", req.Namespace, "name",
		req.Name, "pod", pod.Name, "UID", req.UID, "operation", req.Operation, "userinfo", req.UserInfo)

	// determine whether to perform mutation
	if !mutationRequired(ignoredNamespaces, &pod.ObjectMeta) {
		// whsvr.logger.Infow("Skipping mutation for %s/%s due to policy check", pod.Namespace, pod.Name)
		whsvr.logger.Infow("skipped mutation", "namespace", pod.Namespace, "pod", pod.Name, "reason", "policy check (special namespaces)")
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := whsvr.createPatch(&pod)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	whsvr.logger.Infow("admission response created", "response", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Serve method for webhook server
func (whsvr *WebhookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte

	if whsvr.logger == nil {
		whsvr.logger = zap.NewNop().Sugar()
	}

	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		whsvr.logger.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		whsvr.logger.Errorw("invalid content type", "expected", "application/json", "context type", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		whsvr.logger.Errorw("can't decode body", "err", err, "body", body)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = whsvr.mutate(&ar)
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		whsvr.logger.Errorw("can't decode response", "err", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	whsvr.logger.Info("writing reponse")
	if _, err := w.Write(resp); err != nil {
		whsvr.logger.Errorw("can't write response", "err", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
