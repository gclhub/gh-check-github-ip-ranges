# gh-check-github-ip-ranges Constitution

## Core Principles

### I. Single-Package CLI Extension

All Go source files live at the repository root in a single `main` package. No nested packages, no `internal/`, no `pkg/`. The project is a GitHub CLI extension (`gh` plugin) that checks whether an IP address belongs to GitHub's published IP ranges.

### II. CLI Interface Contract

- Arguments via positional params and Cobra flags; results to stdout; errors to stderr
- **Exit codes are semantic**: `0` = IP found in ranges, `1` = IP not found, `2` = runtime error
- Exit codes are part of the public contract — never change their meaning without a major version bump
- Use `osExit` package-level variable (not `os.Exit` directly) to enable testability

### III. Table-Driven Testing

- All tests use Go table-driven style (`tests := []struct{...}; for _, tt := range tests`)
- HTTP interactions mocked via `net/http/httptest`
- Test files co-located: `foo.go` → `foo_test.go`
- Run tests: `go test ./...`
- No external test frameworks — standard library `testing` package only

### IV. Dependency Injection via Package Vars

- HTTP clients injected via `setClient()` helper for test substitution
- `osExit` variable wraps `os.Exit` so tests can capture exit behavior
- No interfaces-for-the-sake-of-interfaces; prefer simple function/var injection

### V. Error Handling

- Wrap errors with `fmt.Errorf("context: %w", err)` for chain tracing
- Errors that reach the user print to stderr and exit with code `2`
- Never swallow errors silently — log or propagate

### VI. Simplicity & YAGNI

- No abstractions until a second use case demands them
- No ORM, no config files, no dependency injection frameworks
- External data source is GitHub's `/meta` API endpoint — fetch, parse, check

## Technical Standards

### Language & Toolchain

- **Go 1.24.2** (see `go.mod`)
- **Cobra** (`github.com/spf13/cobra`) for CLI argument parsing
- **No CGo** — pure Go for cross-compilation via `gh-extension-precompile`

### Naming Conventions

- **Files**: `snake_case.go` (e.g., `ip_checker.go`, `ip_checker_test.go`)
- **Exported identifiers**: `PascalCase` for types and functions, `camelCase` for unexported
- **Test functions**: `TestFunctionName` or `TestFunctionName_scenario`

### API Contract

- GitHub `/meta` endpoint returns JSON with CIDR ranges keyed by service (actions, hooks, etc.)
- Parse CIDR ranges with `net.ParseCIDR`; check containment with `network.Contains(ip)`
- Handle rate limiting gracefully (exit code 2 with informative message)

## Quality Gates

### Build & Test

```bash
go build ./...     # Must compile cleanly
go vet ./...       # Must pass static analysis
go test ./...      # Must pass all tests
```

### CI Pipeline

- GitHub Actions with `cli/gh-extension-precompile@v2` for cross-platform release builds
- All PRs must pass `go test ./...` before merge

### Code Ownership

- `CODEOWNERS` file defines reviewers
- Architecture documented in `ARCHITECTURE.md`
- Testing strategy documented in `TESTING.md`

## Governance

- This constitution reflects conventions detected in the existing codebase
- Amendments require updating this file and verifying no existing code violates new rules
- When in doubt, follow existing patterns in the codebase over external Go conventions

**Version**: 1.0.0 | **Ratified**: 2025-07-21 | **Last Amended**: 2025-07-21
