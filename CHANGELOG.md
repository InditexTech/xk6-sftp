# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Adopt the InditexTech CI governance Go profile: verification, CodeQL,
  Scorecard, SonarCloud and the release lane are now generated and kept in sync
  centrally.
- Install Go tooling into a repository-local `.bin/` directory so `make verify`
  resolves `xk6` and `golangci-lint` under any Go toolchain setup.
- Scope Dependabot to product-owned dependencies and keep its pull requests out
  of the release lane.

### Removed

- `PR-verify.yml` and `release.yml`, superseded by the governed workflows.

## [1.0.0] - 2025-04-08

### Added

- Initial release of the `xk6-sftp` k6 extension, providing `newClient`,
  `uploadFile`, `downloadFile`, `deleteFile` and `close`.

[Unreleased]: https://github.com/InditexTech/xk6-sftp/compare/1.0.0...HEAD

[1.0.0]: https://github.com/InditexTech/xk6-sftp/releases/tag/1.0.0
