<a href="https://opensource.newrelic.com/oss-category/#community-plus"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Community_Plus.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Plus.png"><img alt="New Relic Open Source community plus project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Plus.png"></picture></a>

# Kubernetes Metadata injection for New Relic APM agents - test

[![Build Status](https://travis-ci.com/newrelic/k8s-metadata-injection.svg?branch=main)](https://travis-ci.com/newrelic/k8s-metadata-injection) [![Go Report Card](https://goreportcard.com/badge/github.com/newrelic/k8s-metadata-injection)](https://goreportcard.com/report/github.com/newrelic/k8s-metadata-injection)

# Table of contents

- [Documentation](#documentation)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Dependency management](#dependency-management)
  - [Configuration](#configuration)
  - [Run](#run)
  - [Tests](#tests)
  - [API Documentation](#api-documentation)
  - [Performance](#performance)
- [Certificates management](#certificates-management)
  - [Automatic](#automatic)
  - [Custom](#custom)
- [Contributing](#contributing)
- [License](#license)
- [Release a new version](#release-a-new-version)

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

## Helm chart

You can install this integration using [`nri-bundle` helm chart](https://github.com/newrelic/helm-charts/tree/master/charts/nri-bundle) located in the
[helm-charts repository](https://github.com/newrelic/helm-charts) or directly from this repository by adding this Helm repository:

```shell
helm repo add nri-metadata-injection https://newrelic.github.io/k8s-metadata-injection
helm upgrade --install nri-metadata-injection/nri-metadata-injection -f your-custom-values.yaml
```

For further information of the configuration needed for the chart just read the [chart's README](/charts/nri-metadata-injection/README.md).

## Development

### Prerequisites

For the development process [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube) and [Skaffold](https://github.com/GoogleCloudPlatform/skaffold) tools are used.

- [Install Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/).
- [Install Skaffold](https://github.com/GoogleCloudPlatform/skaffold#installation).

Currently the project compiles with **Go 1.11.4**.

### Dependency management

[Go modules](https://github.com/golang/go/wiki/Modules) are used for managing dependencies. This project does not need to be in your GOROOT, if you wish so.

Currently for K8s libraries it uses version 1.13.1. Only couple of libraries are direct dependencies, the rest are indirect. You need to point all of them to the same K8s version to make sure that everything works as expected. For the moment this process is manual.

### Configuration

- Copy the deployment file `deploy/newrelic-metadata-injection.yaml` to `deploy/local.yaml`.
- Edit the file and set the following value as container image: `internal/k8s-metadata-injector`.
- Make sure that `imagePullPolicy: Always` is not present in the file (otherwise, the image won't be pulled).

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

The e2e tests make the assumption that you are running on an AMD system so in case the test doesn't generate the needed binary, run the below command. 
For instance if you run on an M2 mac arm64 is the target arch but it is not made by default. 

```bash
make compile-multiarch
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

- `MutatingWebhookConfiguration` - **get**, **create** and **patch**: to be able to create the webhook and patch its CA bundle.
- `CertificateSigningRequests` - **create**, **get** and **delete**: to be able to sign the certificate required for the webhook server without leaving duplicates.
- `CertificateSigningRequests/Approval` - **update**: to be able to approve CertificateSigningRequests.
- `Secrets` - **create**, **get** and **patch**: to be able to manage the TLS secret used to store the key/cert pair used in the webhook server.
- `ConfigMaps` - **get**: to be able go get the k8s api server's CA bundle, used in the MutatingWebhookConfiguration.

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

## Release a new version

- Update the version in `deploy/newrelic-metadata-injection.yaml`.
- Update the version in `WEBHOOK_DOCKER_IMAGE_TAG` in the `Makefile`.
- Create a Github release.
- Launch the `k8s-metadata-injection-release` job in Jenkins.

## Support

Should you need assistance with New Relic products, you are in good hands with several support diagnostic tools and support channels.

>New Relic offers NRDiag, [a client-side diagnostic utility](https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/troubleshooting/new-relic-diagnostics) that automatically detects common problems with New Relic agents. If NRDiag detects a problem, it suggests troubleshooting steps. NRDiag can also automatically attach troubleshooting data to a New Relic Support ticket. Remove this section if it doesn't apply.

If the issue has been confirmed as a bug or is a feature request, file a GitHub issue.

**Support Channels**

- [New Relic Documentation](https://docs.newrelic.com): Comprehensive guidance for using our platform
- [New Relic Community](https://forum.newrelic.com/t/new-relic-kubernetes-open-source-integration/109093): The best place to engage in troubleshooting questions
- [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
- [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level
- [New Relic Technical Support](https://support.newrelic.com/) 24/7/365 ticketed support. Read more about our [Technical Support Offerings](https://docs.newrelic.com/docs/licenses/license-information/general-usage-licenses/support-plan).

## Privacy

At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).

## Contribute

We encourage your contributions to improve this project! Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA (which is required if your contribution is on behalf of a company), drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project would not be what it is today.

## License

Kubernetes Metadata injection is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
