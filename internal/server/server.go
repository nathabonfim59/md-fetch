package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

type Server struct {
	port int
}

type FetchRequest struct {
	URLs    []string `json:"urls"`
	Browser string  `json:"browser,omitempty"`
}

type FetchResponse struct {
	Results map[string]string `json:"results"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func New(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/fetch", s.handleFetch)
	http.HandleFunc("/openapi.yaml", s.handleOpenAPI)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := make(map[string]string)
	errors := make(map[string]string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, url := range req.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			content, err := fetcher.FetchContent(url, req.Browser)
			
			mu.Lock()
			defer mu.Unlock()
			
			if err != nil {
				errors[url] = err.Error()
				return
			}
			results[url] = content
		}(url)
	}

	wg.Wait()

	response := FetchResponse{
		Results: results,
	}
	if len(errors) > 0 {
		response.Errors = errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleOpenAPI(w http.ResponseWriter, r *http.Request) {
	openAPISpec := `openapi: 3.0.0
info:
  title: md-fetch API
  description: API for fetching web content and converting it to Markdown
  version: 1.0.0

servers:
  - url: http://localhost:{port}
    variables:
      port:
        default: "8080"

paths:
  /fetch:
    post:
      summary: Fetch content from URLs and convert to Markdown
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                urls:
                  type: array
                  items:
                    type: string
                  description: List of URLs to fetch
                browser:
                  type: string
                  enum: [chrome, firefox, curl]
                  description: Browser to use for fetching (optional)
              required:
                - urls
      responses:
        '200':
          description: Successful operation
          content:
            application/json:
              schema:
                type: object
                properties:
                  results:
                    type: object
                    additionalProperties:
                      type: string
                    description: Map of URLs to their fetched content
                  errors:
                    type: object
                    additionalProperties:
                      type: string
                    description: Map of URLs to error messages (if any)
        '400':
          description: Invalid request
        '405':
          description: Method not allowed`

	w.Header().Set("Content-Type", "text/yaml")
	fmt.Fprint(w, openAPISpec)
}
