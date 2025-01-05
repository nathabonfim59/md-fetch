package browser

import (
	"os/exec"
)

type Curl struct {
	execPath string
}

func NewCurl() (*Curl, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"curl"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &Curl{execPath: path}, nil
}

func (c *Curl) Name() string {
	return "Curl"
}

func (c *Curl) Fetch(url string) ([]byte, error) {
	cmd := exec.Command(c.execPath, "-L", "-s", url)
	content, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Strip JavaScript from the content since Curl can't execute it
	return StripJavaScript(content), nil
}
