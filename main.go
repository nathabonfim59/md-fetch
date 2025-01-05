package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
)

func main() {
	browserFlag := flag.String("browser", "", "Browser to use (chrome, firefox, or curl)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: md-fetch [-browser chrome|firefox|curl] <url>")
		os.Exit(1)
	}

	url := flag.Arg(0)
	content, err := fetcher.FetchContent(url, *browserFlag)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Fetched Content:")
	fmt.Println(content)
}
