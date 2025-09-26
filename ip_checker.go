package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

var githubMetaURL = "https://api.github.com/meta"

// GitHubMeta represents the response from GitHub's /meta API endpoint
type GitHubMeta struct {
	Hooks          []string `json:"hooks"`
	Web            []string `json:"web"`
	Api            []string `json:"api"`
	Git            []string `json:"git"`
	Packages       []string `json:"packages"`
	Pages          []string `json:"pages"`
	Importer       []string `json:"importer"`
	Actions        []string `json:"actions"`
	Dependabot     []string `json:"dependabot"`
	ActionsIPv4    []string `json:"actions_ipv4"`
	// IPv6 ranges
	HooksIPv6      []string `json:"hooks_ipv6"`
	WebIPv6        []string `json:"web_ipv6"`
	ApiIPv6        []string `json:"api_ipv6"`
	GitIPv6        []string `json:"git_ipv6"`
	PackagesIPv6   []string `json:"packages_ipv6"`
	PagesIPv6      []string `json:"pages_ipv6"`
	ImporterIPv6   []string `json:"importer_ipv6"`
	ActionsIPv6    []string `json:"actions_ipv6"`
	DependabotIPv6 []string `json:"dependabot_ipv6"`
}

// IPChecker provides functionality to check IP addresses against GitHub's ranges
type IPChecker struct {
	meta   *GitHubMeta
	client *http.Client // Add client field
}

// CheckResult contains the result of an IP check
type CheckResult struct {
	IsGitHubIP     bool
	FunctionalArea string
	Range          string
}

// NewIPChecker creates a new IPChecker instance
func NewIPChecker() *IPChecker {
	return &IPChecker{
		client: http.DefaultClient,
	}
}

// For testing purposes
func (c *IPChecker) setClient(client *http.Client) {
	c.client = client
}

// fetchGitHubMeta fetches the IP ranges from GitHub's API
func (c *IPChecker) fetchGitHubMeta() error {
	resp, err := c.client.Get(githubMetaURL) // Use injected client
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub meta: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var meta GitHubMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return fmt.Errorf("failed to decode GitHub meta response: %w", err)
	}

	c.meta = &meta
	return nil
}

// isBroadcastAddress checks if the IP is a broadcast address
func isBroadcastAddress(ip net.IP) bool {
	for i := 0; i < len(ip); i++ {
		if ip[i] != 255 {
			return false
		}
	}
	return true
}

// isPublicIPv6 checks if the IPv6 address is public and routable
func isPublicIPv6(ip net.IP) bool {
	// Check for common non-public IPv6 addresses
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() {
		return false
	}
	
	// Check for link-local addresses (fe80::/10)
	if len(ip) == 16 && ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
		return false
	}
	
	// Check for unique local addresses (fc00::/7)
	if len(ip) == 16 && (ip[0]&0xfe) == 0xfc {
		return false
	}
	
	// Check for documentation addresses (2001:db8::/32)
	if len(ip) == 16 && ip[0] == 0x20 && ip[1] == 0x01 &&
		ip[2] == 0x0d && ip[3] == 0xb8 {
		return false
	}
	
	return true
}

// CheckIP checks if the provided IP address is within GitHub's ranges
func (c *IPChecker) CheckIP(ipStr string) (*CheckResult, error) {
	// Parse and validate the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address format")
	}

	// Determine if it's IPv4 or IPv6
	var isIPv6 bool
	ipv4 := ip.To4()
	if ipv4 != nil {
		// It's IPv4
		ip = ipv4
		isIPv6 = false
		
		// Check if it's a public IPv4 address
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || isBroadcastAddress(ip) {
			return nil, fmt.Errorf("IP address must be a public, routable address")
		}
	} else {
		// It's IPv6
		isIPv6 = true
		
		// Check if it's a public IPv6 address
		if !isPublicIPv6(ip) {
			return nil, fmt.Errorf("IP address must be a public, routable address")
		}
	}

	// Fetch GitHub meta if not already cached
	if c.meta == nil {
		if err := c.fetchGitHubMeta(); err != nil {
			return nil, fmt.Errorf("failed to fetch GitHub meta: %w", err)
		}
	}

	// Check each range category
	var ranges []struct {
		name   string
		ranges []string
	}
	
	if isIPv6 {
		ranges = []struct {
			name   string
			ranges []string
		}{
			{"Hooks", c.meta.HooksIPv6},
			{"Web", c.meta.WebIPv6},
			{"API", c.meta.ApiIPv6},
			{"Git", c.meta.GitIPv6},
			{"Packages", c.meta.PackagesIPv6},
			{"Pages", c.meta.PagesIPv6},
			{"Importer", c.meta.ImporterIPv6},
			{"Actions", c.meta.ActionsIPv6},
			{"Dependabot", c.meta.DependabotIPv6},
		}
	} else {
		ranges = []struct {
			name   string
			ranges []string
		}{
			{"Hooks", c.meta.Hooks},
			{"Web", c.meta.Web},
			{"API", c.meta.Api},
			{"Git", c.meta.Git},
			{"Packages", c.meta.Packages},
			{"Pages", c.meta.Pages},
			{"Importer", c.meta.Importer},
			{"Actions", c.meta.Actions},
			{"Dependabot", c.meta.Dependabot},
			{"Actions IPv4", c.meta.ActionsIPv4},
		}
	}

	for _, category := range ranges {
		for _, cidr := range category.ranges {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}

			if ipNet.Contains(ip) {
				return &CheckResult{
					IsGitHubIP:     true,
					FunctionalArea: category.name,
					Range:          cidr,
				}, nil
			}
		}
	}

	return &CheckResult{IsGitHubIP: false}, nil
}
