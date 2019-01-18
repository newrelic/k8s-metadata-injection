# Kubernetes Metadata injection for New Relic APM agents

## How does it work?

New Relic requires the following environment variables to identify Kubernetes objects in the APM agents:

- `NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NODE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NAMESPACE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_DEPLOYMENT_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_POD_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_CONTAINER_NAME`

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

### 4) Certificates

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

## Development

### Prerequisites

For the development process [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube) and [Skaffold](https://github.com/GoogleCloudPlatform/skaffold) tools are used.


* [Install Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/);
* [Install Skaffold](https://github.com/GoogleCloudPlatform/skaffold#installation).

Currently the project compiles with **Go 1.11.4**.

### Dependency management

[Go modules](https://github.com/golang/go/wiki/Modules) are used for managing dependencies. This project does not need to be in your GOROOT, if you wish so.

Currently for K8s libraries it uses version 1.13.1. Only couple of libraries are direct dependencies, the rest are indirect. You need to point all of them to the same K8s version to make sure that everything works as expected. For the moment this process is manual.

### Configuration

* Copy the deployment file `deploy/newrelic-metadata-injection.yaml` to `deploy/local.yaml`.
* Edit the file and set the following value as container image: `internal/k8s-metadata-injector`.
* Make sure that `imagePullPolicy: Always` is not present in the file (otherwise, the image won't be pulled).

### Run

Run `make deploy-dev`. This will compile your binary with compatibility for the container OS architecture, build a temporary docker image and finally deploy it to your Minikube.

If you would like to enable automatic redeploy on changes to the repository, you can run `skaffold dev`.

### Tests

For running unit tests, use

```bash
make test
```

For running benchmark tests, use:

```bash
make benchmark-test
```

### Documentation

Please use the [Open Api 3.0 spec file](openapi.yaml) as documentation reference. Note that it describes the schema of the requests the webhook server replies to. This schema depends on the currently supported Kubernetes versions.

You can go to editor.swagger.io and paste its contents there to see a rendered version.

### Performance

Please refer to [docs/performance.md](docs/performance.md).
