# gh-check-github-ip-ranges

A GitHub CLI extension to check if an IP address is within GitHub's published IP ranges.

## Installation

```bash
gh extension install gclhub/gh-check-github-ip-ranges
```

## Usage

```bash
gh check-github-ip-ranges <ip-address>
```

### Options

- `-s, --silent`: Silent mode - only use exit codes (useful for scripts)

### Exit Codes

- `0`: Success (IP address belongs to GitHub)
- `1`: IP address does not belong to GitHub
- `2`: Invalid input or error condition:
  - Invalid IP address format
  - Non-IPv4 address (IPv6 is not supported)
  - Private, loopback, multicast, or broadcast IP addresses
  - Network errors when fetching GitHub IP ranges
  - API errors from GitHub's meta endpoint
  - Missing command line arguments

### Examples

Check if an IP address belongs to GitHub:
```bash
gh check-github-ip-ranges 192.30.252.1
```

Use in a script with silent mode:
```bash
if gh check-github-ip-ranges -s 192.30.252.1; then
    echo "IP belongs to GitHub"
else
    echo "IP does not belong to GitHub"
fi
```

## Features

- Validates IP address format and routability
- Checks IPv4 addresses against all GitHub IP ranges
- Returns the specific functional area (Actions, API, Git, etc.) for GitHub IPs
- Includes a silent mode for use in scripts
- Supports all GitHub IP range categories from the /meta API endpoint
- **Security hardened** with timeout protection, input validation, and secure HTTP client configuration

## Security

This tool implements several security measures to ensure safe operation:

### Input Validation
- **Length limits**: IP address inputs are limited to 45 characters to prevent memory exhaustion
- **Empty input protection**: Rejects empty or whitespace-only inputs
- **Format validation**: Strict IP address format checking to prevent injection attacks

### HTTP Client Security  
- **Connection timeouts**: 30-second timeout prevents indefinite hangs
- **TLS requirements**: Minimum TLS 1.2 for all HTTPS connections
- **Certificate validation**: Strict certificate validation prevents MITM attacks
- **Response size limits**: API responses limited to 1MB to prevent memory exhaustion

### Network Security
- **Public IP validation**: Only accepts public, routable IPv4 addresses
- **API rate limiting**: Built-in protection against excessive API calls
- **Secure defaults**: No insecure fallbacks or bypass mechanisms
- **Concurrent safety**: Thread-safe operations with proper mutex synchronization

For security issues, please report them through GitHub's security advisory system.

## Requirements

- GitHub CLI (`gh`) version 2.0.0 or higher
- Go 1.16 or higher (for development)

## Development

To work on this locally:

1. Clone the repository:
```bash
git clone https://github.com/gclhub/gh-check-github-ip-ranges
cd gh-check-github-ip-ranges
```

2. Install the extension from your local directory:
```bash
gh extension remove check-github-ip-ranges  # Remove any existing installation
gh extension install .  # Install from current directory
```

3. Build the extension:
```bash
go build
```

4. Run tests:
```bash
go test ./...
```

Now you can run the extension through `gh` and any changes you make will be reflected after rebuilding:
```bash
gh check-github-ip-ranges <ip-address>
```

To iterate on changes:
1. Make your code changes
2. Run `go build`
3. Test the extension with `gh check-github-ip-ranges`

To enable debug logging, set the `GH_DEBUG` environment variable:
```bash
GH_DEBUG=1 gh check-github-ip-ranges <ip-address>
```

## License

MIT