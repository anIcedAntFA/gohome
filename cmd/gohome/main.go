package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/anIcedAntFA/gohome/internal/config"
	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/anIcedAntFA/gohome/internal/parser"
	"github.com/anIcedAntFA/gohome/internal/renderer"
	"github.com/anIcedAntFA/gohome/internal/scanner"
)

func main() {
	// 1. Load Configurations
	cfg := config.Load()

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
	fmt.Println("Period:", period)
	absPath, _ := filepath.Abs(cfg.Path)

	// 5. Execute Git Log Retrieval
	repos, err := scanner.ScanGitRepos(absPath)
	if err != nil {
		log.Fatal(err)
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
			printer.Print(filepath.Base(repo), commits)
		}
	}

	if !foundAny {
		fmt.Println("📭 No commits found.")
	}
}
