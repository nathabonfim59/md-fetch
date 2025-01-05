package browser

import (
	"fmt"
	"os/exec"
)

type Curl struct {
	execPath string
	cleaningOpts *CleaningOptions
}

func NewCurl() (*Curl, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"curl"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &Curl{
		execPath: path,
		cleaningOpts: DefaultCleaningOptions(),
	}, nil
}

func (c *Curl) Name() string {
	return "Curl"
}

func (c *Curl) SetCleaningOptions(opts *CleaningOptions) {
	c.cleaningOpts = opts
}

func (c *Curl) Fetch(url string) ([]byte, error) {
	cmd := exec.Command(c.execPath, "-L", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("curl execution error: %v", err)
	}

	return CleanHTML(output, c.cleaningOpts), nil
}
