package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		silent     bool
		wantError  bool
		wantStdout string
		wantStderr string
	}{
		{
			name:       "GitHub IP with silent mode",
			args:       []string{"192.30.252.1"},
			silent:     true,
			wantError:  false,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:       "GitHub IP without silent mode",
			args:       []string{"192.30.252.1"},
			silent:     false,
			wantError:  false,
			wantStdout: "IP 192.30.252.1 belongs to GitHub's Hooks range (192.30.252.0/22)\n",
			wantStderr: "",
		},
		{
			name:       "Non-GitHub IP with silent mode",
			args:       []string{"8.8.8.8"},
			silent:     true,
			wantError:  true,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:       "Non-GitHub IP without silent mode",
			args:       []string{"8.8.8.8"},
			silent:     false,
			wantError:  true,
			wantStdout: "",
			wantStderr: "",  // Error message is handled by main(), not runCommand
		},
		{
			name:       "Invalid IP with silent mode",
			args:      []string{"invalid-ip"},
			silent:    true,
			wantError: true,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:      "Invalid IP without silent mode",
			args:      []string{"invalid-ip"},
			silent:    false,
			wantError: true,
			wantStdout: "",
			wantStderr: "",  // Error message is handled by main(), not runCommand
		},
	}

	// Mock server for GitHub meta API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			// Create command and set flags
			cmd := &cobra.Command{}
			cmd.Flags().BoolP("silent", "s", false, "")
			if tt.silent {
				cmd.Flags().Set("silent", "true")
			}

			// Override githubMetaURL for testing
			oldURL := githubMetaURL
			githubMetaURL = server.URL
			defer func() { githubMetaURL = oldURL }()

			// Run command
			err := runCommand(cmd, tt.args)

			// Restore stdout and stderr
			wOut.Close()
			wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// Read output
			var bufOut, bufErr bytes.Buffer
			bufOut.ReadFrom(rOut)
			bufErr.ReadFrom(rErr)
			stdout := bufOut.String()
			stderr := bufErr.String()

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("runCommand() error = %v, wantError %v", err, tt.wantError)
			}

			// Check stdout
			if stdout != tt.wantStdout {
				t.Errorf("runCommand() stdout = %q, want %q", stdout, tt.wantStdout)
			}

			// Check stderr
			if stderr != tt.wantStderr {
				t.Errorf("runCommand() stderr = %q, want %q", stderr, tt.wantStderr)
			}
		})
	}
}

func TestMainFunction(t *testing.T) {
	// Mock server for GitHub meta API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	// Save original values
	oldArgs := os.Args
	oldURL := githubMetaURL
	oldOsExit := osExit
	defer func() {
		os.Args = oldArgs
		githubMetaURL = oldURL
		osExit = oldOsExit
	}()

	// Override githubMetaURL for testing
	githubMetaURL = server.URL

	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantErr  bool
		silent   bool
	}{
		{
			name:     "Valid GitHub IP",
			args:     []string{"gh-check-github-ip-ranges", "192.30.252.1"},
			wantCode: 0,
			wantErr:  false,
			silent:   false,
		},
		{
			name:     "Non-GitHub IP",
			args:     []string{"gh-check-github-ip-ranges", "8.8.8.8"},
			wantCode: 1,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "Invalid IP format",
			args:     []string{"gh-check-github-ip-ranges", "invalid-ip"},
			wantCode: 2,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "Private IP",
			args:     []string{"gh-check-github-ip-ranges", "192.168.1.1"},
			wantCode: 2,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "IPv6 address",
			args:     []string{"gh-check-github-ip-ranges", "2001:db8::1"},
			wantCode: 2,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "Broadcast address",
			args:     []string{"gh-check-github-ip-ranges", "255.255.255.255"},
			wantCode: 2,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "No arguments",
			args:     []string{"gh-check-github-ip-ranges"},
			wantCode: 2,
			wantErr:  true,
			silent:   false,
		},
		{
			name:     "Help flag",
			args:     []string{"gh-check-github-ip-ranges", "--help"},
			wantCode: 0,
			wantErr:  false,
			silent:   false,
		},
		{
			name:     "Non-GitHub IP with silent mode",
			args:     []string{"gh-check-github-ip-ranges", "8.8.8.8", "-s"},
			wantCode: 1,
			wantErr:  true,
			silent:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				// Capture stderr for checking error messages
				oldStdout := os.Stdout
				oldStderr := os.Stderr
				_, wOut, _ := os.Pipe()
				rErr, wErr, _ := os.Pipe()
				os.Stdout = wOut
				os.Stderr = wErr

				// Set up test args
				os.Args = tt.args

				// Create exit code channel
				exitCode := make(chan int, 1)
				osExit = func(code int) {
					exitCode <- code
					// Don't actually exit in tests
				}

				// Run main in goroutine
				go func() {
					main()
					exitCode <- 0
				}()

				// Get exit code
				code := <-exitCode
				
				// Close pipes
				wOut.Close()
				wErr.Close()
				os.Stdout = oldStdout
				os.Stderr = oldStderr

				// Read stderr
				var bufErr bytes.Buffer
				bufErr.ReadFrom(rErr)
				stderr := bufErr.String()

				if code != tt.wantCode {
					t.Errorf("main() exitCode = %v, want %v", code, tt.wantCode)
				}

				// For non-GitHub IPs, verify no "Error: " prefix
				if code == 1 && !tt.silent {
					expectedMsg := "The provided IP address is not a GitHub-owned address\n"
					if stderr != expectedMsg {
						t.Errorf("main() stderr = %q, want %q", stderr, expectedMsg)
					}
				}
		})
	}
}