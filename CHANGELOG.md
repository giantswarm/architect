# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [3.7.0] - 2021-06-16

### Changed

- Removed decommissioned installations (`archon`, `davis`, `dinosaur`, `dragon`).

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

[Unreleased]: https://github.com/giantswarm/architect/compare/v3.7.0...HEAD
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
