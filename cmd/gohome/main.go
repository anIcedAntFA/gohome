// Package main is the entry point for the gohome CLI application.
// It aggregates git commit reports and provides formatting options.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/anIcedAntFA/gohome/internal/config"
	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/anIcedAntFA/gohome/internal/parser"
	"github.com/anIcedAntFA/gohome/internal/renderer"
	"github.com/anIcedAntFA/gohome/internal/scanner"
	"github.com/anIcedAntFA/gohome/internal/sys"
)

func main() {
	// 1. Load Configurations
	cfg := config.Load()

	if cfg.SaveConfig {
		if err := cfg.SaveToFile(); err != nil {
			log.Fatalf("❌ Failed to save config: %v", err)
		}
		fmt.Println("✅ Configuration saved successfully!")
		fmt.Println("You can now run 'gohome' without flags to use these settings.")
		os.Exit(0) // Exit the program immediately, don't run scan
	}

	// 2. Initialize Dependencies
	gitClient := git.NewClient()
	parserSvc := parser.NewService()
	printer := renderer.NewPrinter(renderer.Config{
		Format:    cfg.OutputFmt,
		ShowIcon:  cfg.ShowIcon,
		ShowScope: cfg.ShowScope,
	})

	// 3. Detect author
	author := cfg.Author
	if author == "" {
		if val := gitClient.GetUser(); val != "" {
			author = val
		} else {
			log.Fatal("❌ Author not found. Please use -a flag or check git config.")
		}
	}

	// 4. Get Path and Period
	period := cfg.GetPeriod()
	fmt.Println("🗓️ Period:", period)
	absPath, _ := filepath.Abs(cfg.Path)

	// 5. Execute Git Log Retrieval
	repos, err := scanner.ScanGitRepos(absPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a buffer to collect the entire report content
	var clipboardBuffer bytes.Buffer

	// Determine where we will write output
	// Default is stdout (os.Stdout)
	var outputWriter io.Writer = os.Stdout

	// If user wants to copy, use MultiWriter
	// It will write simultaneously to both stdout AND buffer
	if cfg.CopyToClipboard {
		outputWriter = io.MultiWriter(os.Stdout, &clipboardBuffer)
	}

	foundAny := false
	for _, repo := range repos {
		// a. Get Raw Logs
		rawLogs, err := gitClient.GetLogs(repo, author, period)
		if err != nil || len(rawLogs) == 0 {
			continue
		}

		// b. Parse Logs -> Entities
		var commits []entity.Commit
		for _, line := range rawLogs {
			commits = append(commits, parserSvc.Parse(line))
		}

		// c. Render
		if len(commits) > 0 {
			foundAny = true
			printer.Print(outputWriter, filepath.Base(repo), commits)
		}
	}

	if !foundAny {
		fmt.Println("📭 No commits found.")
	} else {
		// If user wants to copy, copy from buffer to clipboard
		if cfg.CopyToClipboard {
			content := clipboardBuffer.String()
			if err := sys.CopyToClipboard(content); err != nil {
				fmt.Printf("\n⚠️  Failed to copy: %v\n", err)
				fmt.Println("   (Linux users: please install 'wl-clipboard' or 'xclip')")
			} else {
				fmt.Println("\n📋 Report copied to clipboard!")
			}
		}
	}
}
