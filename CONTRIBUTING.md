# Contributing

dbui is a small Go project — contributions are welcome.

## Setup

- Go (version from `go.mod`); Docker for integration tests only.
- `make test` — unit tests.
- `make test-integration` — driver + acceptance suites (Docker).
- `make lint` — golangci-lint (v2 config).
- `make snapshot` — cross-platform release build via GoReleaser.

## Pull requests

- Keep PRs small and single-purpose; reference the issue they close.
- Use Conventional Commit messages (`fix:`, `feat:`, `test:`, `ci:`, `docs:`, `chore:`).
- CI must be green: lint, unit, integration, govulncheck, snapshot.
- The cross-engine suite in `internal/acceptance` is the contract every
  `DataSource` change must keep green.

## Roadmap

Work is organized by [milestones](https://github.com/KenanBek/dbui/milestones) —
pick an issue from the current one.
