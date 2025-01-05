package browser

import (
	"fmt"
	"os/exec"
)

type Firefox struct {
	execPath string
}

func NewFirefox() (*Firefox, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"firefox"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &Firefox{execPath: path}, nil
}

func (f *Firefox) Name() string {
	return "Firefox"
}

func (f *Firefox) Fetch(url string) ([]byte, error) {
	// Use Firefox in headless mode to fetch content
	cmd := exec.Command(f.execPath,
		"--headless",
		"--enable-automation",
		"--wait-for-browser",
		"--dump-dom",  // This will output the rendered DOM
		url,
	)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("firefox execution error: %v", err)
	}

	return output, nil
}
