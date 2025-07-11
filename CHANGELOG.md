# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Update Go to 1.24.5

## [7.0.2] - 2025-06-18

- Dependency updates:
  - nancy v1.0.51

## [7.0.1] - 2025-06-03

- Dependency updates:
  - helm-chart-testing v3.13
  - conftest v0.61.0
  - helm v3.18.2
  - kubeconform v0.7.0

## [7.0.0] - 2025-05-08

### Removed

- Remove `golangci-lint` and the `make lint` target

## [6.20.1] - 2025-04-28

- Downgrade golangci-lint v1.64.8 to v1.64.7, to be able to handle v2 configuration files with v1 keys in them.

## [6.20.0] - 2025-04-16

- Add `generators` flag to `create kustomization` command, defaults to `true` for backward compatibility.
  - When set to `false`, the YAML files in the target directory are assumed to be normal resources
    and will update `.resources` with the list in the `kustomization.yaml`.

## [6.19.1] - 2025-03-14

## [6.19.0] - 2025-03-05

- Dependency updates
  - Base image: golang:1.24.0-alpine3.21
  - Helm v3.17.1
  - kubeconform v0.6.7
  - yamale v6.0.0

## [6.18.2] - 2025-02-13

- Dependency updates
  - golangci-lint v1.64.4
  - yamllint v1.35.1
  - yamale v3.0.8

## [6.18.1] - 2025-01-22

- Dependency updates
  - helm-chart-testing v3.12.0
  - app-build-suite v1.2.8
  - golang 1.23.4
  - conftest v0.56.0

## [6.18.0] - 2024-08-26

### Changed

- Upgrade Go to v1.23.0
- Upgrade golangci-lint to v1.60.3
- Upgrade nancy to v1.0.46
- Upgrade kubeconform to v0.4.14

## [6.17.0] - 2024-08-19

### Changed

- Upgrade `golangci-lint` to `1.60.1`

## [6.16.0] - 2024-08-01

### Changed

- Bump `github.com/giantswarm/gitrepo` to `v0.3.0` to add git tag prefix support. See: https://github.com/giantswarm/gitrepo/releases/tag/v0.3.0.

## [6.15.1] - 2024-07-26

### Changed

- Bump base image to `golang:1.21.12-alpine3.19` to fix reported CVEs.

## [6.15.0] - 2024-07-18

### Changed

- Upgrade `golangci-lint` to `1.59.1`

## [6.14.1] - 2024-01-18

### Changed

- Update Go to v1.21.6
- Access base images from `gsoci.azurecr.io`

## [6.14.0] - 2023-12-11

- added `jq` util to the image

## [6.13.0] - 2023-11-08

### Changed

- Update Go to v1.21.3, golangci-lint to v1.55.2, nancy to v1.0.45, yamllint to v1.32.0

## [6.12.1] - 2023-07-20

### Changed

- Update `golangci-lint` to `v1.53.3`.

## [6.12.0] - 2023-07-20

### Changed

- Update Go to v1.20.6

## [6.11.0] - 2023-04-13

### Changed

- Support dots in version suffix when generating PR Releases [#797](https://github.com/giantswarm/architect/pull/797)

## [6.10.0] - 2023-02-21

### Changed

- Update Go to 1.19.6.

## [6.9.0] - 2023-02-15

### Changed

- Switched from kubeval to kubeconform 0.4.13

## [6.8.0] - 2022-11-21

- Update abs to `v1.1.3`.

## [6.7.0] - 2022-10-04

### Changed

- Update Go to 1.19.1.
- Update Alpine to 3.16.2.
- Update Go CI linter to 1.49.0.

## [6.6.0] - 2022-07-11

### Changed

- Update `nancy` to `v1.0.37`.

## [6.5.0] - 2022-05-26

### Changed

- Use go 1.18 for the `architect` module.
- Update `nancy` to `v1.0.33`.

## [6.4.0] - 2022-05-10

### Changed

- Update `helm` version to `v3.8.1`.
- Update Go to 1.18.1.

## [6.3.0] - 2022-03-04

### Changed
- Update `go` version to `v1.17.8`.

### Removed

- Remove `create argoapp` as Argo is replaced by Flux.

## [6.2.0] - 2022-02-11

### Changed

- Update `go` version to `v1.17.7`.

## [6.1.1] - 2022-02-09

### Fixed

- Skip `kustomization.yaml` when running `architect create kustomization`.

## [6.1.0] - 2022-02-07

### Changed

- Update dependencies.

### Added

- Add `--update-changelog` flag to `prepare-release` command to allow disabling changelog update.

## [6.0.0] - 2022-02-07

- Add `architect create fluxgenerator` and `architect create kustomization` for App collections managed by Flux.

## [5.3.0] - 2021-09-29

### Added

- Add `kubeval` version `v0.16.1`.

## [5.2.0] - 2021-09-10

### Changed

- Update `go` version to `v1.17.1`.
- Update `alpine` version to `3.14.2`.
- Update `golangci-lint` to `v1.42.1`.

## [5.1.0] - 2021-08-25

### Changed

- Update `argoapp` version to `v0.1.4` with cascading deletion enabled.

## [5.0.0] - 2021-08-24

### Removed

- Remove `architect legacy deploy` as draughtsman is retired now.

## [4.0.1] - 2021-07-21

### Added

- Add `gopher` installation.

## [4.0.0] - 2021-07-12

### Added

- Install `make`.

### Changed

- Update kubebuilder from `v2.3.1` to `v3.1.0`.

## [3.7.1] - 2021-06-16

### Changed

- Update abs to `v0.2.3`.

## [3.7.0] - 2021-06-16

### Changed

- Removed decommissioned installations (`archon`, `davis`, `dinosaur`, `dragon`).
- Update chart-testing to `v3.4.0`.

## [3.6.0] - 2021-05-20

### Added

- Add `selfHeal: true` and `allowEmpty: true` to the generated Application CR sync policy in `architect create argoapp` (See [argoapp@v0.1.2](https://github.com/giantswarm/argoapp/blob/main/CHANGELOG.md#012---2021-05-20).

### Fixed

- Temporarily don't fail when Chart.yaml doesn't have the config annotation in `architect create argoapp`.

## [3.5.2] - 2021-05-17

### Fixed

- Fix `architect create argoapp` generated Application CR project (renamed from "draughtsman2" to "collections") by updating to `gaintswarm/argoapp@v0.1.1`.

## [3.5.1] - 2021-05-13

### Added

- Add Beaver.
- Add `--config-ref-from-chart` flag to `architect create argoapp` (#624).

## [3.5.0] - 2021-05-13

### Added

- Add `create argoapp` command.
- Add `yq` into Dockerfile.

## [3.4.4] - 2021-05-06

### Added

- Add Otter.

## [3.4.3] - 2021-04-28

### Added

- Add Eagle.

## [3.4.2] - 2021-03-24

### Added

- Add flamingo.

## [3.4.1] - 2021-03-24

### Added

- Added kubebuilder to the image to be able to run integration tests based on
  controller-runtime `envtest`.

## [3.4.0] - 2021-03-17

### Changed

- Update `go` version to `v1.16.2`.
- Update `helm` version to `v3.5.3`.
- Update `alpine` version to `3.13`.
- Update `conftest` version to `v0.21.0`.
- Update `golangci-lint` version to `v1.38.0`.
- Update `nancy` version to `v1.0.17`.
- Update `helm-chart-testing` to `v3.3.1`.
- Print version in `architect version`.

## [3.3.1] - 2021-03-11

### Changed

- Update `giantswarm/app` to `v4.7.0`.
- Update `github.com/google/go-cmp` to `v0.5.5`.

## [3.3.0] - 2021-02-19

### Changed

- Update `go` version to `v1.16`.

## [3.2.2] - 2021-02-08

### Added

- Add Kudu.

### Changed

- Update `giantswarm/app` to `v4.2.0`.
- Remove Axolotl.

## [3.2.1] - 2021-01-11

### Added

- Prevent deployment to `amagon` (decommissioned).

## [3.2.0] - 2020-12-03

### Added

- Add `--config--version` flag to the `create appcr` command.

## [3.1.1] - 2020-11-27

### Fixed

- Fix app CR configmap and secret flags.

## [3.1.0] - 2020-11-27

### Added

- Allow app CR configmap and secret configuration.

## [3.0.6] - 2020-11-13

### Added

- Added `exodus` installation.

## [3.0.5] - 2020-10-20

### Added

- Added `gremlin` installation.

## [3.0.4] - 2020-10-16

### Fixed

- Fix `prepare-release` when running on multi-digit patch version.

## [3.0.3] - 2020-10-14

## [3.0.2] - 2020-10-14

### Fixed

- Accept alphanumeric strings for release suffix rather than only numbers in prepare-release command.

## [3.0.1] - 2020-10-07

### Added

- Added `orion` installation.

## [3.0.0] - 2020-09-24

### Changed

- Move `deploy` to `legacy deploy` and strip down the functionality to only
  creating GitHub deployment events.
- Update `go` version to `v1.15.2`.

### Removed

- Remove updating module line in go.mod file (if it exists) when major version
  is bigger than 1 in `prepare-release` command added in 2.1.0. It was buggy.
  Expectation is to have a validation instead.
- Remove legacy commands:
    - build
    - publish
    - release
    - unpublish

## [2.1.6] - 2020-08-18

### Added

- Add `nancy` binary to image to use for vulnerability scanning.

## [2.1.5] - 2020-08-13

### Changed

- Remove `avatar` installation from `deploy` command.
- Remove `panther` installation from `deploy` command.
- Remove `platypus` installation from `deploy` command.

### Added

- Add `bandicoot` installation to `deploy` command.

## [2.1.4] - 2020-08-12

### Fixed

- Remove version suffix from reference version before updating `project.go`.

## [2.1.3] - 2020-08-11

### Fixed

- Fix `helm template` rendering for reference versions in.

## [2.1.2] - 2020-08-11

### Added

- Handle release versions like `0.1.0-1` in `prepare-release` command.
- Do not update version in `project.go` file for replacement releases (versions `0.1.0-x`).

## [2.1.1] - 2020-08-05

### Added

- Add `visitor` installation to `deploy` command.

## [2.1.0] - 2020-07-21

### Added

- Add `camel` installation to `deploy` command.
- Update module line in go.mod file (if it exists) when major version is bigger
  than 1 in `prepare-release` command.

### Fixed

- Support "Unreleased" link update for first release on non-master branches in
  `prepare-release` command.

## [2.0.0] - 2020-07-03

### Changed

- Update `helm` binary to `v3.2.4`

## [1.2.0] - 2020-06-08

### Added

- Add `prepare-release` command (#442).

## [1.1.3] 2020-06-05

### Changed

-  Update giantswarm/app to 0.2.2 and use 0.0.0 as version for app CRs.

## [1.1.2] 2020-06-03

### Changed

-  Revert giantswarm/app to 0.2.1 and use 1.0.0 as version for app CRs.

## [1.1.1] 2020-06-02

### Changed

-  Update giantswarm/app to 0.2.2 and use 0.0.0 as version for app CRs.

## [1.1.0] 2020-05-28

### Changed

### Added

- Add Gaia (#428)
- Add argali env (#425)
- Add antelope env (#424)
- Add alpaca env (#423)

### Fixed

- sort out installation alphabetically (#443)
- Skip AppVersion check when project file is absent (#421)
- ensure AppVersion for repos without pkg/project (#416)
- Update giantswarm/app to 0.2.1 (#412)
- Hardcode tag to be 1.0.0 (#409)

### Removed

- remove happa (#441)
- remove api from project list (#439)
- remove credentiald from project list (#438)
- remove passage from project list (#437)
- remove route53-manager (#436)
- events: remove vault-exporter (#435)
- remove cert-exporter from baseProjectList (#433)
- events: remove etcd-backup (#434)
- remove ingress-exporter from base project list (#431)
- events: remove net-exporter from the project list (#430)
- Remove companyd (#429)
- Delete tokend from architect (#420)
- Delete userd (#419)
- Delete kubernetesd (#418)
- Remove cluster-service (#417)
- Remove node-operator (#415)
- Remove flannel-operator (#414)
- Remove bridge-operator (#413)
- Delete g8s-oauth2 (#411)
- Remove g8s-grafana (#410)
- remove cert-operator from baseProjectList (#408)

## [1.0.0] 2020-04-23

### Added

- Add changelog.
- Add SemVer versioning.

[Unreleased]: https://github.com/giantswarm/architect/compare/v7.0.2...HEAD
[7.0.2]: https://github.com/giantswarm/architect/compare/v7.0.1...v7.0.2
[7.0.1]: https://github.com/giantswarm/architect/compare/v7.0.0...v7.0.1
[7.0.0]: https://github.com/giantswarm/architect/compare/v6.20.1...v7.0.0
[6.20.1]: https://github.com/giantswarm/architect/compare/v6.20.0...v6.20.1
[6.20.0]: https://github.com/giantswarm/architect/compare/v6.19.1...v6.20.0
[6.19.1]: https://github.com/giantswarm/architect/compare/v6.19.0...v6.19.1
[6.19.0]: https://github.com/giantswarm/architect/compare/v6.18.2...v6.19.0
[6.18.2]: https://github.com/giantswarm/architect/compare/v6.18.1...v6.18.2
[6.18.1]: https://github.com/giantswarm/architect/compare/v6.18.0...v6.18.1
[6.18.0]: https://github.com/giantswarm/architect/compare/v6.17.0...v6.18.0
[6.17.0]: https://github.com/giantswarm/architect/compare/v6.16.0...v6.17.0
[6.16.0]: https://github.com/giantswarm/architect/compare/v6.15.1...v6.16.0
[6.15.1]: https://github.com/giantswarm/architect/compare/v6.15.0...v6.15.1
[6.15.0]: https://github.com/giantswarm/architect/compare/v6.14.1...v6.15.0
[6.14.1]: https://github.com/giantswarm/architect/compare/v6.14.0...v6.14.1
[6.14.0]: https://github.com/giantswarm/architect/compare/v6.13.0...v6.14.0
[6.13.0]: https://github.com/giantswarm/architect/compare/v6.12.1...v6.13.0
[6.12.1]: https://github.com/giantswarm/architect/compare/v6.12.0...v6.12.1
[6.12.0]: https://github.com/giantswarm/architect/compare/v6.11.0...v6.12.0
[6.11.0]: https://github.com/giantswarm/architect/compare/v6.10.0...v6.11.0
[6.10.0]: https://github.com/giantswarm/architect/compare/v6.9.0...v6.10.0
[6.9.0]: https://github.com/giantswarm/architect/compare/v6.8.0...v6.9.0
[6.8.0]: https://github.com/giantswarm/architect/compare/v6.7.0...v6.8.0
[6.7.0]: https://github.com/giantswarm/architect/compare/v6.6.0...v6.7.0
[6.6.0]: https://github.com/giantswarm/architect/compare/v6.5.0...v6.6.0
[6.5.0]: https://github.com/giantswarm/architect/compare/v6.4.0...v6.5.0
[6.4.0]: https://github.com/giantswarm/architect/compare/v6.3.0...v6.4.0
[6.3.0]: https://github.com/giantswarm/architect/compare/v6.2.0...v6.3.0
[6.2.0]: https://github.com/giantswarm/architect/compare/v6.1.1...v6.2.0
[6.1.1]: https://github.com/giantswarm/architect/compare/v6.1.0...v6.1.1
[6.1.0]: https://github.com/giantswarm/architect/compare/v6.0.0...v6.1.0
[6.0.0]: https://github.com/giantswarm/architect/compare/v5.3.0...v6.0.0
[5.3.0]: https://github.com/giantswarm/architect/compare/v5.2.0...v5.3.0
[5.2.0]: https://github.com/giantswarm/architect/compare/v5.1.0...v5.2.0
[5.1.0]: https://github.com/giantswarm/architect/compare/v5.0.0...v5.1.0
[5.0.0]: https://github.com/giantswarm/architect/compare/v4.0.1...v5.0.0
[4.0.1]: https://github.com/giantswarm/architect/compare/v4.0.0...v4.0.1
[4.0.0]: https://github.com/giantswarm/architect/compare/v3.7.1...v4.0.0
[3.7.1]: https://github.com/giantswarm/architect/compare/v3.7.0...v3.7.1
[3.7.0]: https://github.com/giantswarm/architect/compare/v3.6.0...v3.7.0
[3.6.0]: https://github.com/giantswarm/architect/compare/v3.5.2...v3.6.0
[3.5.2]: https://github.com/giantswarm/architect/compare/v3.5.1...v3.5.2
[3.5.1]: https://github.com/giantswarm/architect/compare/v3.5.0...v3.5.1
[3.5.0]: https://github.com/giantswarm/architect/compare/v3.4.4...v3.5.0
[3.4.4]: https://github.com/giantswarm/architect/compare/v3.4.3...v3.4.4
[3.4.3]: https://github.com/giantswarm/architect/compare/v3.4.2...v3.4.3
[3.4.2]: https://github.com/giantswarm/architect/compare/v3.4.1...v3.4.2
[3.4.1]: https://github.com/giantswarm/architect/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/giantswarm/architect/compare/v3.3.1...v3.4.0
[3.3.1]: https://github.com/giantswarm/architect/compare/v3.3.0...v3.3.1
[3.3.0]: https://github.com/giantswarm/architect/compare/v3.2.2...v3.3.0
[3.2.2]: https://github.com/giantswarm/architect/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/giantswarm/architect/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/giantswarm/architect/compare/v3.1.1...v3.2.0
[3.1.1]: https://github.com/giantswarm/architect/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/giantswarm/architect/compare/v3.0.6...v3.1.0
[3.0.6]: https://github.com/giantswarm/architect/compare/v3.0.5...v3.0.6
[3.0.5]: https://github.com/giantswarm/architect/compare/v3.0.4...v3.0.5
[3.0.4]: https://github.com/giantswarm/architect/compare/v3.0.3...v3.0.4
[3.0.3]: https://github.com/giantswarm/architect/compare/v3.0.2...v3.0.3
[3.0.2]: https://github.com/giantswarm/architect/compare/v3.0.1...v3.0.2
[3.0.1]: https://github.com/giantswarm/architect/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/giantswarm/architect/compare/v2.1.6...v3.0.0
[2.1.6]: https://github.com/giantswarm/architect/compare/v2.1.5...v2.1.6
[2.1.5]: https://github.com/giantswarm/architect/compare/v2.1.4...v2.1.5
[2.1.4]: https://github.com/giantswarm/architect/compare/v2.1.3...v2.1.4
[2.1.3]: https://github.com/giantswarm/architect/compare/v2.1.2...v2.1.3
[2.1.2]: https://github.com/giantswarm/architect/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/giantswarm/architect/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/giantswarm/architect/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/giantswarm/architect/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/giantswarm/architect/compare/v1.1.2...v1.2.0
[1.1.2]: https://github.com/giantswarm/architect/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/giantswarm/architect/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/giantswarm/architect/compare/v1.0.0...v1.1.0

[1.0.0]: https://github.com/giantswarm/architect/releases/tag/v1.0.0
