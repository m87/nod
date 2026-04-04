# Contributing to nod

Thank you for your interest in contributing!

## Getting Started

1. Fork and clone the repository
2. Install Go (see `go.mod` for the required version)
3. Run `go test ./...` to verify everything works

## Running Tests

**SQLite** (no setup required):

```bash
go test ./sqlite -v
```

**PostgreSQL** (requires a running instance):

```bash
export NOD_TEST_POSTGRES_DSN='host=localhost port=5432 user=nod password=nod dbname=nod_test sslmode=disable'
go test ./postgres -v
```

## Adding a New Database Adapter

1. Create a new package under the project root (e.g., `mysql/`)
2. Implement a `NewRepository(dsn, log, mappers)` constructor that returns `*nod.Repository`
3. Re-use `nod.NewRepository(db, log, mappers)` with the GORM dialector for your database
4. Add tests using the contract test suite:

```go
contract.RunRepositoryContractTests(t, func(t *testing.T) *nod.Repository {
    repo, err := NewRepository(dsn, slog.Default(), nod.NewMapperRegistry())
    require.NoError(t, err)
    return repo
})
```

## Code Style

- Run `go vet ./...` before submitting
- All exported types and functions must have godoc comments
- Follow standard Go conventions

## Pull Requests

- Keep PRs focused on a single change
- Include tests for new features
- Update `CHANGELOG.md` under the `[Unreleased]` section
