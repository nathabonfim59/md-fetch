package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFetch(t *testing.T) {
	srv := New(8080)

	tests := []struct {
		name           string
		method        string
		requestBody   interface{}
		expectedCode  int
		validateResp  func(t *testing.T, resp *FetchResponse)
	}{
		{
			name:   "valid single URL request",
			method: http.MethodPost,
			requestBody: FetchRequest{
				URLs:    []string{"https://example.com"},
				Browser: "chrome",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Results) != 1 {
					t.Errorf("expected 1 result, got %d", len(resp.Results))
				}
				if content, ok := resp.Results["https://example.com"]; !ok || !strings.Contains(content, "Example Domain") {
					t.Error("expected example.com content in results")
				}
			},
		},
		{
			name:   "valid multiple URLs request",
			method: http.MethodPost,
			requestBody: FetchRequest{
				URLs:    []string{"https://example.com", "https://example.org"},
				Browser: "chrome",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Results) != 2 {
					t.Errorf("expected 2 results, got %d", len(resp.Results))
				}
			},
		},
		{
			name:   "invalid URL request",
			method: http.MethodPost,
			requestBody: FetchRequest{
				URLs:    []string{"not-a-valid-url"},
				Browser: "chrome",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Errors) != 1 {
					t.Errorf("expected 1 error, got %d", len(resp.Errors))
				}
			},
		},
		{
			name:   "empty URLs request",
			method: http.MethodPost,
			requestBody: FetchRequest{
				URLs:    []string{},
				Browser: "chrome",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, resp *FetchResponse) {
				if len(resp.Results) != 0 {
					t.Errorf("expected 0 results, got %d", len(resp.Results))
				}
			},
		},
		{
			name:         "invalid method",
			method:      http.MethodGet,
			requestBody: nil,
			expectedCode: http.StatusMethodNotAllowed,
			validateResp: nil,
		},
		{
			name:   "invalid request body",
			method: http.MethodPost,
			requestBody: "invalid json",
			expectedCode: http.StatusBadRequest,
			validateResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.requestBody != nil {
				if jsonStr, ok := tt.requestBody.(string); ok {
					body = bytes.NewBufferString(jsonStr)
				} else {
					jsonData, err := json.Marshal(tt.requestBody)
					if err != nil {
						t.Fatalf("failed to marshal request body: %v", err)
					}
					body = bytes.NewBuffer(jsonData)
				}
			}

			req := httptest.NewRequest(tt.method, "/fetch", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			srv.handleFetch(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.validateResp != nil {
				var resp FetchResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.validateResp(t, &resp)
			}
		})
	}
}

func TestHandleOpenAPI(t *testing.T) {
	srv := New(8080)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	w := httptest.NewRecorder()

	srv.handleOpenAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/yaml" {
		t.Errorf("expected Content-Type %s, got %s", "text/yaml", contentType)
	}

	body := w.Body.String()
	expectedFields := []string{
		"openapi: 3.0.0",
		"title: md-fetch API",
		"/fetch:",
		"post:",
		"application/json",
	}

	for _, field := range expectedFields {
		if !strings.Contains(body, field) {
			t.Errorf("expected OpenAPI spec to contain %q", field)
		}
	}
}

func TestNew(t *testing.T) {
	port := 8080
	srv := New(port)

	if srv == nil {
		t.Error("expected non-nil server")
	}

	if srv.port != port {
		t.Errorf("expected port %d, got %d", port, srv.port)
	}
}
