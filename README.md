# Kubernetes Metadata injection for New Relic APM agents

## How does it work?

New Relic requires the following environment variables to identify Kubernetes objects in the APM agents:

- `NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NODE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NAMESPACE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_DEPLOYMENT_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_POD_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_CONTAINER_NAME`

These environment variables can be set manually by the customer, or they can be automatically injected using a `MutatingAdmissionWebhook`.

New Relic provides an easy method for deploying this automatic approach.

## Prerequisites

Your cluster need to have the `MutatingAdmissionWebhook` controller enabled. This feature requires Kubernetes 1.9 or later. Verify that the kube-apiserver process has the admission-control flag set by running:

```bash
$ kubectl api-versions | grep admissionregistration.k8s.io/v1beta1
admissionregistration.k8s.io/v1beta1
```

## Installation

There are 2 different behavior options for the metadata injection:

- **The default**, to inject on all pods except those in the `kube-system` and `kube-public` namespaces. This doesn't require any change.
- Only inject on pods created in specific namespaces, based on their labels.

If you wish to use the label-based solution: download the yaml file, look for `namespaceSelector`, uncomment it and the 2 other lines right below (until `newrelic-metadata-injection: enabled`). There are comments in this part of the yaml to help you uncomment the correcty lines.

This webhook needs to be authenticated by the Kubernetes extension API server, so it will need to have a signed certificate from a CA trusted by the extension API server. The certificate management is isolated from the webhook server and a secret is used to mount them. So there are 2 installation options:

1. [With **automatic certificate** management](#automatic-certificate-management), recommended if you don't manually manage your certificates.
2. [With **manual certificate** management](#manual-certificate-management), if you like to manage your certificate manually.

Installing the injection, no matter which options chosen, will end up creating these resources in your Kubernetes cluster:

- `newrelic-metadata-injection-deployment` and `newrelic-metadata-injection-svc`.
- `newrelic-metadata-injection-cfg`, registering the `newrelic-metadata-injection-svc` service as a MutatingAdmissionWebhook with the Kubernetes API.

### Automatic certificate management

The certificate management can be automatic, using the Kubernetes extension API server (recommended, but optional):

```bash
$ kubectl apply -f http://download.newrelic.com/infrastructure_agent/integrations/k8s-metadata-injection-latest.yaml
```

This manifest contains a service account that has the following **cluster** permissions (**RBAC based**) to be capable of automatically manage the certificates:

* `MutatingWebhookConfiguration` - **get**, **create** and **patch**: to be able to create the webhook and patch its CA bundle.
* `CertificateSigningRequests` - **create**, **get** and **delete**: to be able to sign the certificate required for the webhook server without leaving duplicates.
* `CertificateSigningRequests/Approval` - **update**: to be able to approve CertificateSigningRequests.
* `Secrets` - **create**, **get** and **patch**: to be able to manage the TLS secret used to store the key/cert pair used in the webhook server.
* `ConfigMaps` - **get**: to be able to get the k8s api server's CA bundle, used in the MutatingWebhookConfiguration.

It also includes a job will execute the shell script [generate_certificate.sh](https://github.com/newrelic/k8s-webhook-cert-manager/blob/master/generate_certificate.sh) to setup everything. This script will:

1. Generate a server key.
2. If there is any previous CSR (certificate signing request) for this key, it is deleted.
3. Generate a CSR for such key.
4. The signature of the key is then approved.
5. The server's certificate is fetched from the CSR and then encoded.
6. A secret of type `tls` is created with the server certificate and key.
7. The k8s extension api server's CA bundle is fetched.
8. The mutating webhook configuration for the webhook server is patched with the k8s api server's CA bundle from the previous step. This CA bundle will be used by the k8s extension api server when calling our webhook.

If you wish to learn more about TLS certificates management inside Kubernetes, check out [the official documentation for Managing TLS Certificate in a Cluster](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#create-a-certificate-signing-request-object-to-send-to-the-kubernetes-api).

### Manual certificate management

The first step in this case is to have the server's key and certificate already downloaded to your machine and named as `server-key.pem` and `signed-server-cert.pem`.

Then run the following commands to install the injection, create the TLS secret and patch the `MutatingWebhookConfiguration` with the correct **caBundle**:

```bash
$ kubectl apply -f http://download.newrelic.com/infrastructure_agent/integrations/k8s-metadata-injection-custom-certs-latest.yaml

$ kubectl create secret tls newrelic-metadata-injection-secret \
      --key=server-key.pem \
      --cert=signed-server-cert.pem \
      --dry-run -o yaml |
  kubectl -n default apply -f -

$ caBundle=$(cat caBundle.pem | base64 | td -d '\n')
$ kubectl patch mutatingwebhookconfiguration newrelic-metadata-injection-cfg --type='json' -p "[{'op': 'replace', 'path': '/webhooks/0/clientConfig/caBundle', 'value':'${caBundle}'}]"
```

### Certificate rotation

**Important**: the webhook server has a file watcher pointed at the secret's folder that will trigger a certificate reload whenever anything is created or modified inside the secret. This allows easy certificate rotation with an update of the TLS secret. The following snippet is an example of updating the TLS secret with a new server certificate (the key stays the same):

```bash
$ namespace=default # Change the namespace here if you also changed it in the yaml files.
$ serverCert=$(kubectl get csr newrelic-metadata-injection-svc.${namespace} -o jsonpath='{.status.certificate}')
$ tmpdir=$(mktemp -d)
$ echo ${serverCert} | openssl base64 -d -A -out ${tmpdir}/server-cert.pem
$ kubectl patch secret newrelic-metadata-injection-secret --type='json' \
    -p "[{'op': 'replace', 'path':'/data/tls.crt', 'value':'$(serverCert)'}]"
$ rm -rf $(tmpdir)
```

### Usage

If using the default injection installation, all the pods you create will automatically have the correct environment variables set.

If using the label-based injection option, the injection is only applied to namespaces that have the `newrelic-metadata-injection` label set to `enabled`.

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

Run `skaffold run`. This will build a docker image, build the webhook server inside it, and finally deploy the webhook server to your Minikube and use the Kubernetes API server to sign its TLS certificate ([see section about certificates](#automatic-certificate-management)).

To follow the logs, you can run `skaffold run --tail`. To delete the resources created by Skaffold you can run `skaffold delete`.

If you would like to enable automatic redeploy on changes to the repository, you can run `skaffold dev`. It automatically tails the logs and delete the resources when interrupted (i.e. with a `Ctrl + C`).

### Tests

For running unit tests, use

```bash
make test
```

For running benchmark tests, use:

```bash
make benchmark-test
```

There are also some basic E2E tests, they are prepared to run using
[Minikube](https://github.com/kubernetes/minikube). To run them, execute:

``` bash
make e2e-test
```

You can specify against which version of K8s you want to execute the tests:

``` bash
E2E_KUBERNETES_VERSION=v1.10.0 E2E_START_MINIKUBE=yes make e2e-test
```

### Documentation

Please use the [Open Api 3.0 spec file](openapi.yaml) as documentation reference. Note that it describes the schema of the requests the webhook server replies to. This schema depends on the currently supported Kubernetes versions.

You can go to [editor.swagger.io](editor.swagger.io) and paste its contents there to see a rendered version.

### Performance

Please refer to [docs/performance.md](docs/performance.md).
