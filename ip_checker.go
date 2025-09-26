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

// CheckIP checks if the provided IP address is within GitHub's ranges
func (c *IPChecker) CheckIP(ipStr string) (*CheckResult, error) {
	// Parse and validate the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address format")
	}

	// Check if it's a public IP address (both IPv4 and IPv6)
	if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() {
		return nil, fmt.Errorf("IP address must be a public, routable address")
	}

	// Additional checks for IPv4
	if ipv4 := ip.To4(); ipv4 != nil {
		if isBroadcastAddress(ipv4) {
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
	ranges := []struct {
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
		{"Hooks IPv6", c.meta.HooksIPv6},
		{"Web IPv6", c.meta.WebIPv6},
		{"API IPv6", c.meta.ApiIPv6},
		{"Git IPv6", c.meta.GitIPv6},
		{"Packages IPv6", c.meta.PackagesIPv6},
		{"Pages IPv6", c.meta.PagesIPv6},
		{"Importer IPv6", c.meta.ImporterIPv6},
		{"Actions IPv6", c.meta.ActionsIPv6},
		{"Dependabot IPv6", c.meta.DependabotIPv6},
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
