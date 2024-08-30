# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## v1.28.4 - 2024-08-12

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.26.0

## v1.28.3 - 2024-07-29

### ⛓️ Dependencies
- Updated alpine to v3.20.2
- Updated kubernetes packages to v0.30.3

## v1.28.2 - 2024-07-22

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.25.0

## v1.28.1 - 2024-07-08

### ⛓️ Dependencies
- Updated kubernetes packages to v0.30.2

## v1.28.0 - 2024-06-24

### 🚀 Enhancements
- Add 1.29 and 1.30 support and drop 1.25 and 1.24 @dbudziwojskiNR [#551](https://github.com/newrelic/k8s-metadata-injection/pull/551)

### ⛓️ Dependencies
- Updated alpine to v3.20.1

## v1.27.4 - 2024-06-17

### ⛓️ Dependencies
- Updated go to v1.22.4
- Updated golang.org/x/crypto to v0.24.0

## v1.27.3 - 2024-06-10

### ⛓️ Dependencies
- Updated go to v1.22.3

## v1.27.2 - 2024-05-27

### ⛓️ Dependencies
- Updated alpine to v3.20.0

## v1.27.1 - 2024-05-13

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.23.0

## v1.27.0 - 2024-04-29

### ⛓️ Dependencies
- Upgraded golang.org/x/net from 0.21.0 to 0.23.0

## v1.26.4 - 2024-04-15

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.22.0

## v1.26.3 - 2024-03-25

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.3

## v1.26.2 - 2024-03-11

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.21.0

## v1.26.1 - 2024-03-04

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.2

## v1.26.0 - 2024-02-26

### 🚀 Enhancements
- Add linux node selector @dbudziwojskiNR [#523](https://github.com/newrelic/k8s-metadata-injection/pull/523)

### ⛓️ Dependencies
- Updated go.uber.org/zap to v1.27.0

## v1.25.1 - 2024-02-19

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.19.0

## v1.25.0 - 2024-02-05

### 🚀 Enhancements
- Add Codecov @dbudziwojskiNR [#513](https://github.com/newrelic/k8s-metadata-injection/pull/513)

## v1.24.2 - 2024-01-29

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.1
- Updated alpine to v3.19.1

## v1.24.1 - 2024-01-22

### ⛓️ Dependencies
- Updated go to v1.21.6

## v1.24.0 - 2024-01-15

### 🚀 Enhancements
- Trigger release creation by @juanjjaramillo [#506](https://github.com/newrelic/k8s-metadata-injection/pull/506)
- Remove reusable workflows by @juanjjaramillo [#491](https://github.com/newrelic/k8s-metadata-injection/pull/491)

## v1.23.2 - 2024-01-08

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.0
- Updated golang.org/x/crypto to v0.18.0

## v1.23.1 - 2023-12-25

### 🚀 Enhancements
- Update e2e testing workflow to also run on release in [#485](https://github.com/newrelic/k8s-metadata-injection/pull/485)

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.4
- Updated go to v1.21.5
- Updated alpine to v3.19.0

## v1.23.0 - 2023-12-06

### 🚀 Enhancements
- Update reusable workflow dependency by @juanjjaramillo [#490](https://github.com/newrelic/k8s-metadata-injection/pull/490)
- Reusable release workflow now provides a mechanism for opting out of helm chart updates [#488](https://github.com/newrelic/k8s-metadata-injection/pull/488)

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.16.0
- Updated alpine to v3.18.5

## v1.22.1 - 2023-11-16

### ⛓️ Dependencies
- Updated golang.org/x/crypto to v0.15.0

## v1.22.0 - 2023-11-13

### 🚀 Enhancements
- Update k8s version in e2e tests by @svetlanabrennan in [#459](https://github.com/newrelic/k8s-metadata-injection/pull/459)

## v1.21.0 - 2023-11-13

### 🚀 Enhancements
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
