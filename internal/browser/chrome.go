package browser

import (
	"fmt"
	"os/exec"
)

type Chrome struct {
	execPath string
	cleaningOpts *CleaningOptions
}

func NewChrome() (*Chrome, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"google-chrome", "chromium", "chromium-browser"},
	}

	path, err := finder.Find()
	if err != nil {
		return nil, err
	}

	return &Chrome{
		execPath: path,
		cleaningOpts: DefaultCleaningOptions(),
	}, nil
}

func (c *Chrome) Name() string {
	return "Chrome/Chromium"
}

func (c *Chrome) SetCleaningOptions(opts *CleaningOptions) {
	c.cleaningOpts = opts
}

func (c *Chrome) Fetch(url string) ([]byte, error) {
	// Use Chrome in headless mode to fetch content
	cmd := exec.Command(c.execPath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--enable-automation",
		"--virtual-time-budget=5000",  // Allow 5 seconds for JavaScript execution
		"--dump-dom",  // This will output the rendered DOM
		url,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("chrome execution error: %v", err)
	}

	return CleanHTML(output, c.cleaningOpts), nil
}
