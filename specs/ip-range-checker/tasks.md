# Tasks: GitHub IP Range Checker

**Input**: Reverse-engineered from existing implementation
**Prerequisites**: spec.md, plan.md (both migrated)
**Status**: Migrated — all tasks reflect completed work

## Path Conventions

- **All Go files at repository root** (single-package monolith)
- Implementation: `main.go`, `ip_checker.go`
- Tests: `main_test.go`, `ip_checker_test.go`

## Build & Test Commands

```bash
go build ./...      # Compile
go vet ./...        # Static analysis
go test ./...       # Run all tests
```

---

## Phase 1: Setup & Foundation

**Purpose**: Project scaffolding and dependency management

- [x] T001 [P] Initialize Go module: `go mod init github.com/gclhub/gh-check-github-ip-ranges`
- [x] T002 [P] Add Cobra dependency: `go get github.com/spf13/cobra`
- [x] T003 [P] Create `main.go` with package declaration and Cobra command skeleton
- [x] T004 [P] Create `ip_checker.go` with package declaration and type stubs

**Checkpoint**: `go build ./...` compiles cleanly ✅

---

## Phase 2: User Story 1 — IP Range Checking (Priority: P1) 🎯 MVP

**Goal**: Determine if an IPv4 address belongs to GitHub's IP ranges and report the matching service category.

**Independent Test**: `go test -v -run TestIPChecker_CheckIP ./...`

### Tests for User Story 1 ✅

- [x] T005 [P] [US1] Table-driven test for `CheckIP()` in `ip_checker_test.go` — valid GitHub IP returns correct category
- [x] T006 [P] [US1] Table-driven test for `CheckIP()` — non-GitHub IP returns `IsGitHubIP: false`
- [x] T007 [P] [US1] Mock server with `httptest.NewServer` for `/meta` response
- [x] T008 [P] [US1] Table-driven test for invalid CIDR entries (silently skipped)

### Implementation for User Story 1

- [x] T009 [US1] Define `GitHubMeta` struct with JSON tags for all 10 service categories in `ip_checker.go`
- [x] T010 [US1] Define `CheckResult` struct (`IsGitHubIP`, `FunctionalArea`, `Range`) in `ip_checker.go`
- [x] T011 [US1] Implement `fetchGitHubMeta()` — HTTP GET, status check, JSON decode in `ip_checker.go`
- [x] T012 [US1] Implement `CheckIP()` — validate IP, fetch meta, iterate categories, check CIDR containment in `ip_checker.go`
- [x] T013 [US1] Wire `CheckIP()` into Cobra `RunE` handler in `main.go`

**Checkpoint**: Core IP checking works end-to-end ✅

---

## Phase 3: User Story 2 — Silent Mode (Priority: P2)

**Goal**: Support `--silent` flag for exit-code-only operation in scripts.

**Independent Test**: `go test -v -run TestRunCommand ./...` (silent cases)

### Tests for User Story 2 ✅

- [x] T014 [P] [US2] Test silent mode produces no stdout for GitHub IP in `main_test.go`
- [x] T015 [P] [US2] Test silent mode produces no stderr for non-GitHub IP in `main_test.go`

### Implementation for User Story 2

- [x] T016 [US2] Add `--silent` / `-s` bool flag to Cobra command in `main.go`
- [x] T017 [US2] Gate stdout output on `!silent` condition in `runCommand()` in `main.go`
- [x] T018 [US2] Gate stderr output on `!cmd.Flags().Changed("silent")` in `main()` error handler

**Checkpoint**: Silent mode passes all tests ✅

---

## Phase 4: User Story 3 — Input Validation (Priority: P2)

**Goal**: Reject invalid, private, IPv6, and special addresses with clear error messages before making API calls.

**Independent Test**: `go test -v -run TestIPChecker_CheckIP ./...` (validation cases)

### Tests for User Story 3 ✅

- [x] T019 [P] [US3] Test invalid IP format returns error with "invalid IP address format"
- [x] T020 [P] [US3] Test IPv6 address returns error with "only IPv4 addresses are supported"
- [x] T021 [P] [US3] Test private IP (192.168.x.x) returns error with "public, routable address"
- [x] T022 [P] [US3] Test broadcast address (255.255.255.255) returns error

### Implementation for User Story 3

- [x] T023 [US3] Implement IP parsing with `net.ParseIP()` — reject nil (malformed) in `ip_checker.go`
- [x] T024 [US3] Implement IPv4 check with `.To4()` — reject nil (IPv6) in `ip_checker.go`
- [x] T025 [US3] Implement public address check: reject `.IsPrivate()`, `.IsLoopback()`, `.IsUnspecified()`, `.IsMulticast()` in `ip_checker.go`
- [x] T026 [US3] Implement `isBroadcastAddress()` helper — reject 255.255.255.255 in `ip_checker.go`

**Checkpoint**: All validation tests pass ✅

---

## Phase 5: User Story 4 — API Error Handling (Priority: P3)

**Goal**: Gracefully handle network failures, non-200 responses, and malformed JSON from GitHub's API.

**Independent Test**: `go test -v -run TestIPChecker_CheckIP ./...` (error server cases)

### Tests for User Story 4 ✅

- [x] T027 [P] [US4] Test HTTP 500 response triggers error in `ip_checker_test.go`
- [x] T028 [P] [US4] Test invalid JSON response triggers decode error in `ip_checker_test.go`
- [x] T029 [P] [US4] Test network failure (failingTransport) triggers fetch error in `ip_checker_test.go`

### Implementation for User Story 4

- [x] T030 [US4] Check `resp.StatusCode != http.StatusOK` → return error with status code in `ip_checker.go`
- [x] T031 [US4] Wrap JSON decode errors with `fmt.Errorf("failed to decode GitHub meta response: %w", err)` in `ip_checker.go`
- [x] T032 [US4] Wrap HTTP client errors with `fmt.Errorf("failed to fetch GitHub meta: %w", err)` in `ip_checker.go`

**Checkpoint**: All API error scenarios handled gracefully ✅

---

## Phase 6: Polish & Integration

**Purpose**: CLI-level integration, exit code routing, and documentation

- [x] T033 [P] Implement exit code routing in `main()`: code 1 for "not found", code 2 for all other errors
- [x] T034 [P] Implement stderr message formatting: no `Error:` prefix for "not found", `Error:` prefix for real errors
- [x] T035 [P] Add `--version` support via Cobra `Version` field
- [x] T036 [P] Integration tests for full CLI lifecycle (`TestMainFunction`) with exit code capture in `main_test.go`
- [x] T037 [P] Create `README.md` with installation and usage instructions
- [x] T038 [P] Create `ARCHITECTURE.md` documenting project structure
- [x] T039 [P] Create `TESTING.md` documenting test strategy
- [x] T040 Final verification: `go build ./... && go vet ./... && go test ./...` all pass

**Checkpoint**: Full CLI operational, all tests green, docs complete ✅

---

## Dependencies & Execution Order

### Phase Dependencies (as built)

- **Phase 1 (Setup)**: No dependencies — foundation
- **Phase 2 (US1 - Core)**: Depends on Phase 1; provides MVP
- **Phase 3 (US2 - Silent)**: Depends on Phase 2 (needs working CLI)
- **Phase 4 (US3 - Validation)**: Independent of Phases 3; depends on Phase 2
- **Phase 5 (US4 - Errors)**: Depends on Phase 2 (needs API integration)
- **Phase 6 (Polish)**: Depends on all previous phases

### Parallel Opportunities

- Phases 3, 4, 5 could have been developed in parallel (independent concerns)
- All test files and implementation files within each phase are parallelizable

---

## Identified Gaps

| ID | Gap | Severity | Recommendation |
|----|-----|----------|----------------|
| G001 | No HTTP client timeout | Medium | Add `http.Client{Timeout: 10*time.Second}` to `NewIPChecker()` |
| G002 | No rate limit detection (403) | Low | Check for 403 and suggest `GITHUB_TOKEN` usage |
| G003 | No response caching | Low | Consider local file cache with 1hr TTL for repeated use |
| G004 | No IPv6 support | Low | Future enhancement; `/meta` provides IPv6 ranges |
| G005 | No authenticated API requests | Low | Support `GITHUB_TOKEN` env var for higher rate limits |

---

## Summary

| Metric | Value |
|--------|-------|
| Total tasks | 40 |
| Completed | 40/40 (100%) |
| User stories covered | 4/4 |
| Test coverage | All stories have dedicated test tasks |
| Gaps identified | 5 |
