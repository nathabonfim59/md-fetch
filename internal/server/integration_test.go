package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServerIntegration(t *testing.T) {
	// Start server in a goroutine
	port := 8081
	srv := New(port)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test cases
	tests := []struct {
		name        string
		request     FetchRequest
		validateResp func(t *testing.T, resp *FetchResponse)
	}{
		{
			name: "fetch single URL",
			request: FetchRequest{
				URLs:    []string{"https://example.com"},
				Browser: "chrome",
			},
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Results) != 1 {
					t.Errorf("expected 1 result, got %d", len(resp.Results))
				}
				content, ok := resp.Results["https://example.com"]
				if !ok {
					t.Error("expected result for example.com")
					return
				}
				if !contains(content, []string{"Example Domain", "illustrative examples"}) {
					t.Error("expected content not found in response")
				}
			},
		},
		{
			name: "fetch multiple URLs",
			request: FetchRequest{
				URLs:    []string{"https://example.com", "https://example.org"},
				Browser: "chrome",
			},
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Results) != 2 {
					t.Errorf("expected 2 results, got %d", len(resp.Results))
				}
			},
		},
		{
			name: "handle invalid URL",
			request: FetchRequest{
				URLs:    []string{"not-a-valid-url"},
				Browser: "chrome",
			},
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Errors) != 1 {
					t.Errorf("expected 1 error, got %d", len(resp.Errors))
				}
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			resp, err := http.Post(
				fmt.Sprintf("http://localhost:%d/fetch", port),
				"application/json",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK, got %v", resp.Status)
			}

			var result FetchResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			tt.validateResp(t, &result)
		})
	}
}

// Helper function to check if content contains all expected strings
func contains(content string, expected []string) bool {
	for _, exp := range expected {
		if !strings.Contains(content, exp) {
			return false
		}
	}
	return true
}
