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
- `1`: Error or IP address does not belong to GitHub

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

## Requirements

- GitHub CLI (`gh`) version 2.0.0 or higher
- Go 1.16 or higher (for development)

## Development

To work on this locally:

```bash
git clone https://github.com/gclhub/gh-check-github-ip-ranges
cd gh-check-github-ip-ranges
```

2. Build the extension
```bash
go build
```

3. Run tests
```bash
go test ./...
```

## License

MIT