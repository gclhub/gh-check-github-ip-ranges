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
	Hooks        []string `json:"hooks"`
	Web          []string `json:"web"`
	Api          []string `json:"api"`
	Git          []string `json:"git"`
	Packages     []string `json:"packages"`
	Pages        []string `json:"pages"`
	Importer     []string `json:"importer"`
	Actions      []string `json:"actions"`
	Dependabot   []string `json:"dependabot"`
	ActionsIPv4  []string `json:"actions_ipv4"`
}

// IPChecker provides functionality to check IP addresses against GitHub's ranges
type IPChecker struct {
	meta *GitHubMeta
}

// CheckResult contains the result of an IP check
type CheckResult struct {
	IsGitHubIP     bool
	FunctionalArea string
	Range          string
}

// NewIPChecker creates a new IPChecker instance
func NewIPChecker() *IPChecker {
	return &IPChecker{}
}

// fetchGitHubMeta fetches the IP ranges from GitHub's API
func (c *IPChecker) fetchGitHubMeta() error {
	resp, err := http.Get(githubMetaURL)
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

	// Ensure it's an IPv4 address
	ip = ip.To4()
	if ip == nil {
		return nil, fmt.Errorf("only IPv4 addresses are supported")
	}

	// Check if it's a public IP address
	if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || isBroadcastAddress(ip) {
		return nil, fmt.Errorf("IP address must be a public, routable address")
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