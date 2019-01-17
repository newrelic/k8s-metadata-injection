# Kubernetes Metadata injection for New Relic APM agents

## How does it work ?

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

```bash
$ kubectl apply -f newrelic-metadata-injection.yaml
```

Executing this
- creates `newrelic-metadata-injection-deployment` and `newrelic-metadata-injection-svc`;
- registers the `newrelic-metadata-injection-svc` service as a MutatingAdmissionWebhook with the Kubernetes API.

Then, if you wish to let the certificate management be automatic using the Kubernetes extension API server (recommended, but optional):

```bash
$ kubectl apply -f deploy/job.yaml
```

Otherwise, if you are managing the certificate manually you will have to create the TLS secret with the signed certificate/key pair:

```bash
kubectl create secret tls newrelic-metadata-injection-secret \
      --key=server-key.pem \
      --cert=signed-server-cert.pem \
      --dry-run -o yaml |
  kubectl -n default apply -f -
```

Either certificate management choice made, the important thing is to have the secret created with the correct name and namespace. As long as this is done the webhook server will be able to pick it up.

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

### Performance

Please refer to [docs/performance.md](docs/performance.md).