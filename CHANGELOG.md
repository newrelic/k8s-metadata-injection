# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.9.0

- Updated dependencies

## 1.8.0

- Updated dependencies
- Fix: Resolve the issue about MutatingWebhookConfiguration being not supported in v1beta1

## 1.7.5

- Updated dependencies

## 1.7.4

- Fix: Update dependencies to address vulnerability issue (#234)

## 1.7.3

- Updated dependencies
- Fix: Re-enable trivy for high vulnerabilities (#202)

## 1.7.2

- Fix: Update transitive dependencies to address trivy vulnerability issue (#164)

## 1.7.1

- Updated dependencies

## 1.7.0

- Updated dependencies

## 1.6.0

- Adds support for Kubernetes 1.22

## 1.5.0

- Dependencies have been updated to their latest versions (#93)

## 1.4.0

- Support multiarch images

## 1.3.2

- Update k8s-webhook-cert-manager to 1.3.2
  This new version introduces a fix to extend support to version 1.19.x of Kubernetes

## 1.3.1

- Use Github Actions for releasing

## 1.3.0

- Update k8s-webhook-cert-manager to 1.3.0
- Update to golang version 1.14.6 and alpine 3.12.0

## 1.2.0

- Update deployment apiVersion to apps/v1
- Add kubernetes.io/legacy-unknown signer with approve permission to rbac for 1.18 compatibility

## 1.1.4

- Update k8s-webhook-cert-manager to 1.2.1

## 1.1.3

- Update k8s-webhook-cert-manager to 1.2.0

## 1.1.2

- Update k8s-webhook-cert-manager to 1.1.1

## 1.1.1

- Change default server timeout to 1s

## 1.1.0

### Added

- OpenShift support!

### Changed

- Deployment and Service resources are now explicitly assigned to a namespace.
- The webhook server now listens on a non-root port by default: 8443.

## 1.0.0

- Initial version of the webhook.
