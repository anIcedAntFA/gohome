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
	"github.com/anIcedAntFA/gohome/internal/editor"
	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/anIcedAntFA/gohome/internal/parser"
	"github.com/anIcedAntFA/gohome/internal/renderer"
	"github.com/anIcedAntFA/gohome/internal/scanner"
	"github.com/anIcedAntFA/gohome/internal/spinner"
	"github.com/anIcedAntFA/gohome/internal/sys"
)

var reportCmd = &cobra.Command{
	Use:           "report",
	Short:         "Generate activity report (default command)",
	Long:          `Scan repositories and generate a formatted activity report.`,
	Args:          cobra.NoArgs, // Reject any positional arguments
	SilenceErrors: true,         // We handle error formatting ourselves
	SilenceUsage:  false,        // Show usage on errors
	RunE:          runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Define flags as PERSISTENT on root command so they're inherited by all subcommands
	defineReportFlags(rootCmd)

	// Bind flags to viper ONCE
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Printf("Warning: failed to bind flags: %v\n", err)
	}

	// Make report the default command if no subcommand specified
	rootCmd.RunE = runReport
}

func defineReportFlags(cmd *cobra.Command) {
	// Time period flags
	cmd.PersistentFlags().IntP("hours", "H", 0, "Number of hours to look back")
	cmd.PersistentFlags().IntP("days", "d", 1, "Number of days to look back")
	cmd.PersistentFlags().IntP("weeks", "w", 0, "Number of weeks to look back")
	cmd.PersistentFlags().IntP("months", "M", 0, "Number of months to look back")
	cmd.PersistentFlags().IntP("years", "y", 0, "Number of years to look back")
	cmd.PersistentFlags().BoolP("today", "T", false, "Report from midnight to now")

	// Path and author
	cmd.PersistentFlags().StringP("path", "p", ".", "Root path to scan for repositories")
	cmd.PersistentFlags().IntP("max-depth", "m", 2, "Maximum depth to scan for repositories")
	cmd.PersistentFlags().StringP("author", "a", "", "Git author name (auto-detected if empty)")

	// Output formatting
	cmd.PersistentFlags().StringP("format", "f", "text", "Output format: text, table")
	cmd.PersistentFlags().StringP("style", "s", "normal", "Table style: normal, markdown")
	cmd.PersistentFlags().BoolP("icon", "i", false, "Show commit icon")
	cmd.PersistentFlags().BoolP("scope", "c", false, "Show commit scope")

	// Register dynamic flag completions
	_ = cmd.RegisterFlagCompletionFunc("format", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "table"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = cmd.RegisterFlagCompletionFunc("style", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"normal", "markdown"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Branch filtering
	cmd.PersistentFlags().BoolP("all-branches", "A", false, "Include commits from all local branches")
	cmd.PersistentFlags().StringP("branch", "b", "", "Filter commits by specific branch")

	// Clipboard and tasks
	cmd.PersistentFlags().BoolP("copy", "C", false, "Copy output to clipboard")
	cmd.PersistentFlags().StringSliceP("task", "t", []string{}, "Add custom task")

	// Editor flag
	cmd.PersistentFlags().BoolP("edit", "E", false, "Edit output in editor before displaying/copying")

	// Save config flag
	cmd.PersistentFlags().BoolP("save", "S", false, "Save current arguments as default configuration")

	// Flag groups: Mark mutually exclusive flags
	// User cannot use --all-branches and --branch together (conflicting branch filters)
	cmd.MarkFlagsMutuallyExclusive("all-branches", "branch")
}

// handlePeriodFlags applies period flag overrides to config.
func handlePeriodFlags(cmd *cobra.Command, cfg *viperconfig.Config) {
	periodFlags := []string{"hours", "days", "weeks", "months", "years", "today"}
	anyPeriodChanged := false
	for _, flag := range periodFlags {
		if cmd.Flags().Changed(flag) {
			anyPeriodChanged = true
			break
		}
	}

	if !anyPeriodChanged {
		return
	}

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

// handleSaveConfig saves config to file and prints success message.
func handleSaveConfig(cfg *viperconfig.Config) error {
	if err := cfg.SaveToFile(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
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

// getAuthorName determines the git author name from config or git client.
func getAuthorName(cfg *viperconfig.Config, gitClient *git.Client) (string, error) {
	author := cfg.Author
	if author != "" {
		return author, nil
	}

	if val := gitClient.GetUser(context.Background()); val != "" {
		return val, nil
	}

	return "", fmt.Errorf("author not found. Please use -a flag or check git config")
}

// scanRepositories scans for git repos and displays results.
func scanRepositories(absPath string, maxDepth int) ([]string, error) {
	sp := spinner.New("🔍 Scanning repositories...").
		WithFrames(spinner.PacmanGhost).
		WithInterval(100 * time.Millisecond)
	sp.Start()

	repos, err := scanner.ScanGitRepos(absPath, maxDepth)
	sp.Stop()

	if err != nil {
		return nil, err
	}

	fmt.Printf("✓ Found %d repositories\n", len(repos))
	return repos, nil
}

// processCommits fetches and prints commits from repositories.
func processCommits(repos []string, author, period string, cfg *viperconfig.Config, gitClient *git.Client, parserSvc *parser.Service, printer *renderer.Printer, writer io.Writer) bool {
	foundAny := false
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
			printer.Print(writer, filepath.Base(repo), commits)
		}
	}
	return foundAny
}

// collectActiveTasks gathers enabled static tasks and all dynamic tasks.
func collectActiveTasks(cfg *viperconfig.Config) []entity.Task {
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

	return activeTasks
}

// handleClipboard copies content to clipboard and displays result.
func handleClipboard(content string) {
	if err := sys.CopyToClipboard(context.Background(), content); err != nil {
		fmt.Printf("\n⚠️  Failed to copy: %v\n", err)
		fmt.Println("   (Linux users: please install 'wl-clipboard' or 'xclip')")
	} else {
		fmt.Println("\n📋 Report copied to clipboard!")
	}
}

// handleEditMode opens the content in editor and returns the edited version.
func handleEditMode(content string, cfg *viperconfig.Config) (string, error) {
	// Create editor client (respects VISUAL/EDITOR env or uses config)
	var editorClient *editor.Client
	if cfg.EditorCommand != "" {
		editorClient = editor.NewClientWithEditor(cfg.EditorCommand)
	} else {
		editorClient = editor.NewClient()
	}

	// Show brief message before opening editor
	fmt.Printf("\n📝 Opening editor (%s)... Save and close when done.\n", editorClient.GetEditor())
	fmt.Println("💡 Tip: Delete all content or close without saving to cancel")

	// Brief pause to let user read the message
	time.Sleep(500 * time.Millisecond)

	// Open editor and get modified content
	editedContent, err := editorClient.Open(content)
	if err != nil {
		return "", err
	}

	// Show preview of edited content
	fmt.Println("\n✅ Edited Report:")
	fmt.Println(editedContent)

	return editedContent, nil
}

func runReport(cmd *cobra.Command, _ []string) error {
	// Load configuration (Viper handles flag > env > config > default)
	cfg := viperconfig.LoadFromViper()

	// Handle period flag overrides
	handlePeriodFlags(cmd, cfg)

	// Get dynamic tasks from flag
	cfg.DynamicTasks = viper.GetStringSlice("task")

	// Validate configuration (separate validation layer)
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Handle save config
	if viper.GetBool("save") {
		return handleSaveConfig(cfg)
	}

	// Validate format and style compatibility
	if cfg.Format != "table" && cfg.Style != "normal" {
		return fmt.Errorf("--style flag only works with --format table\n   Current: --format %s --style %s\n   Hint: Use '--format table --style %s' or remove --style flag", cfg.Format, cfg.Style, cfg.Style)
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
	author, err := getAuthorName(cfg, gitClient)
	if err != nil {
		return err
	}

	// Get period and scan repos
	period := cfg.GetPeriod()
	fmt.Println("🗓️  Period:", period)

	absPath, _ := filepath.Abs(cfg.Path)
	repos, err := scanRepositories(absPath, cfg.MaxDepth)
	if err != nil {
		return err
	}

	// Setup buffer for output generation
	var outputBuffer bytes.Buffer

	// Process commits and tasks to buffer (don't print to stdout yet if edit mode is on)
	foundCommits := processCommits(repos, author, period, cfg, gitClient, parserSvc, printer, &outputBuffer)

	activeTasks := collectActiveTasks(cfg)
	if len(activeTasks) > 0 {
		printer.PrintTasks(&outputBuffer, activeTasks)
	}

	// Check if anything was found
	foundAny := foundCommits || len(activeTasks) > 0
	if !foundAny {
		fmt.Println("📭 No commits or tasks found.")
		return nil
	}

	// Get the generated output
	output := outputBuffer.String()

	// Handle edit mode (if enabled)
	editMode := viper.GetBool("edit")
	if editMode {
		editedOutput, err := handleEditMode(output, cfg)
		if err != nil {
			return fmt.Errorf("edit mode failed: %w", err)
		}
		output = editedOutput
	} else {
		// Only print to stdout if not in edit mode (edit mode prints after editing)
		fmt.Print(output)
	}

	// Handle clipboard (use edited output if edit mode was used)
	if cfg.CopyToClipboard {
		handleClipboard(output)
	}

	return nil
}
