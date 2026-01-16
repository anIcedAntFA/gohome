package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	viperconfig "github.com/anIcedAntFA/gohome/internal/config/viper"
	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/anIcedAntFA/gohome/internal/parser"
	"github.com/anIcedAntFA/gohome/internal/renderer"
	"github.com/anIcedAntFA/gohome/internal/scanner"
	"github.com/anIcedAntFA/gohome/internal/spinner"
	"github.com/anIcedAntFA/gohome/internal/sys"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate activity report (default command)",
	Long:  `Scan repositories and generate a formatted activity report.`,
	RunE:  runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Define flags as PERSISTENT on root command so they're inherited by all subcommands
	defineReportFlags(rootCmd)

	// Bind flags to viper ONCE
	viper.BindPFlags(rootCmd.PersistentFlags())

	// Make report the default command if no subcommand specified
	rootCmd.RunE = runReport
}

func defineReportFlags(cmd *cobra.Command) {
	// Time period flags
	cmd.PersistentFlags().IntP("hours", "H", 0, "Number of hours to look back")
	cmd.PersistentFlags().IntP("days", "d", 1, "Number of days to look back")
	cmd.PersistentFlags().IntP("weeks", "w", 0, "Number of weeks to look back")
	cmd.PersistentFlags().IntP("months", "m", 0, "Number of months to look back")
	cmd.PersistentFlags().IntP("years", "y", 0, "Number of years to look back")
	cmd.PersistentFlags().Bool("today", false, "Report from midnight to now")

	// Path and author
	cmd.PersistentFlags().StringP("path", "p", ".", "Root path to scan for repositories")
	cmd.PersistentFlags().Int("max-depth", 2, "Maximum depth to scan for repositories")
	cmd.PersistentFlags().StringP("author", "a", "", "Git author name (auto-detected if empty)")

	// Output formatting
	cmd.PersistentFlags().StringP("format", "f", "text", "Output format: text, table")
	cmd.PersistentFlags().StringP("style", "s", "normal", "Table style: normal, markdown")
	cmd.PersistentFlags().BoolP("icon", "i", false, "Show commit icon")
	cmd.PersistentFlags().BoolP("scope", "c", false, "Show commit scope")

	// Branch filtering
	cmd.PersistentFlags().BoolP("all-branches", "b", false, "Include commits from all local branches")
	cmd.PersistentFlags().String("branch", "", "Filter commits by specific branch")

	// Clipboard and tasks
	cmd.PersistentFlags().BoolP("copy", "C", false, "Copy output to clipboard")
	cmd.PersistentFlags().StringSliceP("task", "t", []string{}, "Add custom task")

	// Save config flag
	cmd.PersistentFlags().Bool("save", false, "Save current arguments as default configuration")
}

func runReport(cmd *cobra.Command, args []string) error {
	// Load configuration (Viper handles flag > env > config > default)
	cfg := viperconfig.LoadFromViper()

	// IMPORTANT: Handle period flag overrides
	// If user explicitly sets any period flag via CLI, it should override ALL config period values
	// This prevents conflicts like: config has today=true, but user sets -d 3
	periodFlags := []string{"hours", "days", "weeks", "months", "years", "today"}
	anyPeriodChanged := false
	for _, flag := range periodFlags {
		if cmd.Flags().Changed(flag) {
			anyPeriodChanged = true
			break
		}
	}

	if anyPeriodChanged {
		// Reset all period values to 0/false first (clear config values)
		cfg.Hours = 0
		cfg.Days = 0
		cfg.Weeks = 0
		cfg.Months = 0
		cfg.Years = 0
		cfg.Today = false

		// Then set only the flags that were explicitly provided by user
		if cmd.Flags().Changed("hours") {
			cfg.Hours = viper.GetInt("hours")
		}
		if cmd.Flags().Changed("days") {
			cfg.Days = viper.GetInt("days")
		}
		if cmd.Flags().Changed("weeks") {
			cfg.Weeks = viper.GetInt("weeks")
		}
		if cmd.Flags().Changed("months") {
			cfg.Months = viper.GetInt("months")
		}
		if cmd.Flags().Changed("years") {
			cfg.Years = viper.GetInt("years")
		}
		if cmd.Flags().Changed("today") {
			cfg.Today = viper.GetBool("today")
		}

		// Normalize to keep only the highest priority
		cfg.NormalizePeriod()
	}

	// Get dynamic tasks from flag
	cfg.DynamicTasks = viper.GetStringSlice("task")

	// Handle save config
	if viper.GetBool("save") {
		if err := cfg.SaveToFile(); err != nil {
			return fmt.Errorf("❌ Failed to save config: %w", err)
		}
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			home, _ := os.UserHomeDir()
			configPath = filepath.Join(home, ".gohome.json")
		}
		fmt.Println("✅ Configuration saved successfully!")
		fmt.Printf("💡 Tip: You can edit this file to customize your daily recurring tasks.\n   Config location: %s\n", configPath)
		fmt.Println("You can now run 'gohome' without flags to use these settings.")
		return nil
	}

	// Validate format and style compatibility
	if cfg.Format != "table" && cfg.Style != "normal" {
		return fmt.Errorf("❌ Error: --style flag only works with --format table\n   Current: --format %s --style %s\n   Hint: Use '--format table --style %s' or remove --style flag", cfg.Format, cfg.Style, cfg.Style)
	}

	// Initialize dependencies
	gitClient := git.NewClient()
	parserSvc := parser.NewService()
	printer := renderer.NewPrinter(renderer.Config{
		Format:    cfg.Format,
		Style:     cfg.Style,
		ShowIcon:  cfg.ShowIcon,
		ShowScope: cfg.ShowScope,
	})

	// Determine author
	author := cfg.Author
	if author == "" {
		if val := gitClient.GetUser(context.Background()); val != "" {
			author = val
		} else {
			return fmt.Errorf("❌ Author not found. Please use -a flag or check git config")
		}
	}

	// Get period and scan repos
	period := cfg.GetPeriod()
	fmt.Println("🗓️  Period:", period)

	absPath, _ := filepath.Abs(cfg.Path)

	sp := spinner.New("🔍 Scanning repositories...").
		WithFrames(spinner.PacmanGhost).
		WithInterval(100 * time.Millisecond)
	sp.Start()

	repos, err := scanner.ScanGitRepos(absPath, cfg.MaxDepth)
	sp.Stop()

	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d repositories\n", len(repos))

	// Setup writer for clipboard
	var clipboardBuffer bytes.Buffer
	var outputWriter io.Writer = os.Stdout

	if cfg.CopyToClipboard {
		outputWriter = io.MultiWriter(os.Stdout, &clipboardBuffer)
	}

	// Process commits and tasks
	foundAny := false

	// Process commits
	for _, repo := range repos {
		repoName := filepath.Base(repo)
		sp := spinner.New(fmt.Sprintf("📥 Fetching commits from %s...", repoName))
		sp.Start()

		rawLogs, err := gitClient.GetLogs(context.Background(), repo, author, period, cfg.AllBranches, cfg.Branch)
		sp.Stop()

		if err != nil || len(rawLogs) == 0 {
			continue
		}

		commits := make([]entity.Commit, 0, len(rawLogs))
		for _, line := range rawLogs {
			commits = append(commits, parserSvc.Parse(line))
		}

		if len(commits) > 0 {
			foundAny = true
			printer.Print(outputWriter, filepath.Base(repo), commits)
		}
	}

	// Process tasks
	activeTasks := make([]entity.Task, 0, len(cfg.Tasks))

	// Filter static tasks (only enabled ones)
	for _, t := range cfg.Tasks {
		if t.Enabled {
			activeTasks = append(activeTasks, t)
		}
	}

	// Add dynamic tasks (always displayed)
	for _, msg := range cfg.DynamicTasks {
		activeTasks = append(activeTasks, entity.Task{
			Message: msg,
			Type:    "misc",
			Icon:    "📌",
		})
	}

	if len(activeTasks) > 0 {
		printer.PrintTasks(outputWriter, activeTasks)
		foundAny = true
	}

	// Handle clipboard
	if !foundAny {
		fmt.Println("📭 No commits or tasks found.")
		return nil
	}

	if cfg.CopyToClipboard {
		content := clipboardBuffer.String()
		if err := sys.CopyToClipboard(context.Background(), content); err != nil {
			fmt.Printf("\n⚠️  Failed to copy: %v\n", err)
			fmt.Println("   (Linux users: please install 'wl-clipboard' or 'xclip')")
		} else {
			fmt.Println("\n📋 Report copied to clipboard!")
		}
	}

	return nil
}
