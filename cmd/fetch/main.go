package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	browserType := flag.String("browser", "curl", "Browser to use (chrome, firefox, links, lynx, w3m, curl)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: fetch [-browser=type] URL")
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
