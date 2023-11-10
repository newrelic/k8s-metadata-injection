# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### enhancement
- Replace k8s v1.28.0-rc.1 with k8s 1.28.3 support by @svetlanabrennan in [#458](https://github.com/newrelic/k8s-metadata-injection/pull/458)

## v1.20.0 - 2023-11-06

### 🛡️ Security notices
- Pin Slack notification action to a hash, not to a tag by @juanjjaramillo in [#447](https://github.com/newrelic/k8s-metadata-injection/pull/447)

## v1.19.0 - 2023-10-30

### 🚀 Enhancements
- Remove 1.23 support by @svetlanabrennan in [#441](https://github.com/newrelic/k8s-metadata-injection/pull/441)
- Add k8s 1.28.0-rc.1 support by @svetlanabrennan in [#443](https://github.com/newrelic/k8s-metadata-injection/pull/443)
- Upload sarif when running periodically or pushing to main by @juanjaramillo in [#444](https://github.com/newrelic/k8s-metadata-injection/pull/444)
- Improve Trivy scan by using Docker image by @juanjjaramillo in [#446](https://github.com/newrelic/k8s-metadata-injection/pull/446)

## v1.18.4 - 2023-10-23

### 🐞 Bug fixes
- Trivy scans should only run on the 'Security' workflow by juanjjaramillo in [#436](https://github.com/newrelic/k8s-metadata-injection/pull/436)

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.3
- Updated github.com/fsnotify/fsnotify to v1.7.0 - [Changelog 🔗](https://github.com/fsnotify/fsnotify/releases/tag/v1.7.0)

## v1.18.3 - 2023-10-16

### 🐞 Bug fixes
- Address CVE-2023-44487 and CVE-2023-39325 by juanjjaramillo in [#434](https://github.com/newrelic/k8s-metadata-injection/pull/434)

## v1.18.2 - 2023-10-09

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.14.0

## v1.18.1 - 2023-10-02

### 🐞 Bug fixes
- Fix release workflow to include build-time metadata on release image by juanjjaramillo in [#425](https://github.com/newrelic/k8s-metadata-injection/pull/425)

## v1.18.0 - 2023-09-29

### 🚀 Enhancements
- Improve readability of `release-integration-reusable.yml` by @juanjjaramillo in [#422](https://github.com/newrelic/k8s-metadata-injection/pull/422)

## v1.17.0 - 2023-09-29

### 🚀 Enhancements
- Make explicit that we are only using a single file by @juanjjaramillo in [#416](https://github.com/newrelic/k8s-metadata-injection/pull/416)

### 🐞 Bug fixes
- Fix action to fetch `version-update.go` by @juanjjaramillo in [#420](https://github.com/newrelic/k8s-metadata-injection/pull/420)
- Add quotation to variables to handle spaces by @juanjjaramillo in [#417](https://github.com/newrelic/k8s-metadata-injection/pull/417)

### ⛓️ Dependencies
- Updated alpine to v3.18.4

## v1.16.1 - 2023-09-26

### ⛓️ Dependencies
- Updated go.uber.org/zap to v1.26.0

## v1.16.0 - 2023-09-21

### 🚀 Enhancements
- update contributing.md docs by @svetlanabrennan in [#389](https://github.com/newrelic/k8s-metadata-injection/pull/389)

## v1.15.2 - 2023-09-20

### ⛓️ Dependencies
- Updated go.uber.org/zap to v1.26.0

## v1.15.1 - 2023-09-18

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.2
- Updated go to 1.21
- Updated golang.org/x/crypto to v0.13.0

## v1.15.0 - 2023-09-11

### 🚀 Enhancements
- Update K8s Versions in E2E Tests by @xqi-nr in [#369](https://github.com/newrelic/k8s-metadata-injection/pull/369)

## v1.14.1 - 2023-09-04

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.1

## v1.14.0 - 2023-08-31

### 🚀 Enhancements
- Remove old maintainers @svetlanabrennan [#355](https://github.com/newrelic/k8s-metadata-injection/pull/355)

## v1.13.0 - 2023-08-28

### 🚀 Enhancements
- Define GitHub bot name and email @juanjjaramillo [#343](https://github.com/newrelic/k8s-metadata-injection/pull/343)

### ⛓️ Dependencies
- Updated alpine to v3.18.3

## v1.12.0 - 2023-08-23

### 🛡️ Security notices
- Meet internal security standards @juanjjaramillo [#334](https://github.com/newrelic/k8s-metadata-injection/pull/334)

## 1.11.0
## What's Changed
- Add configuration of certmanager durations @cdobbyn [#323](https://github.com/newrelic/k8s-metadata-injection/pull/323)
- Add changelog workflow @svetlanabrennan [#316](https://github.com/newrelic/k8s-metadata-injection/pull/316)
- Update code owners @juanjjaramillo [#318](https://github.com/newrelic/k8s-metadata-injection/pull/318)
- Add pull request template @svetlanabrennan [#317](https://github.com/newrelic/k8s-metadata-injection/pull/317)
- Add More Logs for NEW_RELIC_METADATA_KUBERNETES_CLUSTER_NAME Injection @xqi-nr [#325](https://github.com/newrelic/k8s-metadata-injection/pull/325)

**Full Changelog**: https://github.com/newrelic/k8s-metadata-injection/compare/v1.10.2...v1.11.0

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
