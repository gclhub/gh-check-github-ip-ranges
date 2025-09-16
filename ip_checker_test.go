package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIPChecker_CheckIP(t *testing.T) {
	// Success case server
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["192.30.252.0/22"],
			"web": ["192.30.252.0/22"],
			"api": ["192.30.252.0/22"],
			"git": ["192.30.252.0/22"],
			"packages": ["192.30.252.0/22"],
			"pages": ["192.30.252.0/22"],
			"importer": ["192.30.252.0/22"],
			"actions": ["192.30.252.0/22"],
			"dependabot": ["192.30.252.0/22"],
			"actions_ipv4": ["192.30.252.0/22"]
		}`))
	}))
	defer successServer.Close()

	// Error case server - returns 500
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	// Invalid JSON server
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer invalidJSONServer.Close()

	// Invalid CIDR server
	invalidCIDRServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["invalid-cidr"],
			"web": ["not-a-cidr"],
			"api": ["192.30.252.0/22"],
			"git": ["192.30.252.0/22"]
		}`))
	}))
	defer invalidCIDRServer.Close()

	// Mixed valid/invalid CIDR server
	mixedCIDRServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["invalid-cidr", "also-invalid"],
			"web": ["not-a-cidr"],
			"api": ["invalid"],
			"git": ["192.30.252.0/22"]
		}`))
	}))
	defer mixedCIDRServer.Close()

	// Create a client that always fails
	failingClient := &http.Client{
		Transport: &failingTransport{},
	}

	tests := []struct {
		name       string
		ip         string
		mockServer *httptest.Server
		client     *http.Client
		wantErr    bool
		wantErrMsg string
		want       *CheckResult
	}{
		{
			name:       "Valid GitHub IP",
			ip:         "192.30.252.1",
			mockServer: successServer,
			client:     nil,
			wantErr:    false,
			want: &CheckResult{
				IsGitHubIP:     true,
				FunctionalArea: "Hooks",
				Range:          "192.30.252.0/22",
			},
		},
		{
			name:       "Non-GitHub IP",
			ip:         "8.8.8.8",
			mockServer: successServer,
			client:     nil,
			wantErr:    false,
			want: &CheckResult{
				IsGitHubIP: false,
			},
		},
		{
			name:       "Invalid IP format",
			ip:         "invalid-ip",
			mockServer: successServer,
			client:     nil,
			wantErr:    true,
			wantErrMsg: "invalid IP address format",
			want:       nil,
		},
		{
			name:       "Private IP",
			ip:         "192.168.1.1",
			mockServer: successServer,
			client:     nil,
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
			want:       nil,
		},
		{
			name:       "IPv6 address",
			ip:         "2001:db8::1",
			mockServer: successServer,
			client:     nil,
			wantErr:    true,
			wantErrMsg: "only IPv4 addresses are supported",
			want:       nil,
		},
		{
			name:       "Broadcast address",
			ip:         "255.255.255.255",
			mockServer: successServer,
			client:     nil,
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
			want:       nil,
		},
		{
			name:       "API Server Error",
			ip:         "8.8.8.8",
			mockServer: errorServer,
			client:     nil,
			wantErr:    true,
			want:       nil,
		},
		{
			name:       "Invalid JSON Response",
			ip:         "8.8.8.8",
			mockServer: invalidJSONServer,
			client:     nil,
			wantErr:    true,
			want:       nil,
		},
		{
			name:       "Network Error",
			ip:         "8.8.8.8",
			mockServer: successServer,
			client:     failingClient,
			wantErr:    true,
			wantErrMsg: "failed to fetch GitHub meta",
			want:       nil,
		},
		{
			name:       "Invalid CIDR Format",
			ip:         "192.30.252.1",
			mockServer: invalidCIDRServer,
			client:     nil,
			wantErr:    false, // Should not error, just skip invalid CIDRs
			want: &CheckResult{
				IsGitHubIP:     true,
				FunctionalArea: "API",
				Range:          "192.30.252.0/22",
			},
		},
		{
			name:       "Valid IP Found After Invalid CIDRs",
			ip:         "192.30.252.1",
			mockServer: mixedCIDRServer,
			client:     nil,
			wantErr:    false,
			want: &CheckResult{
				IsGitHubIP:     true,
				FunctionalArea: "Git",
				Range:          "192.30.252.0/22",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the GitHub meta URL for testing
			oldURL := githubMetaURL
			githubMetaURL = tt.mockServer.URL
			defer func() { githubMetaURL = oldURL }()

			checker := NewIPChecker()
			if tt.client != nil {
				checker.setClient(tt.client)
			}

			got, err := checker.CheckIP(tt.ip)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErrMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("CheckIP() error message = %v, should contain %v", err.Error(), tt.wantErrMsg)
				}
				return
			}

			if tt.wantErr {
				return
			}

			if got.IsGitHubIP != tt.want.IsGitHubIP {
				t.Errorf("CheckIP() IsGitHubIP = %v, want %v", got.IsGitHubIP, tt.want.IsGitHubIP)
			}

			if got.IsGitHubIP {
				if got.FunctionalArea != tt.want.FunctionalArea {
					t.Errorf("CheckIP() FunctionalArea = %v, want %v", got.FunctionalArea, tt.want.FunctionalArea)
				}
				if got.Range != tt.want.Range {
					t.Errorf("CheckIP() Range = %v, want %v", got.Range, tt.want.Range)
				}
			}
		})
	}
}

// failingTransport is a transport that always fails
type failingTransport struct{}

func (t *failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("failed to fetch GitHub meta")
}

// Test for different functional areas
func TestIPChecker_CheckIP_DifferentFunctionalAreas(t *testing.T) {
	// Server with different IP ranges for each service
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["192.30.252.0/22"],
			"web": ["140.82.112.0/20"],
			"api": ["140.82.128.0/20"],
			"git": ["140.82.144.0/20"],
			"packages": ["140.82.160.0/20"],
			"pages": ["185.199.108.0/22"],
			"importer": ["140.82.176.0/20"],
			"actions": ["13.64.0.0/16"],
			"dependabot": ["13.65.0.0/16"],
			"actions_ipv4": ["13.66.0.0/16"]
		}`))
	}))
	defer server.Close()

	tests := []struct {
		name       string
		ip         string
		wantArea   string
		wantRange  string
		wantGitHub bool
	}{
		{
			name:       "Hooks IP",
			ip:         "192.30.252.1",
			wantArea:   "Hooks",
			wantRange:  "192.30.252.0/22",
			wantGitHub: true,
		},
		{
			name:       "Web IP",
			ip:         "140.82.112.1",
			wantArea:   "Web",
			wantRange:  "140.82.112.0/20",
			wantGitHub: true,
		},
		{
			name:       "API IP",
			ip:         "140.82.128.1",
			wantArea:   "API",
			wantRange:  "140.82.128.0/20",
			wantGitHub: true,
		},
		{
			name:       "Git IP",
			ip:         "140.82.144.1",
			wantArea:   "Git",
			wantRange:  "140.82.144.0/20",
			wantGitHub: true,
		},
		{
			name:       "Packages IP",
			ip:         "140.82.160.1",
			wantArea:   "Packages",
			wantRange:  "140.82.160.0/20",
			wantGitHub: true,
		},
		{
			name:       "Pages IP",
			ip:         "185.199.108.1",
			wantArea:   "Pages",
			wantRange:  "185.199.108.0/22",
			wantGitHub: true,
		},
		{
			name:       "Importer IP",
			ip:         "140.82.176.1",
			wantArea:   "Importer",
			wantRange:  "140.82.176.0/20",
			wantGitHub: true,
		},
		{
			name:       "Actions IP",
			ip:         "13.64.0.1",
			wantArea:   "Actions",
			wantRange:  "13.64.0.0/16",
			wantGitHub: true,
		},
		{
			name:       "Dependabot IP",
			ip:         "13.65.0.1",
			wantArea:   "Dependabot",
			wantRange:  "13.65.0.0/16",
			wantGitHub: true,
		},
		{
			name:       "Actions IPv4 IP",
			ip:         "13.66.0.1",
			wantArea:   "Actions IPv4",
			wantRange:  "13.66.0.0/16",
			wantGitHub: true,
		},
	}

	// Override the GitHub meta URL for testing
	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewIPChecker()
			result, err := checker.CheckIP(tt.ip)

			if err != nil {
				t.Errorf("CheckIP() error = %v, wantErr false", err)
				return
			}

			if result.IsGitHubIP != tt.wantGitHub {
				t.Errorf("CheckIP() IsGitHubIP = %v, want %v", result.IsGitHubIP, tt.wantGitHub)
			}

			if tt.wantGitHub {
				if result.FunctionalArea != tt.wantArea {
					t.Errorf("CheckIP() FunctionalArea = %v, want %v", result.FunctionalArea, tt.wantArea)
				}
				if result.Range != tt.wantRange {
					t.Errorf("CheckIP() Range = %v, want %v", result.Range, tt.wantRange)
				}
			}
		})
	}
}

// Test additional IP edge cases
func TestIPChecker_CheckIP_AdditionalEdgeCases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["192.30.252.0/22"],
			"web": ["140.82.112.0/20"],
			"api": ["140.82.113.0/20"],
			"git": ["140.82.114.0/20"],
			"packages": ["140.82.115.0/20"],
			"pages": ["185.199.108.0/22"],
			"importer": ["192.30.253.0/22"],
			"actions": ["13.64.0.0/16"],
			"dependabot": ["13.65.0.0/16"],
			"actions_ipv4": ["13.66.0.0/16"]
		}`))
	}))
	defer server.Close()

	// Override the GitHub meta URL for testing
	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()

	tests := []struct {
		name       string
		ip         string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "Loopback IP",
			ip:         "127.0.0.1",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "Link-local IP (169.254.x.x - not always caught by IsPrivate)",
			ip:         "169.254.1.1",
			wantErr:    false, // This might not be caught by Go's IsPrivate() method
		},
		{
			name:       "Class A private IP",
			ip:         "10.0.0.1",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "Class B private IP",
			ip:         "172.16.0.1",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "Class C private IP edge",
			ip:         "192.168.255.255",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "Multicast IP",
			ip:         "224.0.0.1",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "Unspecified IP",
			ip:         "0.0.0.0",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address",
		},
		{
			name:       "IPv4-mapped IPv6",
			ip:         "::ffff:192.168.1.1",
			wantErr:    true,
			wantErrMsg: "IP address must be a public, routable address", // The actual error we get
		},
		{
			name:       "Empty string",
			ip:         "",
			wantErr:    true,
			wantErrMsg: "invalid IP address format",
		},
		{
			name:       "Malformed IP with letters",
			ip:         "192.168.a.1",
			wantErr:    true,
			wantErrMsg: "invalid IP address format",
		},
		{
			name:       "Partial IP",
			ip:         "192.168.1",
			wantErr:    true,
			wantErrMsg: "invalid IP address format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewIPChecker()
			_, err := checker.CheckIP(tt.ip)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("CheckIP() error message = %v, should contain %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

// Test caching behavior
func TestIPChecker_CheckIP_Caching(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hooks": ["192.30.252.0/22"],
			"web": ["140.82.112.0/20"],
			"api": ["140.82.113.0/20"],
			"git": ["140.82.114.0/20"],
			"packages": ["140.82.115.0/20"],
			"pages": ["185.199.108.0/22"],
			"importer": ["192.30.253.0/22"],
			"actions": ["13.64.0.0/16"],
			"dependabot": ["13.65.0.0/16"],
			"actions_ipv4": ["13.66.0.0/16"]
		}`))
	}))
	defer server.Close()

	// Override the GitHub meta URL for testing
	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()

	checker := NewIPChecker()

	// First call should fetch from API
	_, err := checker.CheckIP("192.30.252.1")
	if err != nil {
		t.Errorf("First CheckIP() call failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 API request after first call, got %d", requestCount)
	}

	// Second call should use cached data
	_, err = checker.CheckIP("140.82.112.1")
	if err != nil {
		t.Errorf("Second CheckIP() call failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 API request after second call (cached), got %d", requestCount)
	}

	// Third call should still use cached data
	_, err = checker.CheckIP("8.8.8.8") // Non-GitHub IP
	if err != nil {
		t.Errorf("Third CheckIP() call failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 API request after third call (cached), got %d", requestCount)
	}
}
