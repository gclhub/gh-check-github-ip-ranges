package main

import (
	"net/http"
	"net/http/httptest"
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

	tests := []struct {
		name       string
		ip         string
		mockServer *httptest.Server
		wantErr    bool
		wantErrMsg string
		want       *CheckResult
	}{
		{
			name:       "Valid GitHub IP",
			ip:         "192.30.252.1",
			mockServer: successServer,
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
			wantErr:    false,
			want: &CheckResult{
				IsGitHubIP: false,
			},
		},
		{
			name:       "Invalid IP format",
			ip:         "invalid-ip",
			mockServer: successServer,
			wantErr:    true,
			wantErrMsg: "Error: invalid IP address format",
			want:       nil,
		},
		{
			name:       "Private IP",
			ip:         "192.168.1.1",
			mockServer: successServer,
			wantErr:    true,
			wantErrMsg: "Error: IP address must be a public, routable address",
			want:       nil,
		},
		{
			name:       "IPv6 address",
			ip:         "2001:db8::1",
			mockServer: successServer,
			wantErr:    true,
			wantErrMsg: "Error: only IPv4 addresses are supported",
			want:       nil,
		},
		{
			name:       "Broadcast address",
			ip:         "255.255.255.255",
			mockServer: successServer,
			wantErr:    true,
			wantErrMsg: "Error: IP address must be a public, routable address",
			want:       nil,
		},
		{
			name:       "API Server Error",
			ip:         "8.8.8.8",
			mockServer: errorServer,
			wantErr:    true,
			want:       nil,
		},
		{
			name:       "Invalid JSON Response",
			ip:         "8.8.8.8",
			mockServer: invalidJSONServer,
			wantErr:    true,
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the GitHub meta URL for testing
			oldURL := githubMetaURL
			githubMetaURL = tt.mockServer.URL
			defer func() { githubMetaURL = oldURL }()

			checker := NewIPChecker()
			got, err := checker.CheckIP(tt.ip)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErrMsg != "" && err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("CheckIP() error message = %v, want %v", err.Error(), tt.wantErrMsg)
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