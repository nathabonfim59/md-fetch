package browser

import (
	"os/exec"
)

// Links browser implementation
type Links struct {
	execPath string
}

func NewLinks() (*Links, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"links"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &Links{execPath: path}, nil
}

func (l *Links) Name() string {
	return "Links"
}

func (l *Links) Fetch(url string) ([]byte, error) {
	cmd := exec.Command(l.execPath, "-dump", url)
	return cmd.Output()
}

// Lynx browser implementation
type Lynx struct {
	execPath string
}

func NewLynx() (*Lynx, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"lynx"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &Lynx{execPath: path}, nil
}

func (l *Lynx) Name() string {
	return "Lynx"
}

func (l *Lynx) Fetch(url string) ([]byte, error) {
	cmd := exec.Command(l.execPath, "-dump", "-nolist", url)
	return cmd.Output()
}

// W3m browser implementation
type W3m struct {
	execPath string
}

func NewW3m() (*W3m, error) {
	finder := &DefaultExecutableFinder{
		names: []string{"w3m"},
	}
	
	path, err := finder.Find()
	if err != nil {
		return nil, err
	}
	
	return &W3m{execPath: path}, nil
}

func (w *W3m) Name() string {
	return "W3m"
}

func (w *W3m) Fetch(url string) ([]byte, error) {
	cmd := exec.Command(w.execPath, "-dump", url)
	return cmd.Output()
}
