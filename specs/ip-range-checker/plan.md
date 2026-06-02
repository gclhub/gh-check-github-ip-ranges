# Implementation Plan: GitHub IP Range Checker

**Branch**: `main` (pre-existing) | **Date**: 2025-07-21 | **Spec**: `specs/ip-range-checker/spec.md`
**Input**: Reverse-engineered from existing implementation
**Status**: Migrated

## Summary

A GitHub CLI extension that checks whether a given IPv4 address belongs to GitHub's published IP ranges by fetching the `/meta` API endpoint, parsing CIDR ranges by service category, and reporting the matching category. Implemented as a single-package Go CLI using Cobra for argument parsing.

## Technical Context

**Language/Version**: Go 1.24.2  
**Primary Dependencies**: github.com/spf13/cobra v1.9.1 (CLI framework)  
**Storage**: N/A (stateless — fetches from GitHub /meta API on each invocation)  
**Testing**: `go test ./...` with standard library `testing` package, `net/http/httptest` for mocks  
**Target Platform**: Cross-platform (Linux, macOS, Windows) via gh-extension-precompile  
**Project Type**: CLI extension (GitHub CLI plugin)  
**Performance Goals**: Single API call + CIDR check; < 2 seconds at p95 on standard broadband  
**Constraints**: No CGo (pure Go for cross-compilation); minimal dependencies  
**Scale/Scope**: Single-package monolith; all `.go` files at repository root

## Constitution Check

*GATE: Verified against `.specify/memory/constitution.md`*

| Rule | Status |
|------|--------|
| Single-package at root | ✅ All code in `main` package at root |
| Semantic exit codes (0/1/2) | ✅ Implemented correctly |
| Table-driven tests | ✅ Both test files use `[]struct{}` pattern |
| DI via package vars | ✅ `osExit`, `githubMetaURL`, `setClient()` |
| Error wrapping with `%w` | ✅ All errors use `fmt.Errorf("context: %w", err)` |
| No external test frameworks | ✅ Standard `testing` package only |
| `go vet` clean | ✅ Passes static analysis |

## Project Structure

### Documentation (this feature)

```text
specs/ip-range-checker/
├── spec.md              # Feature specification (migrated)
├── plan.md              # This file (migrated)
└── tasks.md             # Task list (migrated, all complete)
```

### Source Code (repository root)

```text
main.go                 # Cobra command setup, entry point, error routing, exit codes (70 lines)
ip_checker.go           # Core logic: IP validation, API fetch, CIDR matching (143 lines)
ip_checker_test.go      # Table-driven tests for IPChecker with httptest mocks (245 lines)
main_test.go            # Integration tests for CLI behavior, exit codes, stdout/stderr (319 lines)
go.mod                  # Module: github.com/gclhub/gh-check-github-ip-ranges
go.sum                  # Dependency checksums
```

**Structure Decision**: Single-package monolith at repository root. No subdirectories for source code. Test files co-located with implementation files (`*_test.go` convention).

## Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| CLI framework | Cobra | Industry standard for Go CLIs; provides flags, args validation, help/version for free |
| HTTP mocking | `net/http/httptest` | Standard library; no external deps; test server pattern |
| Exit code testing | `osExit` package var | Allows tests to capture exit without `os.Exit` killing test process |
| URL override | `githubMetaURL` package var | Simple injection point for test servers |
| IP parsing | `net.ParseIP` + `net.ParseCIDR` | Standard library; well-tested; handles edge cases |
| Error output | stderr for errors, stdout for success | Unix convention; enables piping and scripting |

## Complexity Assessment

| Metric | Value |
|--------|-------|
| Total source lines | 213 |
| Total test lines | 564 |
| Test-to-code ratio | 2.6:1 |
| External dependencies | 1 (Cobra) |
| API integrations | 1 (GitHub /meta) |
| Service categories checked | 10 |
| Cyclomatic complexity | Low (linear flow with validation gates) |

## Architecture

```
User Input (IPv4 string)
    │
    ▼
┌─────────────────┐
│  main.go        │  Cobra command setup, flag parsing
│  runCommand()   │  Orchestrates check, routes output
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  ip_checker.go  │  
│  CheckIP()      │  1. Validate IP (format, IPv4, public)
│                 │  2. Fetch /meta (if not cached)
│                 │  3. Iterate categories, check CIDR containment
│                 │  4. Return CheckResult
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  main.go        │  Format output, determine exit code
│  main()         │  Route to stdout/stderr based on result
└─────────────────┘
```

## Identified Gaps & Future Improvements

| Gap | Impact | Suggested Fix |
|-----|--------|---------------|
| No HTTP timeout | Could hang indefinitely on slow networks | Set `http.Client{Timeout: 10*time.Second}` |
| No rate limit detection | Generic error for 403 responses | Check status 403 and suggest auth token |
| No response caching | Fresh API call every invocation | Optional local cache with TTL |
| No IPv6 support | Users with IPv6 get rejection | Parse IPv6 and check IPv6 ranges from /meta |
| No authenticated requests | Limited to 60 req/hr | Support `GITHUB_TOKEN` env var for auth |
