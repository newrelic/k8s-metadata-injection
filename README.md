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

## Setup

### 1) Check if MutatingAdmissionWebhook is enabled on your cluster

This feature requires Kubernetes 1.9 or later. Verify that the kube-apiserver process has the admission-control flag set.

```
$ kubectl api-versions | grep admissionregistration.k8s.io/v1beta1
admissionregistration.k8s.io/v1beta1
```

### 2) Install the injection

```bash
$ kubectl apply -f deploy/newrelic-metadata-injection.yaml
```

Executing this:

- creates `newrelic-metadata-injection-deployment` and `newrelic-metadata-injection-svc`.
- registers the `newrelic-metadata-injection-svc` service as a MutatingAdmissionWebhook with the Kubernetes API.

### 3) Install the certificates

This webhook needs to be authenticated by the Kubernetes extension API server, so it will need to have a signed certificate from a CA trusted by the extension API server. The certificate management is isolated from the webhook server and a secret is used to mount them. 

Important: the webhook server has a file watcher pointed at the secret's folder that will trigger a certificate reload whenever anything is created or modified inside the secret. This allows certificate rotation to be transparent and cause no downtime.

#### Automatic management

The certificate management can be automatic, using the Kubernetes extension API server (recommended, but optional):

```bash
$ kubectl apply -f deploy/job.yaml
```

This manifest contains a service account that has the following **cluster** permissions (**RBAC based**) to be capable of automatically management the certificates:

* MutatingWebhookConfiguration - **get**, **create** and **patch**: to be able to create the webhook and patch its CA bundle.
* CertificateSigningRequests - **create**, **get** and **delete**: to be able to sign the certificate required for the webhook server without leaving duplicates.
* CertificateSigningRequests/Approval - **update**: to be able to approve CertificateSigningRequests.
* Secrets - **create**, **get** and **patch**: to be able to manage the TLS secret used to store the key/cert pair used in the webhook server.
* ConfigMaps - **get**: to be able go get the k8s api server's CA bundle, used in the MutatingWebhookConfiguration.

This job will execute the shell script [k8s-cert-signer/generate_certificate.sh](./k8s-cert-signer/generate_certificate.sh) to setup everything. This script will:

1. Generate a server key.
2. If there is any previous CSR (certificate signing request) for this key, it is deleted.
3. Generate a CSR for such key.
4. The signature of the key is then approved.
5. The server's certificate is fetched from the CSR and then encoded.
6. A secret of type `tls` is created with the server certificate and key.
7. The k8s extension api server's CA bundle is fetched.
8. The mutating webhook configuration for the webhook server is patched with the k8s api server's CA bundle from the previous step. This CA bundle will be used by the k8s extension api server when calling our webhook.

#### Manual management

Otherwise, if you are managing the certificate manually you will have to create the TLS secret with the signed certificate/key pair and patch the webhook's CA bundle:

```bash
$ kubectl create secret tls newrelic-metadata-injection-secret \
      --key=server-key.pem \
      --cert=signed-server-cert.pem \
      --dry-run -o yaml |
  kubectl -n default apply -f -

$ caBundle=$(cat caBundle.pem | base64 | td -d '\n')
$ kubectl patch mutatingwebhookconfiguration newrelic-metadata-injection-cfg --type='json' -p "[{'op': 'replace', 'path': '/webhooks/0/clientConfig/caBundle', 'value':'${caBundle}'}]"
```

Either certificate management choice made, the important thing is to have the secret created with the correct name and namespace. As long as this is done the webhook server will be able to pick it up.

### 3) Enable the automatic Kubernetes metadata injection on your namespaces

The injection is only applied to namespaces that have the `newrelic-metadata-injection` label set to `enabled`.

```
$ kubectl label namespace <namespace> newrelic-metadata-injection=enabled
```

## Development

### Prerequisites

For the development process [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube) and [Skaffold](https://github.com/GoogleCloudPlatform/skaffold) tools are used.

* [Install Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/).
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

Run `make deploy-dev`. This will compile your binary with compatibility for the container OS architecture, build a temporary docker image, and finally deploy the webhook server to your Minikube and use the Kubernetes API server to sign its TLS certificate ([see section about certificates](#3-install-the-certificates)).

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

You can go to [editor.swagger.io](editor.swagger.io) and paste its contents there to see a rendered version.

### Performance

Please refer to [docs/performance.md](docs/performance.md).
