# k8s-host-change Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.7.1] - 2025-03-28
### Removed
- [#22] etcd-dependency from helm-chart 

## [v0.7.0] - 2024-11-19
### Added
- [#20] Add RBAC permission to read global-config

## [v0.6.0] - 2024-10-28
### Changed
- [#18] Make imagePullSecrets configurable via helm values and use `ces-container-registries` as default.

## [v0.5.0] - 2024-09-19
### Changed
- [#16] Relicense to AGPL-3.0-only

## [v0.4.0] - 2024-07-25
### Changed
- [#13] Use k8s-registry-lib to read and write configs

## [v0.3.2] - 2023-12-11
### Fixed
- [#11] Use correct key in patch templates.

## [v0.3.1] - 2023-12-07
### Added
- [#9] Add component patch templates used for airgapped environments.
- [#9] Update makefiles to 9.0.1 and hold all yaml resources in a single helm chart.

## [v0.3.0] - 2023-09-15
### Changed
- [#7] Move component-dependencies to helm-annotations

## [v0.2.0] - 2023-08-31
### Added
- [#5] Add helm-chart with dependencies and release the helm-chart in the build-pipeline

## [v0.1.1] - 2023-07-13

### Changed
- [#3] Improved doc comments and split complex components

## [v0.1.0] - 2023-03-27

### Added
- initial release
- [#1] support changes in split dns configuration
