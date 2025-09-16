package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSecurityEnhancements(t *testing.T) {
	t.Run("Input Length Validation", func(t *testing.T) {
		checker := NewIPChecker()
		
		// Test extremely long input
		longInput := strings.Repeat("1", 50) // More than maxIPStringLength (45)
		_, err := checker.CheckIP(longInput)
		if err == nil {
			t.Error("Expected error for overly long input")
		}
		if !strings.Contains(err.Error(), "IP address string too long") {
			t.Errorf("Expected length validation error, got: %v", err)
		}
	})
	
	t.Run("Empty Input Validation", func(t *testing.T) {
		checker := NewIPChecker()
		
		_, err := checker.CheckIP("")
		if err == nil {
			t.Error("Expected error for empty input")
		}
		if !strings.Contains(err.Error(), "IP address cannot be empty") {
			t.Errorf("Expected empty input error, got: %v", err)
		}
	})
	
	t.Run("HTTP Client Security Configuration", func(t *testing.T) {
		checker := NewIPChecker()
		
		// Verify the client has timeout set
		transport, ok := checker.client.Transport.(*http.Transport)
		if !ok {
			t.Fatal("Expected HTTP transport")
		}
		
		// Check TLS configuration
		if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
			t.Errorf("Expected minimum TLS version 1.2, got %d", transport.TLSClientConfig.MinVersion)
		}
		
		// Check timeout
		if checker.client.Timeout != 30*time.Second {
			t.Errorf("Expected 30 second timeout, got %v", checker.client.Timeout)
		}
	})
	
	t.Run("Response Size Limit", func(t *testing.T) {
		// Create a server that returns a response larger than maxResponseSize
		largeResponseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			
			// Write a response larger than 1MB
			largeData := strings.Repeat("x", 1<<20+1) // 1MB + 1 byte
			w.Write([]byte(`{"hooks": ["` + largeData + `"]}`))
		}))
		defer largeResponseServer.Close()
		
		// Override githubMetaURL for testing
		oldURL := githubMetaURL
		githubMetaURL = largeResponseServer.URL
		defer func() { githubMetaURL = oldURL }()
		
		checker := NewIPChecker()
		_, err := checker.CheckIP("8.8.8.8")
		if err == nil {
			t.Error("Expected error for oversized response")
		}
		if !strings.Contains(err.Error(), "failed to fetch GitHub meta") {
			t.Errorf("Expected fetch error due to size limit, got: %v", err)
		}
	})
	
	t.Run("Timeout Protection", func(t *testing.T) {
		// Create a server that hangs indefinitely
		hangingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Don't respond, causing a timeout
			time.Sleep(35 * time.Second) // Longer than client timeout
		}))
		defer hangingServer.Close()
		
		// Override githubMetaURL for testing
		oldURL := githubMetaURL
		githubMetaURL = hangingServer.URL
		defer func() { githubMetaURL = oldURL }()
		
		checker := NewIPChecker()
		start := time.Now()
		_, err := checker.CheckIP("8.8.8.8")
		elapsed := time.Since(start)
		
		if err == nil {
			t.Error("Expected timeout error")
		}
		
		// Should timeout within reasonable time (allowing some buffer)
		if elapsed > 35*time.Second {
			t.Errorf("Request took too long: %v", elapsed)
		}
		
		if !strings.Contains(err.Error(), "failed to fetch GitHub meta") {
			t.Errorf("Expected timeout error, got: %v", err)
		}
	})
	
	t.Run("TLS Certificate Validation", func(t *testing.T) {
		// Create an HTTPS server with self-signed certificate
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"hooks": ["192.30.252.0/22"]}`))
		}))
		defer server.Close()
		
		// Override githubMetaURL for testing
		oldURL := githubMetaURL
		githubMetaURL = server.URL
		defer func() { githubMetaURL = oldURL }()
		
		checker := NewIPChecker()
		_, err := checker.CheckIP("192.30.252.1")
		
		// Should fail due to certificate validation (self-signed cert)
		if err == nil {
			t.Error("Expected TLS certificate validation error")
		}
		if !strings.Contains(err.Error(), "failed to fetch GitHub meta") {
			t.Errorf("Expected TLS error, got: %v", err)
		}
	})
}

func TestSecurityFuzzInputs(t *testing.T) {
	// Test various malformed inputs that could cause issues
	testCases := []struct {
		name  string
		input string
	}{
		{"Unicode characters", "192.168.1.1\u0000"},
		{"Control characters", "192.168.1.1\n\r"},
		{"Mixed valid/invalid", "192.168.1.1.extra"},
		{"Special characters", "192.168.1.1;whoami"},
		{"SQL injection attempt", "192.168.1.1'; DROP TABLE--;"},
		{"Path traversal", "192.168.1.1/../../../etc/passwd"},
		{"Buffer overflow attempt", strings.Repeat("A", 1000)},
		{"Null bytes", "192.168.1.1\x00"},
		{"Format string", "192.168.1.%s%s%s"},
		{"JavaScript", "192.168.1.1<script>alert(1)</script>"},
	}
	
	checker := NewIPChecker()
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := checker.CheckIP(tc.input)
			// All of these should return errors, none should panic or hang
			if err == nil {
				t.Errorf("Expected error for malicious input: %q", tc.input)
			}
		})
	}
}

func TestConcurrentSafety(t *testing.T) {
	// Test that concurrent access doesn't cause race conditions
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
	
	// Override githubMetaURL for testing
	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()
	
	checker := NewIPChecker()
	
	// Run multiple goroutines concurrently
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	errChan := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			_, err := checker.CheckIP("192.30.252.1")
			errChan <- err
		}()
	}
	
	// Collect results
	for i := 0; i < 10; i++ {
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("Concurrent access failed: %v", err)
			}
		case <-ctx.Done():
			t.Fatal("Test timed out")
		}
	}
}