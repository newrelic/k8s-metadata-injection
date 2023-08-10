# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### enhancement
- Add changelog workflow @svetlanabrennan [#316](https://github.com/newrelic/k8s-metadata-injection/pull/316)
- Update code owners @jjaramillo [#318](https://github.com/newrelic/k8s-metadata-injection/pull/318)
- Add pull request template @svetlanabrennan [#317](https://github.com/newrelic/k8s-metadata-injection/pull/317)

## 1.10.2
## What's Changed
* Update CHANGELOG.md by @juanjjaramillo in https://github.com/newrelic/k8s-metadata-injection/pull/302
* Bump versions by @juanjjaramillo in https://github.com/newrelic/k8s-metadata-injection/pull/303
* chore(deps): bump aquasecurity/trivy-action from 0.10.0 to 0.11.2 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/304
* chore(deps): bump alpine from 3.18.0 to 3.18.2 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/305
* chore(deps): bump k8s.io/apimachinery from 0.27.2 to 0.27.3 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/306
* chore(deps): bump k8s.io/api from 0.27.2 to 0.27.3 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/307
* upgrade go version by @xqi-nr in https://github.com/newrelic/k8s-metadata-injection/pull/308

**Full Changelog**: https://github.com/newrelic/k8s-metadata-injection/compare/v1.10.1...v1.10.2

## 1.10.1

## What's Changed
* Fix helm unittests by @htroisi in https://github.com/newrelic/k8s-metadata-injection/pull/292
* Bump app and chart versions by @juanjjaramillo in https://github.com/newrelic/k8s-metadata-injection/pull/293
* Update Helm unit test reference by @juanjjaramillo in https://github.com/newrelic/k8s-metadata-injection/pull/294
* chore(deps): bump alpine from 3.17.3 to 3.18.0 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/295
* chore(deps): bump k8s.io/api from 0.27.1 to 0.27.2 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/296
* chore(deps): bump github.com/stretchr/testify from 1.8.2 to 1.8.4 by @dependabot in https://github.com/newrelic/k8s-metadata-injection/pull/301


**Full Changelog**: https://github.com/newrelic/k8s-metadata-injection/compare/v1.10.0...v1.10.1

## 1.10.0

- Update dependencies
- Update rennovate workflow
- Bump Helm chart version

## 1.9.0

- Update dependencies
- Update chart maintainers
- Add support for Pod annotations in batch job pods (#261)

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
