# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`  
**Created**: [DATE]  
**Status**: Draft  
**Input**: User description: "$ARGUMENTS"

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - [Brief Title] (Priority: P1)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently - e.g., "Can be fully tested by [specific action] and delivers [specific value]"]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## CLI Interface *(mandatory for CLI features)*

### Commands & Flags

| Command/Flag | Type | Required | Description |
|---|---|---|---|
| `[positional arg]` | string | yes/no | [what it represents] |
| `--flag-name` | string/bool/int | no | [what it controls] |

### Exit Codes

| Code | Meaning | When |
|---|---|---|
| 0 | Success / IP found | [condition] |
| 1 | Not found | [condition] |
| 2 | Error | [condition — API failure, invalid input, etc.] |

### Example Usage

```bash
gh check-github-ip-ranges [args] [--flags]
```

## GitHub API Contract *(if feature touches /meta endpoint)*

- **Endpoint**: `https://api.github.com/meta`
- **Rate Limits**: Unauthenticated = 60/hr; Authenticated = 5000/hr
- **Response Format**: JSON with CIDR ranges keyed by service category
- **Error Scenarios**: Rate limited (403), endpoint unavailable (5xx)
- **Caching**: [describe caching strategy if applicable]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST [specific capability]
- **FR-002**: CLI MUST [specific capability]
- **FR-003**: CLI MUST [error handling behavior]

*Mark unclear requirements:*

- **FR-004**: CLI MUST [NEEDS CLARIFICATION: detail not specified]

### Error Handling

- All network errors → stderr message + exit code 2
- Invalid IP input → stderr message + exit code 2
- IP not in any range → informative stdout message + exit code 1
- [Additional error scenarios specific to feature]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: [Measurable metric, e.g., "Users can complete account creation in under 2 minutes"]
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
- **SC-003**: [User satisfaction metric, e.g., "90% of users successfully complete primary task on first attempt"]
- **SC-004**: [Business metric, e.g., "Reduce support tickets related to [X] by 50%"]

## Assumptions

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right assumptions based on reasonable defaults
  chosen when the feature description did not specify certain details.
-->

- [Assumption about target users, e.g., "Users have stable internet connectivity"]
- [Assumption about scope boundaries, e.g., "Mobile support is out of scope for v1"]
- [Assumption about data/environment, e.g., "Existing authentication system will be reused"]
- [Dependency on existing system/service, e.g., "Requires access to the existing user profile API"]
