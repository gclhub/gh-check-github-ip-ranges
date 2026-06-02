---

description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **All Go files at repository root** (single-package monolith)
- Implementation: `main.go`, `ip_checker.go`, or new `feature_name.go`
- Tests: co-located `feature_name_test.go`
- No subdirectories for source code

## Build & Test Commands

```bash
go build ./...      # Compile
go vet ./...        # Static analysis
go test ./...       # Run all tests
go test -v -run TestSpecificFunction ./...  # Run specific test
```

<!-- 
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration purposes only.
  
  The /speckit.tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (with their priorities P1, P2, P3...)
  - Feature requirements from plan.md
  - Entities from data-model.md
  - Endpoints from contracts/
  
  Tasks MUST be organized by user story so each story can be:
  - Implemented independently
  - Tested independently
  - Delivered as an MVP increment
  
  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## Phase 1: Setup & Foundation

**Purpose**: Any prerequisite changes before implementing the feature

- [ ] T001 Add any new dependencies to go.mod (`go get [package]`)
- [ ] T002 [P] Create new Go source file at root: `[feature_name].go` with package declaration
- [ ] T003 [P] Create corresponding test file: `[feature_name]_test.go`

**Checkpoint**: `go build ./...` compiles cleanly

---

## Phase 2: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Test**: `go test -v -run TestFeatureName ./...`

### Tests for User Story 1 (OPTIONAL - only if tests requested) ⚠️

> **NOTE: Write tests FIRST using table-driven style, ensure they FAIL before implementation**

- [ ] T004 [P] [US1] Table-driven test for [function] in `[feature_name]_test.go` using httptest
- [ ] T005 [P] [US1] Test exit code behavior for [scenario] in `main_test.go`

### Implementation for User Story 1

- [ ] T006 [US1] Implement core logic in `[feature_name].go` (types, functions)
- [ ] T007 [US1] Wire into Cobra command in `main.go` (flags, args, invocation)
- [ ] T008 [US1] Add error handling with `fmt.Errorf("context: %w", err)` and correct exit codes
- [ ] T009 [US1] Verify: `go test ./...` passes, `go vet ./...` clean

**Checkpoint**: User Story 1 functional — `go build ./... && go test ./...` green

---

## Phase 3: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: `go test -v -run TestStory2Function ./...`

### Tests for User Story 2 (OPTIONAL - only if tests requested) ⚠️

- [ ] T010 [P] [US2] Table-driven test for [function] in `[feature_name]_test.go`

### Implementation for User Story 2

- [ ] T011 [US2] Implement [feature] in `[feature_name].go` or `ip_checker.go`
- [ ] T012 [US2] Update Cobra command if new flags/args needed in `main.go`
- [ ] T013 [US2] Verify: `go test ./...` passes

**Checkpoint**: User Stories 1 AND 2 both pass independently

---

[Add more user story phases as needed, following the same pattern]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Update README.md with new usage examples
- [ ] TXXX [P] Update ARCHITECTURE.md if structural changes made
- [ ] TXXX Verify all exit codes documented and tested
- [ ] TXXX Final check: `go build ./... && go vet ./... && go test ./...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **User Stories (Phase 2+)**: Depend on Setup completion; can run sequentially (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Core logic before Cobra wiring
- Implementation before error handling polish
- Story complete before moving to next priority

### Parallel Opportunities

- All tasks marked [P] can run in parallel (different files, no conflicts)
- Test files and implementation files can be created in parallel
- Independent user stories can proceed in parallel if no shared state

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (`go build ./...` compiles)
2. Complete Phase 2: User Story 1
3. **STOP and VALIDATE**: `go test ./...` green
4. Commit and verify CI passes

### Incremental Delivery

1. Setup → compiles cleanly
2. User Story 1 → `go test ./...` green → commit (MVP!)
3. User Story 2 → `go test ./...` green → commit
4. Polish → final `go vet ./... && go test ./...` → commit

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- All source files live at repository root (no subdirectories)
- Verify tests fail before implementing (table-driven style)
- Commit after each task or logical group
- Always run `go vet ./...` before final commit
