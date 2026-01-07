// Package main is the entry point for the gohome CLI application.
// It aggregates git commit reports and provides formatting options.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/anIcedAntFA/gohome/internal/config"
	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/anIcedAntFA/gohome/internal/parser"
	"github.com/anIcedAntFA/gohome/internal/renderer"
	"github.com/anIcedAntFA/gohome/internal/scanner"
	"github.com/anIcedAntFA/gohome/internal/spinner"
	"github.com/anIcedAntFA/gohome/internal/sys"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Handle config save and exit early
	if cfg.SaveConfig {
		handleSaveConfig(cfg)
	}

	// 3. Initialize dependencies
	deps := initDependencies(cfg)

	// 4. Setup output writer
	outputWriter, clipboardBuffer := setupWriter(cfg.CopyToClipboard)

	// 5. Process and render
	foundAny := processAndRender(deps, cfg, outputWriter)

	// 6. Handle clipboard copy
	handleClipboard(foundAny, cfg.CopyToClipboard, clipboardBuffer)
}

// handleSaveConfig saves configuration to file and exits.
func handleSaveConfig(cfg *config.AppConfig) {
	if err := cfg.SaveToFile(); err != nil {
		log.Fatalf("❌ Failed to save config: %v", err)
	}

	configPath := config.GetConfigPath()
	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("💡 Tip: You can edit this file to customize your daily recurring tasks.\n   Config location: %s\n", configPath)
	fmt.Println("You can now run 'gohome' without flags to use these settings.")
	os.Exit(0)
}

// dependencies holds all service instances.
type dependencies struct {
	gitClient *git.Client
	parser    *parser.Service
	printer   *renderer.Printer
	author    string
	period    string
	repos     []string
}

// initDependencies creates and initializes all required services.
func initDependencies(cfg *config.AppConfig) *dependencies {
	gitClient := git.NewClient()
	parserSvc := parser.NewService()
	printer := renderer.NewPrinter(renderer.Config{
		Format:    cfg.OutputFmt,
		Style:     cfg.Preset,
		ShowIcon:  cfg.ShowIcon,
		ShowScope: cfg.ShowScope,
	})

	// Determine author
	author := cfg.Author
	if author == "" {
		if val := gitClient.GetUser(context.Background()); val != "" {
			author = val
		} else {
			log.Fatal("❌ Author not found. Please use -a flag or check git config.")
		}
	}

	// Get period and scan repos
	period := cfg.GetPeriod()
	fmt.Println("🗓️ Period:", period)

	absPath, _ := filepath.Abs(cfg.Path)

	sp := spinner.New("🔍 Scanning repositories...").
		WithFrames(spinner.PacmanGhost).
		WithInterval(100 * time.Millisecond)
	sp.Start()

	repos, err := scanner.ScanGitRepos(absPath)
	sp.Stop()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Found %d repositories\n", len(repos))

	return &dependencies{
		gitClient: gitClient,
		parser:    parserSvc,
		printer:   printer,
		author:    author,
		period:    period,
		repos:     repos,
	}
}

// setupWriter creates output writer and optional clipboard buffer.
func setupWriter(copyToClipboard bool) (io.Writer, *bytes.Buffer) {
	var clipboardBuffer bytes.Buffer
	var outputWriter io.Writer = os.Stdout

	if copyToClipboard {
		outputWriter = io.MultiWriter(os.Stdout, &clipboardBuffer)
	}

	return outputWriter, &clipboardBuffer
}

// processAndRender handles git commits and tasks rendering.
func processAndRender(deps *dependencies, cfg *config.AppConfig, w io.Writer) bool {
	foundCommits := processCommits(deps, w)
	foundTasks := processTasks(deps.printer, cfg, w)

	return foundCommits || foundTasks
}

// processCommits fetches and renders git commits.
func processCommits(deps *dependencies, w io.Writer) bool {
	foundAny := false

	for _, repo := range deps.repos {
		repoName := filepath.Base(repo)
		sp := spinner.New(fmt.Sprintf("📥 Fetching commits from %s...", repoName))
		sp.Start()

		rawLogs, err := deps.gitClient.GetLogs(context.Background(), repo, deps.author, deps.period)
		sp.Stop()

		if err != nil || len(rawLogs) == 0 {
			continue
		}

		var commits []entity.Commit
		for _, line := range rawLogs {
			commits = append(commits, deps.parser.Parse(line))
		}

		if len(commits) > 0 {
			foundAny = true
			deps.printer.Print(w, filepath.Base(repo), commits)
		}
	}

	return foundAny
}

// processTasks renders static (enabled only) and dynamic tasks.
func processTasks(printer *renderer.Printer, cfg *config.AppConfig, w io.Writer) bool {
	//nolint:prealloc // size unknown at compile time
	var activeTasks []entity.Task

	// 1. Filter Static Tasks: Only include tasks where Enabled = true
	for _, t := range cfg.Tasks {
		if t.Enabled {
			activeTasks = append(activeTasks, t)
		}
	}

	// 2. Add Dynamic Tasks: Tasks from CLI are always displayed by default
	for _, msg := range cfg.DynamicTasks {
		activeTasks = append(activeTasks, entity.Task{
			Message: msg,
			Type:    "misc",
			Icon:    "📌",
			// Enabled: true, // (Optional: set true if struct requires it, but typically not necessary for printing)
		})
	}

	// 3. Print if any tasks exist
	if len(activeTasks) > 0 {
		printer.PrintTasks(w, activeTasks)
		return true
	}

	return false
}

// handleClipboard copies content to clipboard if enabled.
func handleClipboard(foundAny, copyEnabled bool, buffer *bytes.Buffer) {
	if !foundAny {
		fmt.Println("📭 No commits or tasks found.")
		return
	}

	if copyEnabled {
		content := buffer.String()
		if err := sys.CopyToClipboard(context.Background(), content); err != nil {
			fmt.Printf("\n⚠️  Failed to copy: %v\n", err)
			fmt.Println("   (Linux users: please install 'wl-clipboard' or 'xclip')")
		} else {
			fmt.Println("\n📋 Report copied to clipboard!")
		}
	}
}
