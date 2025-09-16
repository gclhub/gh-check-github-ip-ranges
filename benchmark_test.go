package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Benchmark testing IP checking performance
func BenchmarkIPChecker_CheckIP_GitHubIP(b *testing.B) {
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

	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()

	checker := NewIPChecker()
	
	// Pre-load cache
	_, err := checker.CheckIP("192.30.252.1")
	if err != nil {
		b.Fatal("Failed to pre-load cache:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checker.CheckIP("192.30.252.1")
		if err != nil {
			b.Fatal("CheckIP failed:", err)
		}
	}
}

func BenchmarkIPChecker_CheckIP_NonGitHubIP(b *testing.B) {
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

	oldURL := githubMetaURL
	githubMetaURL = server.URL
	defer func() { githubMetaURL = oldURL }()

	checker := NewIPChecker()
	
	// Pre-load cache
	_, err := checker.CheckIP("8.8.8.8")
	if err != nil {
		b.Fatal("Failed to pre-load cache:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checker.CheckIP("8.8.8.8")
		if err != nil {
			b.Fatal("CheckIP failed:", err)
		}
	}
}

func BenchmarkIPChecker_CheckIP_InvalidIP(b *testing.B) {
	checker := NewIPChecker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = checker.CheckIP("invalid-ip")
		// We expect an error, so we ignore it
	}
}

func BenchmarkIsBroadcastAddress(b *testing.B) {
	broadcastIP := []byte{255, 255, 255, 255}
	nonBroadcastIP := []byte{192, 168, 1, 1}

	b.Run("BroadcastIP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isBroadcastAddress(broadcastIP)
		}
	})

	b.Run("NonBroadcastIP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isBroadcastAddress(nonBroadcastIP)
		}
	})
}

func BenchmarkNewIPChecker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewIPChecker()
	}
}