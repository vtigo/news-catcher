package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with default values", func(t *testing.T) {
		client := NewClient()

		if client == nil {
			t.Fatal("expected non-nil client")
		}

		if client.client.Timeout != 30*time.Second {
			t.Errorf("expected timeout of 30s, got %v", client.client.Timeout)
		}

		if client.maxBytes != 10*1024*1024 {
			t.Errorf("expected maxBytes of 10MB, got %d", client.maxBytes)
		}
	})

	t.Run("applies timeout option", func(t *testing.T) {
		timeout := 5 * time.Second
		client := NewClient(WithTimeout(timeout))

		if client.client.Timeout != timeout {
			t.Errorf("expected timeout of %v, got %v", timeout, client.client.Timeout)
		}
	})

	t.Run("applies max bytes option", func(t *testing.T) {
		maxBytes := int64(1024)
		client := NewClient(WithMaxBytes(maxBytes))

		if client.maxBytes != maxBytes {
			t.Errorf("expected maxBytes of %d, got %d", maxBytes, client.maxBytes)
		}
	})

	t.Run("applies multiple options", func(t *testing.T) {
		timeout := 15 * time.Second
		maxBytes := int64(2048)
		client := NewClient(
			WithTimeout(timeout),
			WithMaxBytes(maxBytes),
		)

		if client.client.Timeout != timeout {
			t.Errorf("expected timeout of %v, got %v", timeout, client.client.Timeout)
		}

		if client.maxBytes != maxBytes {
			t.Errorf("expected maxBytes of %d, got %d", maxBytes, client.maxBytes)
		}
	})
}

func TestFetcherClient_Fetch(t *testing.T) {
	t.Run("successfully fetches content", func(t *testing.T) {
		expectedBody := "test content"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedBody))
		}))
		defer server.Close()

		client := NewClient()
		ctx := context.Background()

		body, err := client.Fetch(ctx, server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(body) != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, string(body))
		}
	})

	t.Run("sets user agent header", func(t *testing.T) {
		var receivedUA string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedUA = r.Header.Get("User-Agent")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient()
		ctx := context.Background()

		_, err := client.Fetch(ctx, server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedUA := "NewsCatcher/1.0"
		if receivedUA != expectedUA {
			t.Errorf("expected User-Agent %q, got %q", expectedUA, receivedUA)
		}
	})

	t.Run("returns error on non-200 status code", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
		}{
			{"not found", http.StatusNotFound},
			{"internal server error", http.StatusInternalServerError},
			{"bad request", http.StatusBadRequest},
			{"forbidden", http.StatusForbidden},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
				}))
				defer server.Close()

				client := NewClient()
				ctx := context.Background()

				_, err := client.Fetch(ctx, server.URL)
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				expectedErr := "bad status code"
				if !strings.Contains(err.Error(), expectedErr) {
					t.Errorf("expected error to contain %q, got %q", expectedErr, err.Error())
				}
			})
		}
	})

	t.Run("respects max bytes limit", func(t *testing.T) {
		largeContent := strings.Repeat("a", 2048)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(largeContent))
		}))
		defer server.Close()

		maxBytes := int64(1024)
		client := NewClient(WithMaxBytes(maxBytes))
		ctx := context.Background()

		body, err := client.Fetch(ctx, server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if int64(len(body)) != maxBytes {
			t.Errorf("expected body length of %d, got %d", maxBytes, len(body))
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.Fetch(ctx, server.URL)
		if err == nil {
			t.Fatal("expected error due to cancelled context, got nil")
		}
	})

	t.Run("respects context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := client.Fetch(ctx, server.URL)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	})

	t.Run("returns error on invalid URL", func(t *testing.T) {
		client := NewClient()
		ctx := context.Background()

		_, err := client.Fetch(ctx, "://invalid-url")
		if err == nil {
			t.Fatal("expected error for invalid URL, got nil")
		}
	})

	t.Run("handles server connection error", func(t *testing.T) {
		client := NewClient(WithTimeout(1 * time.Second))
		ctx := context.Background()

		// Use a URL that won't respond
		_, err := client.Fetch(ctx, "http://localhost:9999")
		if err == nil {
			t.Fatal("expected connection error, got nil")
		}
	})
}

func TestFetcherClient_Fetch_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("handles large response body", func(t *testing.T) {
		largeContent := strings.Repeat("x", 5*1024*1024) // 5MB
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, largeContent)
		}))
		defer server.Close()

		client := NewClient()
		ctx := context.Background()

		body, err := client.Fetch(ctx, server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(body) != len(largeContent) {
			t.Errorf("expected body length %d, got %d", len(largeContent), len(body))
		}
	})
}

// Benchmark tests
func BenchmarkFetcherClient_Fetch(b *testing.B) {
	content := "benchmark content"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(content))
	}))
	defer server.Close()

	client := NewClient()
	ctx := context.Background()

	for b.Loop() {
		_, err := client.Fetch(ctx, server.URL)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
