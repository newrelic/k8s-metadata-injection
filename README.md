# Kubernetes Metadata injection for New Relic APM agents

[![Build Status](https://travis-ci.com/newrelic/k8s-metadata-injection.svg?branch=master)](https://travis-ci.com/newrelic/k8s-metadata-injection) [![Go Report Card](https://goreportcard.com/badge/github.com/newrelic/k8s-metadata-injection)](https://goreportcard.com/report/github.com/newrelic/k8s-metadata-injection)

## Documentation

If you wish to read higher-level documentation about this project, please, visit the [official documentation site](https://docs.newrelic.com/docs/integrations/kubernetes-integration/metadata-injection/kubernetes-apm-metadata-injection).

# How does it work?

New Relic APM agents requires the following environment variables to provide Kubernetes object information in the context of an specific application distributed trace, transaction trace or error trace.

- `NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NODE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_NAMESPACE_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_DEPLOYMENT_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_POD_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_CONTAINER_NAME`
- `NEW_RELIC_METADATA_KUBERNETES_CONTAINER_IMAGE_NAME`

These environment variables are automatically injected in the pods using a MutatingAdmissionWebhook provided by this project.

Please refer to the [official documentation](https://docs.newrelic.com/docs/integrations/kubernetes-integration/metadata-injection/kubernetes-apm-metadata-injection) to learn more about the reasoning behind this project.

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

Run `skaffold run`. This will build a docker image, build the webhook server inside it, and finally deploy the webhook server to your Minikube and use the Kubernetes API server to sign its TLS certificate ([see section about certificates](#3-install-the-certificates)).

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

### API Documentation

Please use the [Open Api 3.0 spec file](openapi.yaml) as documentation reference. Note that it describes the schema of the requests the webhook server replies to. This schema depends on the currently supported Kubernetes versions.

You can go to [editor.swagger.io](editor.swagger.io) and paste its contents there to see a rendered version.

### Performance

Please refer to [docs/performance.md](docs/performance.md).

## Certificates management

Admission webhooks are called by the Kubernetes API server and it needs to authenticate the webhooks using TLS. In this project we offer 2 different options of certificate management.

Either certificate management choice made, the important thing is to have the secret created with the correct name and namespace, and also to have the correct CA bundle in the MutatingWebhookConfiguration resource. As long as this is done the webhook server will be able to pick it up.

### Automatic 

Please refer to the [setup instructions in the official documentation](https://docs.newrelic.com/docs/integrations/kubernetes-integration/metadata-injection/kubernetes-apm-metadata-injection#install).

For the automatic certificate management, the [k8s-webhook-cert-manager](https://github.com/newrelic/k8s-webhook-cert-manager) is used. Feel free to check the repository to know more about it.

The manifest file at [deploy/job.yaml](./deploy/job.yaml) contains a service account that has the following **cluster** permissions (**RBAC based**) to be capable of automatically manage the certificates:

* `MutatingWebhookConfiguration` - **get**, **create** and **patch**: to be able to create the webhook and patch its CA bundle.
* `CertificateSigningRequests` - **create**, **get** and **delete**: to be able to sign the certificate required for the webhook server without leaving duplicates.
* `CertificateSigningRequests/Approval` - **update**: to be able to approve CertificateSigningRequests.
* `Secrets` - **create**, **get** and **patch**: to be able to manage the TLS secret used to store the key/cert pair used in the webhook server.
* `ConfigMaps` - **get**: to be able go get the k8s api server's CA bundle, used in the MutatingWebhookConfiguration.

If you wish to learn more about TLS certificates management inside Kubernetes, check out [the official documentation for Managing TLS Certificates in a Cluster](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#create-a-certificate-signing-request-object-to-send-to-the-kubernetes-api).

### Custom

Otherwise, if you want to use the custom certificate management option you have to create the TLS secret with the signed certificate/key pair and patch the webhook's CA bundle:

```bash
$ kubectl create secret tls newrelic-metadata-injection-secret \
      --key=server-key.pem \
      --cert=signed-server-cert.pem \
      --dry-run -o yaml |
  kubectl -n default apply -f -

$ caBundle=$(cat caBundle.pem | base64 | td -d '\n')
$ kubectl patch mutatingwebhookconfiguration newrelic-metadata-injection-cfg --type='json' -p "[{'op': 'replace', 'path': '/webhooks/0/clientConfig/caBundle', 'value':'${caBundle}'}]"
```

## Contributing

We welcome code contributions (in the form of pull requests) from our user community. Before submitting a pull request please review [these guidelines](./CONTRIBUTING.md).

Following these helps us efficiently review and incorporate your contribution and avoid breaking your code with future changes.

## Release a new version

- Update the version in `deploy/newrelic-metadata-injection.yaml`.
- Create a Github release.
- Launch the `k8s-metadata-injection-release` job in Jenkins.
