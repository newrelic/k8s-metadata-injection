# Kubernetes Metadata injection for New Relic APM agents

## How does it work ?

New Relic requires the following environment variables to identify Kubernetes objects in the APM agents:
- `K8S_CLUSTER_NAME`
- `K8S_NODE_NAME`
- `K8S_NAMESPACE_NAME`
- `K8S_DEPLOYMENT_NAME`
- `K8S_POD_NAME`
- `K8S_CONTAINER_NAME`

These environment variables can be set manually by the customer, or they can be automatically injected using a MutatingAdmissionWebhook.
New Relic provides an easy method for deploying this automatic approach.

## Automatic environment variable injection

### 1) Check if MutatingAdmissionWebhook is enabled on your cluster

This feature requires Kubernetes 1.9 or later. Verify that the kube-apiserver process has the admission-control flag set.

```
$ kubectl api-versions | grep admissionregistration.k8s.io/v1beta1
admissionregistration.k8s.io/v1beta1
```

### 2) Install the injection

```
$ kubectl apply -f newrelic-metadata-injection.yaml
```

Executing this
- creates `newrelic-metadata-injection-deployment` and `newrelic-metadata-injection-svc`,
- registers the `newrelic-metadata-injection-svc` service as a MutatingAdmissionWebhook with the Kubernetes api

### 3) Enable the automatic Kubernetes metadata injection on your namespaces

The injection is only applied to namespaces that have the `newrelic-metadata-injection` label set to `enabled`.

```
$ kubectl label namespace <namespace> newrelic-metadata-injection=enabled
```

## Prototype

This repo contains a prototype based on https://github.com/morvencao/kube-mutating-webhook-tutorial/.
*Important note:* this is just a prototype and not production ready!
We need readinesschecks, healthchecks and a lot of testing since this will run in the Pod deployment flow on the Kubernetes cluster.

### Build

The prototype uses dep as the dependency management tool:

```
go get -u github.com/golang/dep/cmd/dep
```

Build and push the docker image, this currently pushes to an AWS machine from fryckbosch (which is not in the DNS - so it will fail):

```
./build.sh
```

### Certificates

To make the deployment as easy as possible, the certificates for the webhook are generated inside the container.
The webhook container then uses the kubernetes api to update the caBundle on the MutatingAdmissionWebhook.

This is more-or-less how Istio does things. It does mean that we need to have a service account that can update MutatingAdmissionWebhook.
This also means that only 1 replica of the webhook service can be running.

Customers that don't want this automatic approach, can create certificates themselves:

```
./create-certs.sh
kubectl create -f newrelic-metadata-injection-manual-cert.yaml
```

The command above requires the following files from this repo: `create-certs.sh`, `newrelic-metadata-injection-manual-nocert.yaml`.
