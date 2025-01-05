package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nathabonfim59/md-fetch/internal/browser"
	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	browserType := flag.String("browser", "", "Browser to use (optional, defaults to chrome > firefox > curl)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Usage: fetch [-browser=type] URL\n\n")
		fmt.Printf("Supported browsers (in order of preference):\n")
		fmt.Printf("  %s\n\n", strings.Join(browser.DefaultBrowsers, ", "))
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)
	content, err := fetcher.FetchContent(url, *browserType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(content)
}
