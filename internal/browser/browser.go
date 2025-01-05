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
