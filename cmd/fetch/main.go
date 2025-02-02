package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/nathabonfim59/md-fetch/internal/browser"
	"github.com/nathabonfim59/md-fetch/internal/fetcher"
	"github.com/nathabonfim59/md-fetch/internal/server"
	"github.com/spf13/cobra"
)

var (
	browserType string
	save        bool
	filename    string
	port        int
)

var rootCmd = &cobra.Command{
	Use:   "md-fetch [url]",
	Short: "Fetch web content and convert it to Markdown",
	Long: `A CLI tool that fetches web content and converts it to clean, readable Markdown format.
Supports multiple browsers and can bypass anti-scraping measures.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		content, err := fetcher.FetchContent(url, browserType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if save {
			if filename == "" {
				filename = slug.Make(url) + ".md"
			}

			fmt.Printf("Saving content to %s\n", filename)
			file, err := os.Create(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()

			if _, err := file.WriteString(content); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(content)
		}
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server mode",
	Long: `Start md-fetch in HTTP server mode. This provides a REST API for fetching content
from multiple URLs in parallel.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(port)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Root command flags
	rootCmd.Flags().StringVarP(&browserType, "browser", "b", "", fmt.Sprintf("Browser to use (optional, defaults to %s)", strings.Join(browser.DefaultBrowsers, " > ")))
	rootCmd.Flags().BoolVarP(&save, "save", "s", false, "Save content to a file with slugified URL name")
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "", "Custom filename to save the content (optional, defaults to slugified URL)")

	// Server command flags
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for HTTP server")
	rootCmd.AddCommand(serveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
