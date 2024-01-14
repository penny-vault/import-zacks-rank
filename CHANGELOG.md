# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added

- Download of balance sheet for a speicific ticker

### Changed

- Updated to reflect latest playwright API

### Deprecated

### Removed

### Fixed

### Security

## [0.2.1] - 2023-07-04
### Security
- Upgrade to latest version of playwright/chromium

## [0.2.0] - 2022-10-09
### Added
- Debugging option that saves a PDF
- Flag to output logs as JSON
- Test command that only downloads and parses data from zacks but does not save to DB or upload to backblaze

### Changed

- Blocked additional ad websites that were interfering with zacks rank download

### Deprecated

### Removed

### Fixed
- Added missing error handlers for some playwright calls

### Security

## [0.1.0] - 2022-05-26
### Added
- Download zacks ranks and other fundamental data from zacks rank
- Upload parquet file to backblaze
- Update pvdb database with downloaded data
- Add Dockerfile for building a container

[Unreleased]: https://github.com/penny-vault/import-zacks-rank/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/penny-vault/import-zacks-rank/releases/tag/v0.2.1
[0.2.0]: https://github.com/penny-vault/import-zacks-rank/releases/tag/v0.2.0
[0.1.0]: https://github.com/penny-vault/import-zacks-rank/releases/tag/v0.1.0
