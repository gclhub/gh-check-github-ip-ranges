# Feature Specification: GitHub IP Range Checker

**Feature Branch**: `main` (pre-existing)  
**Created**: 2025-07-21  
**Status**: Migrated  
**Input**: Reverse-engineered from existing codebase (`main.go`, `ip_checker.go`)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Check if an IP belongs to GitHub (Priority: P1)

A DevOps engineer or security analyst wants to determine whether a given IPv4 address belongs to GitHub's published IP ranges. They run the CLI with a single IP argument and receive a clear answer indicating which GitHub service category (Hooks, Actions, API, etc.) owns that address.

**Why this priority**: This is the core value proposition — without this, the tool has no purpose.

**Independent Test**: Run `gh check-github-ip-ranges 192.30.252.1` and verify it reports the matching service category and CIDR range.

**Acceptance Scenarios**:

1. **Given** a valid public IPv4 address that is in GitHub's ranges, **When** the user runs `gh check-github-ip-ranges <ip>`, **Then** stdout displays `IP <ip> belongs to GitHub's <category> range (<cidr>)` and exits with code 0.
2. **Given** a valid public IPv4 address that is NOT in GitHub's ranges, **When** the user runs `gh check-github-ip-ranges <ip>`, **Then** stderr displays `the provided IP address is not a GitHub-owned address` and exits with code 1.

---

### User Story 2 - Silent mode for scripting (Priority: P2)

A CI/CD pipeline or shell script needs to check IP ownership programmatically using only exit codes, without parsing stdout/stderr output.

**Why this priority**: Exit-code-only operation enables automation use cases — a natural extension of the core check.

**Independent Test**: Run `gh check-github-ip-ranges -s 192.30.252.1` and verify no output is produced but exit code is 0.

**Acceptance Scenarios**:

1. **Given** a GitHub IP and `--silent` flag, **When** the command executes, **Then** no output is produced and exit code is 0.
2. **Given** a non-GitHub IP and `--silent` flag, **When** the command executes, **Then** no output is produced and exit code is 1.
3. **Given** an invalid IP and `--silent` flag, **When** the command executes, **Then** no output is produced and exit code is 2.

---

### User Story 3 - Input validation with clear errors (Priority: P2)

A user provides an invalid, private, or IPv6 address. The tool should reject it immediately with a clear error message (before making any API calls) and exit with code 2.

**Why this priority**: Good error messages prevent user confusion and unnecessary API calls.

**Independent Test**: Run `gh check-github-ip-ranges 192.168.1.1` and verify the error message and exit code 2.

**Acceptance Scenarios**:

1. **Given** a malformed string (not an IP), **When** the command runs, **Then** stderr shows `Error: invalid IP address format` and exit code is 2.
2. **Given** an IPv6 address, **When** the command runs, **Then** stderr shows `Error: only IPv4 addresses are supported` and exit code is 2.
3. **Given** a private/loopback/broadcast IP, **When** the command runs, **Then** stderr shows `Error: IP address must be a public, routable address` and exit code is 2.

---

### User Story 4 - Graceful API failure handling (Priority: P3)

When the GitHub `/meta` API is unreachable, returns an error status, or returns invalid data, the tool exits with code 2 and a descriptive error message rather than panicking or hanging.

**Why this priority**: Network failures are inevitable; graceful degradation builds trust.

**Independent Test**: Run with network disconnected or against a mock returning 500, verify exit code 2 with error context.

**Acceptance Scenarios**:

1. **Given** the GitHub API returns HTTP 500, **When** the command runs with a valid IP, **Then** stderr shows an error containing "GitHub API returned status code 500" and exit code is 2.
2. **Given** the GitHub API returns invalid JSON, **When** the command runs, **Then** stderr shows an error about decoding failure and exit code is 2.
3. **Given** a network timeout/connection failure, **When** the command runs, **Then** stderr shows an error containing "failed to fetch GitHub meta" and exit code is 2.

---

### Edge Cases

- Invalid CIDR entries in the API response are silently skipped (not treated as errors)
- The first matching category is returned (deterministic order: Hooks → Web → API → Git → Packages → Pages → Importer → Actions → Dependabot → Actions IPv4)
- No arguments provided → Cobra prints usage error and exits with code 2
- `--help` and `--version` flags work and exit with code 0

## CLI Interface *(mandatory for CLI features)*

### Commands & Flags

| Command/Flag | Type | Required | Description |
|---|---|---|---|
| `<ip-address>` | string (positional) | yes | IPv4 address to check against GitHub's ranges |
| `--silent`, `-s` | bool | no | Suppress all output; communicate only via exit codes |
| `--version`, `-v` | bool | no | Print version and exit |
| `--help`, `-h` | bool | no | Print usage information and exit |

### Exit Codes

| Code | Meaning | When |
|---|---|---|
| 0 | IP found in GitHub ranges | IP matches a CIDR in any service category |
| 1 | IP not found | Valid public IPv4 but not in any GitHub range |
| 2 | Error | Invalid input, API failure, network error, or missing arguments |

### Example Usage

```bash
# Check a specific IP
gh check-github-ip-ranges 192.30.252.1
# Output: IP 192.30.252.1 belongs to GitHub's Hooks range (192.30.252.0/22)

# Silent mode for scripting
gh check-github-ip-ranges -s 140.82.112.1 && echo "GitHub IP" || echo "Not GitHub"

# Invalid input
gh check-github-ip-ranges 192.168.1.1
# Error: IP address must be a public, routable address
```

## GitHub API Contract *(if feature touches /meta endpoint)*

- **Endpoint**: `https://api.github.com/meta`
- **Rate Limits**: Unauthenticated = 60/hr; Authenticated = 5000/hr
- **Response Format**: JSON with CIDR ranges keyed by service category
- **Parsed Categories**: `hooks`, `web`, `api`, `git`, `packages`, `pages`, `importer`, `actions`, `dependabot`, `actions_ipv4`
- **Error Scenarios**: Rate limited (403), endpoint unavailable (5xx), invalid JSON body
- **Caching**: None — fresh fetch on every invocation (no persistent state)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST accept exactly one positional argument (IPv4 address string)
- **FR-002**: CLI MUST validate the IP is a well-formed IPv4 address before making API calls
- **FR-003**: CLI MUST reject private, loopback, multicast, broadcast, and unspecified addresses
- **FR-004**: CLI MUST reject IPv6 addresses with a clear error message
- **FR-005**: CLI MUST fetch GitHub's `/meta` API endpoint and parse CIDR ranges by service category
- **FR-006**: CLI MUST check the IP against all 10 service categories in deterministic order: Hooks → Web → API → Git → Packages → Pages → Importer → Actions → Dependabot → Actions IPv4
- **FR-007**: CLI MUST skip invalid CIDR entries silently (no error, continue checking)
- **FR-008**: CLI MUST report the first matching category name and CIDR range on stdout
- **FR-009**: CLI MUST support `--silent` / `-s` flag to suppress all output
- **FR-010**: CLI MUST use semantic exit codes: 0 (found), 1 (not found), 2 (error)
- **FR-011**: CLI MUST print errors to stderr (not stdout)
- **FR-012**: CLI MUST distinguish "not found" messages (no `Error:` prefix) from actual errors (`Error:` prefix)
- **FR-013**: CLI MUST support `--version` and `--help` flags via Cobra

### Error Handling

- All network errors → stderr message + exit code 2
- Invalid IP input → stderr message + exit code 2
- IP not in any range → informative stderr message (no `Error:` prefix) + exit code 1
- API non-200 status → stderr message with status code + exit code 2
- Invalid JSON response → stderr message about decode failure + exit code 2

### Key Entities

- **GitHubMeta**: Struct representing the `/meta` API response; contains 10 string slices of CIDR ranges keyed by service category
- **IPChecker**: Service object that fetches the meta response and checks IP containment (stateless — no persistent cache)
- **CheckResult**: Result struct with `IsGitHubIP` bool, `FunctionalArea` string, and `Range` string

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can determine GitHub IP ownership in a single command invocation (< 2 seconds at p95 on standard broadband)
- **SC-002**: Exit codes are reliable enough for use in CI/CD pipelines and shell conditionals
- **SC-003**: All 10 GitHub service categories are checked (no partial coverage)
- **SC-004**: Invalid inputs are rejected before any network call is made (fast failure)
- **SC-005**: Test suite achieves full coverage of validation logic, API interaction, and CLI behavior

## Assumptions

- Users have internet connectivity to reach `api.github.com`
- GitHub's `/meta` endpoint remains publicly accessible without authentication
- IPv6 support is explicitly out of scope for v1.0.0
- The tool is used as a `gh` CLI extension (installed via `gh extension install`)
- Rate limiting is unlikely for typical usage (< 60 checks/hour per user)
- Response caching across invocations is not needed (stateless design preferred)
