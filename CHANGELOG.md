# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/architect/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/giantswarm/architect/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/giantswarm/architect/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/giantswarm/architect/compare/v1.1.2...v1.2.0
[1.1.2]: https://github.com/giantswarm/architect/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/giantswarm/architect/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/giantswarm/architect/compare/v1.0.0...v1.1.0

[1.0.0]: https://github.com/giantswarm/architect/releases/tag/v1.0.0
