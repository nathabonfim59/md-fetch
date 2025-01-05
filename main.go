package main

import (
	"fmt"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	url := "https://example.com"
	// Browser will be automatically selected based on availability
	content, err := fetcher.FetchContent(url, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Fetched Content:")
	fmt.Println(content)
}
