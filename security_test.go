package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestHTTPClientTimeout verifies that the HTTP client has proper timeout configuration
func TestHTTPClientTimeout(t *testing.T) {
	checker := NewIPChecker()
	
	// Check that the client has a timeout configured
	if checker.client.Timeout == 0 {
		t.Error("HTTP client should have a timeout configured")
	}
	
	// Verify timeout is reasonable (not too short, not too long)
	expectedTimeout := 30 * time.Second
	if checker.client.Timeout != expectedTimeout {
		t.Errorf("Expected timeout of %v, got %v", expectedTimeout, checker.client.Timeout)
	}
}

// TestHTTPTransportSecurity verifies that the HTTP transport has security configurations
func TestHTTPTransportSecurity(t *testing.T) {
	checker := NewIPChecker()
	
	transport, ok := checker.client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected *http.Transport, got different type")
	}
	
	// Verify TLS handshake timeout
	if transport.TLSHandshakeTimeout == 0 {
		t.Error("TLS handshake timeout should be configured")
	}
	
	// Verify response header timeout
	if transport.ResponseHeaderTimeout == 0 {
		t.Error("Response header timeout should be configured")
	}
}

// TestResponseSizeLimit verifies that large responses are handled safely
func TestResponseSizeLimit(t *testing.T) {
	// Create a server that returns a very large response
	largeResponse := strings.Repeat("x", 2*1024*1024) // 2MB response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"hooks": ["` + largeResponse + `"]}`))
	}))
	defer server.Close()
	
	// Replace the GitHub meta URL with our test server
	originalURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = originalURL }()
	
	checker := NewIPChecker()
	_, err := checker.CheckIP("8.8.8.8")
	
	// Should get an error due to response size limit
	if err == nil {
		t.Error("Expected error for oversized response")
	}
	
	if !strings.Contains(err.Error(), "failed to fetch GitHub IP ranges") {
		t.Errorf("Expected sanitized error message, got: %v", err)
	}
}

// TestSlowServerTimeout verifies that slow servers are handled with timeouts
func TestSlowServerTimeout(t *testing.T) {
	// Create a server that takes longer than our timeout to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(35 * time.Second) // Longer than our 30s timeout
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"hooks": ["192.30.252.0/22"]}`))
	}))
	defer server.Close()
	
	// Replace the GitHub meta URL with our test server
	originalURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = originalURL }()
	
	checker := NewIPChecker()
	start := time.Now()
	_, err := checker.CheckIP("8.8.8.8")
	duration := time.Since(start)
	
	// Should timeout and return an error
	if err == nil {
		t.Error("Expected timeout error for slow server")
	}
	
	// Should timeout within reasonable time (allowing some buffer for test execution)
	if duration > 35*time.Second {
		t.Errorf("Request took too long to timeout: %v", duration)
	}
	
	// Should get sanitized error message
	if !strings.Contains(err.Error(), "failed to fetch GitHub IP ranges") {
		t.Errorf("Expected sanitized error message, got: %v", err)
	}
}

// TestErrorMessageSanitization verifies that error messages don't leak sensitive information
func TestErrorMessageSanitization(t *testing.T) {
	checker := NewIPChecker()
	
	// Test with failing client to trigger network error
	failingClient := &http.Client{
		Transport: &failingTransport{},
	}
	checker.setClient(failingClient)
	
	_, err := checker.CheckIP("8.8.8.8")
	
	if err == nil {
		t.Error("Expected error from failing client")
	}
	
	errMsg := err.Error()
	
	// Should not contain detailed network error information
	sensitivePatterns := []string{
		"connection refused",
		"no such host",
		"timeout",
		"network is unreachable",
	}
	
	for _, pattern := range sensitivePatterns {
		if strings.Contains(strings.ToLower(errMsg), pattern) {
			t.Errorf("Error message contains sensitive information: %v", errMsg)
		}
	}
	
	// Should contain generic error message
	if !strings.Contains(errMsg, "failed to fetch GitHub IP ranges") {
		t.Errorf("Error message should be sanitized, got: %v", errMsg)
	}
}

// failingTransport is a transport that always fails (already exists in ip_checker_test.go)
// We can reuse it here by making it part of the test helper functions