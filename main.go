package main

import (
	"fmt"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	url := "https://example.com"
	browserType := "curl" // Default to curl

	content, err := fetcher.FetchContent(url, browserType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Fetched Content:")
	fmt.Println(content)
}
