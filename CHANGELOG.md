# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- PostgreSQL adapter (`nod/postgres`)
- Typed repository API with `TypedRepository` and `TypedQuery`
- Tree traversal: ancestor and descendant queries with typed mapping
- Key-value, tag, and content support on nodes
- Contract test suite for adapter conformance
- CI with SQLite and PostgreSQL test jobs

### Changed
- Repository fields are now private; use `DB()`, `Log()`, `Mappers()` accessors
- GORM type annotations are now dialect-agnostic (removed `type:datetime`, `type:real`)
- SQL queries use parameterized values for cross-database compatibility
