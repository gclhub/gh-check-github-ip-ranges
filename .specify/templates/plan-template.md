# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

**Language/Version**: Go 1.24.2  
**Primary Dependencies**: github.com/spf13/cobra (CLI framework)  
**Storage**: N/A (stateless — fetches from GitHub /meta API on each invocation)  
**Testing**: `go test ./...` with standard library `testing` package, `net/http/httptest` for mocks  
**Target Platform**: Cross-platform (Linux, macOS, Windows) via gh-extension-precompile  
**Project Type**: CLI extension (GitHub CLI plugin)  
**Performance Goals**: Single API call + CIDR check; sub-second response expected  
**Constraints**: No CGo (pure Go for cross-compilation); minimal dependencies  
**Scale/Scope**: Single-package monolith; all `.go` files at repository root

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Flat single-package layout (all files at root)
main.go                 # Cobra command setup, entry point
ip_checker.go           # Core IP checking logic (fetch /meta, parse CIDR, check containment)
ip_checker_test.go      # Table-driven tests for IP checker with httptest mocks
main_test.go            # Integration-level tests for CLI behavior and exit codes
go.mod                  # Module definition and dependencies
go.sum                  # Dependency checksums
```

**Structure Decision**: Single-package monolith at repository root. No subdirectories for source code. Test files co-located with implementation files (`*_test.go` convention).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

### Go CLI Complexity Factors

Consider these when estimating effort:

- **CIDR parsing edge cases**: IPv4 vs IPv6, malformed ranges, overlapping subnets
- **API reliability**: GitHub /meta endpoint availability, rate limiting, response format changes
- **Cross-platform**: Path handling, network behavior differences across OS
- **Exit code contract**: Must preserve semantic exit codes — breaking change if altered
- **Testability**: New features must remain testable via httptest mocks and osExit capture
