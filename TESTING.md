# Testing Documentation

## Unit Test Summary

This project uses Go's built-in testing framework to ensure code quality and reliability. The test suite covers all major components and edge cases.

### Test Files

#### `ip_checker_test.go`
Tests for the IP checker functionality that validates and checks IP addresses against GitHub's IP ranges.

**Test Function: `TestIPChecker_CheckIP`**

This comprehensive test validates the core IP checking logic with 11 different test cases:

1. **Valid GitHub IP** - Verifies that a known GitHub IP (192.30.252.1) is correctly identified and returns the appropriate functional area (Hooks) and CIDR range.

2. **Non-GitHub IP** - Tests that a public IP not owned by GitHub (8.8.8.8) is correctly identified as not belonging to GitHub.

3. **Invalid IP Format** - Ensures that malformed IP addresses (e.g., "invalid-ip") return an appropriate error message.

4. **Private IP** - Validates that private IP addresses (e.g., 192.168.1.1) are rejected with an error indicating they must be public, routable addresses.

5. **IPv6 Address** - Confirms that IPv6 addresses (e.g., 2001:db8::1) are properly rejected since only IPv4 is supported.

6. **Broadcast Address** - Tests that broadcast addresses (255.255.255.255) are correctly identified as non-routable and rejected.

7. **API Server Error** - Simulates a GitHub API server error (HTTP 500) to ensure proper error handling when the upstream API is unavailable.

8. **Invalid JSON Response** - Tests resilience when the GitHub API returns malformed JSON data.

9. **Network Error** - Uses a failing HTTP transport to verify proper error handling when network connectivity fails.

10. **Invalid CIDR Format** - Ensures the checker gracefully handles invalid CIDR notation in the API response by skipping invalid entries and continuing to check valid ones.

11. **Valid IP Found After Invalid CIDRs** - Verifies that the checker can successfully match an IP even when some CIDR ranges in the response are malformed.

**Test Infrastructure:**
- Uses `httptest.NewServer` to create mock GitHub API servers
- Tests multiple server scenarios (success, error, invalid JSON, invalid CIDR)
- Implements a custom `failingTransport` to simulate network failures
- Properly overrides the `githubMetaURL` variable for testing

#### `main_test.go`
Tests for the CLI interface and main application logic.

**Test Function: `TestRunCommand`**

Tests the `runCommand` function with 6 test cases covering different CLI scenarios:

1. **GitHub IP with silent mode** - Verifies silent mode produces no output when a GitHub IP is found.

2. **GitHub IP without silent mode** - Confirms proper output message format when a GitHub IP is identified.

3. **Non-GitHub IP with silent mode** - Tests that errors are properly returned without output in silent mode.

4. **Non-GitHub IP without silent mode** - Validates error handling for non-GitHub IPs.

5. **Invalid IP with silent mode** - Tests invalid input handling in silent mode.

6. **Invalid IP without silent mode** - Validates error messages for invalid IP addresses.

**Test Infrastructure:**
- Uses `os.Pipe` to capture stdout and stderr
- Creates Cobra command instances with appropriate flags
- Overrides `githubMetaURL` for testing
- Validates both output content and error conditions

**Test Function: `TestMainFunction`**

Integration tests for the complete CLI application with 11 test cases:

1. **Valid GitHub IP** - End-to-end test with a valid GitHub IP, expecting exit code 0.

2. **Non-GitHub IP** - Tests exit code 1 for IPs not owned by GitHub.

3. **Invalid IP format** - Validates exit code 2 for malformed IP addresses.

4. **Private IP** - Ensures private IPs result in exit code 2.

5. **IPv6 address** - Confirms IPv6 addresses result in exit code 2.

6. **Broadcast address** - Tests broadcast addresses result in exit code 2.

7. **No arguments** - Validates proper error when no IP is provided (exit code 2).

8. **Help flag** - Ensures `--help` works correctly with exit code 0.

9. **Non-GitHub IP with silent mode** - Tests silent mode behavior with exit code 1.

10. **Version flag** - Validates `--version` flag functionality.

11. **Version flag short form** - Tests `-v` shorthand for version.

**Test Infrastructure:**
- Runs `main()` function in a goroutine
- Overrides `osExit` to capture exit codes without terminating tests
- Uses channels to communicate exit codes between goroutines
- Validates both exit codes and stderr output

### Test Coverage Summary

The test suite provides comprehensive coverage of:

- ✅ **IP Validation**: Format validation, IPv4/IPv6 detection, public/private IP checks
- ✅ **GitHub IP Range Checking**: Matching against all functional areas (Hooks, Web, API, Git, Packages, Pages, Importer, Actions, Dependabot, Actions IPv4)
- ✅ **Error Handling**: Network errors, API errors, invalid data, malformed input
- ✅ **CLI Interface**: Command-line argument parsing, flags, output formatting
- ✅ **Exit Codes**: Proper exit codes for success, not found, and error conditions
- ✅ **Silent Mode**: No output behavior when silent flag is used
- ✅ **Edge Cases**: Invalid CIDRs, broadcast addresses, private IPs, empty input

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate detailed coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Design Principles

1. **Isolation**: Each test uses mock HTTP servers to avoid external dependencies
2. **Deterministic**: Tests produce consistent results regardless of external state
3. **Comprehensive**: Cover happy paths, error conditions, and edge cases
4. **Fast**: All tests use in-memory mocks and complete in milliseconds
5. **Maintainable**: Clear test names and well-structured test data
