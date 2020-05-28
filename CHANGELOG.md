# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/architect/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/giantswarm/architect/releases/tag/v1.1.0
[1.0.0]: https://github.com/giantswarm/architect/releases/tag/v1.0.0
