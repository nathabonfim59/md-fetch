package browser

import (
	"fmt"
	"os/exec"
)

// Browser represents a web browser interface
type Browser interface {
	// Fetch retrieves content from a URL
	Fetch(url string) ([]byte, error)
	// Name returns the browser's name
	Name() string
}

// ExecutableFinder is an interface for finding browser executables
type ExecutableFinder interface {
	Find() (string, error)
}

// DefaultExecutableFinder implements ExecutableFinder
type DefaultExecutableFinder struct {
	names []string
}

func (f *DefaultExecutableFinder) Find() (string, error) {
	for _, name := range f.names {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no browser executable found for: %v", f.names)
}

// DefaultBrowsers defines the priority order for browsers
var DefaultBrowsers = []string{"chrome", "firefox", "curl"}

// GetDefaultBrowser tries browsers in order of preference and returns the first available one
func GetDefaultBrowser() (Browser, error) {
	var lastErr error
	for _, browserType := range DefaultBrowsers {
		browser, err := NewBrowser(browserType)
		if err == nil {
			return browser, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("no supported browsers found: %v", lastErr)
}

// NewBrowser creates a new browser instance based on the browser type
func NewBrowser(browserType string) (Browser, error) {
	switch browserType {
	case "chrome", "chromium":
		return NewChrome()
	case "firefox":
		return NewFirefox()
	case "links":
		return NewLinks()
	case "lynx":
		return NewLynx()
	case "w3m":
		return NewW3m()
	case "curl":
		return NewCurl()
	default:
		return nil, fmt.Errorf("unsupported browser type: %s", browserType)
	}
}
