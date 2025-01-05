package main

import (
	"fmt"
	"os"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: md-fetch <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	// Prefer Chrome for better JavaScript support
	content, err := fetcher.FetchContent(url, "chrome")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Fetched Content:")
	fmt.Println(content)
}
