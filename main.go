package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nathabonfim59/md-fetch/internal/fetcher"
	"github.com/nathabonfim59/md-fetch/internal/server"
)

func main() {
	// Command-line flags
	browserFlag := flag.String("browser", "", "Browser to use (chrome, firefox, or curl)")
	serveFlag := flag.Bool("serve", false, "Start HTTP server")
	portFlag := flag.Int("port", 8080, "Port for HTTP server (when using --serve)")
	flag.Parse()

	// Server mode
	if *serveFlag {
		srv := server.New(*portFlag)
		if err := srv.Start(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// CLI mode
	if flag.NArg() < 1 {
		fmt.Println("Usage:")
		fmt.Println("  md-fetch [-browser chrome|firefox|curl] <url>")
		fmt.Println("  md-fetch --serve [-port 8080]")
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
