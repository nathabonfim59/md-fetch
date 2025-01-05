package browser

import (
	"os/exec"
)

type Chrome struct {
	execPath string
}

func NewChrome() (*Chrome, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"google-chrome", "chromium", "chromium-browser"},
	}

	path, err := finder.Find()
	if err != nil {
		return nil, err
	}

	return &Chrome{execPath: path}, nil
}

func (c *Chrome) Name() string {
	return "Chrome/Chromium"
}

func (c *Chrome) Fetch(url string) ([]byte, error) {
	// Use Chrome in headless mode to fetch content
	cmd := exec.Command(c.execPath,
		"--headless",
		"--disable-gpu",
		"--dump-dom",
		"--no-sandbox",
		"--enable-automation",  // This flag helps with redirects and automation
		url,
	)

	return cmd.Output()
}
