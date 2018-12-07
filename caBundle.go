package main

import (
	"encoding/json"
	"fmt"
	b64 "encoding/base64"

	"github.com/golang/glog"

	"k8s.io/api/admissionregistration/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	admissionregistrationv1beta1client "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
)

func UpdateCaBundle(webhookConfigName, webhookName, caBundle string) (error) {
	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// Create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Decode the base64 encoded caBundle
	caDecoded, err := b64.StdEncoding.DecodeString(caBundle)
	if err != nil {
		return err
	}

	// Patch the MutatingWebhookConfig
	return PatchMutatingWebhookConfig(client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations(),
		                              webhookConfigName, webhookName, caDecoded)
}

// PatchMutatingWebhookConfig patches a CA bundle into the specified webhook config.
// Taken from istio: https://github.com/istio/istio/blob/master/pkg/util/webhookpatch.go
func PatchMutatingWebhookConfig(client admissionregistrationv1beta1client.MutatingWebhookConfigurationInterface,
	webhookConfigName, webhookName string, caBundle []byte) error {

	config, err := client.Get(webhookConfigName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("Error on client get")
		return err
	}
	prev, err := json.Marshal(config)
	if err != nil {
		return err
	}
	found := false
	for i, w := range config.Webhooks {
		if w.Name == webhookName {
			config.Webhooks[i].ClientConfig.CABundle = caBundle[:]
			found = true
			break
		}
	}
	if !found {
		return apierrors.NewInternalError(fmt.Errorf(
			"webhook entry %q not found in config %q", webhookName, webhookConfigName))
	}
	curr, err := json.Marshal(config)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(prev, curr, v1beta1.MutatingWebhookConfiguration{})
	if err != nil {
		return err
	}

	if string(patch) != "{}" {
		glog.Infof("Performing WebhookConfiguration patch: %s", string(patch))
		_, err = client.Patch(webhookConfigName, types.StrategicMergePatchType, patch)
	}
	return err
}
