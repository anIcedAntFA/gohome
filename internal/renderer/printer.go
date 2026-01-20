// Package renderer handles formatting and displaying commit data in various output formats.
package renderer

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/ui"
)

// Config holds printer configuration options.
type Config struct {
	Format    string // "text" or "table"
	Style     string // "normal" or "markdown"
	Theme     string // "default", "dracula", "catppuccin-latte", "catppuccin-mocha"
	ShowIcon  bool
	ShowScope bool
}

// Printer formats and outputs commit data according to configuration.
type Printer struct {
	cfg Config
}

// NewPrinter creates a new Printer instance with the given configuration.
func NewPrinter(cfg Config) *Printer {
	// Apply default theme
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}

	// Apply theme if it's a styled theme (not default)
	if ui.IsStyledTheme(cfg.Theme) {
		ui.ApplyTheme(cfg.Theme)
	}

	return &Printer{cfg: cfg}
}

// Print outputs formatted commit data to the provided writer.
func (p *Printer) Print(w io.Writer, repoName string, commits []entity.Commit) {
	if len(commits) == 0 {
		return
	}

	// Use styled output if theme is not default
	if ui.IsStyledTheme(p.cfg.Theme) {
		p.printStyled(w, repoName, commits)
		return
	}

	// Use plain output (backward compatible)
	switch p.cfg.Format {
	case "table":
		p.printTable(w, repoName, commits)
	default:
		p.printText(w, repoName, commits)
	}
}

// printText outputs commits in plain text format.
func (p *Printer) printText(w io.Writer, repoName string, commits []entity.Commit) {
	fmt.Fprintf(w, "\n📁 Repository: %s\n", repoName)

	for _, c := range commits {
		line := "- "

		if p.cfg.ShowIcon {
			line += c.Icon + " "
		}

		line += c.Type

		if p.cfg.ShowScope && c.Scope != "" {
			line += "(" + c.Scope + ")"
		}

		line += ": " + c.Message

		fmt.Fprintln(w, line)
	}

	fmt.Fprintln(w, "------------------------------------------")
}

// printStyled outputs commits using lipgloss styling with the configured theme.
func (p *Printer) printStyled(w io.Writer, repoName string, commits []entity.Commit) {
	// Get terminal width for responsive layout
	width := 80 // Default width

	// Render commit list with modern styling
	output := ui.RenderCommitList(repoName, commits, p.cfg.ShowIcon, p.cfg.ShowScope, width)
	fmt.Fprintln(w, output)
}

// printTable outputs commits in table format using lipgloss table.
func (p *Printer) printTable(w io.Writer, repoName string, commits []entity.Commit) {
	fmt.Fprintf(w, "\n📁 %s\n\n", repoName)

	// Determine border style
	var border lipgloss.Border
	switch p.cfg.Style {
	case "markdown":
		border = lipgloss.MarkdownBorder()
	default: // "normal" or "ascii"
		border = lipgloss.NormalBorder()
	}

	// Create lipgloss table
	t := table.New().
		Border(border).
		BorderStyle(lipgloss.NewStyle().Foreground(ui.AccentColor))

	// Build headers
	headers := []string{}
	if p.cfg.ShowIcon {
		headers = append(headers, "Icon")
	}
	headers = append(headers, "Type")
	if p.cfg.ShowScope {
		headers = append(headers, "Scope")
	}
	headers = append(headers, "Message")
	t.Headers(headers...)

	// Add rows
	for _, c := range commits {
		row := []string{}
		if p.cfg.ShowIcon {
			row = append(row, c.Icon)
		}
		row = append(row, c.Type)
		if p.cfg.ShowScope {
			row = append(row, c.Scope)
		}
		row = append(row, c.Message)
		t.Row(row...)
	}

	// Apply style function for alternate row colors
	t.StyleFunc(func(row, _ int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 1)

		switch {
		case row == table.HeaderRow:
			// Header style
			return style.Bold(true).Foreground(ui.AccentColor).Align(lipgloss.Center)
		case row%2 == 0:
			// Even row
			return style.Foreground(lipgloss.Color("252"))
		default:
			// Odd row
			return style.Foreground(ui.MutedColor)
		}
	})

	// Render table
	fmt.Fprintln(w, t.Render())
}

// PrintTasks outputs formatted task data to the provided writer.
func (p *Printer) PrintTasks(w io.Writer, tasks []entity.Task) {
	if len(tasks) == 0 {
		return
	}

	title := "📝 Additional Tasks"

	// Use styled output if theme is not default
	if ui.IsStyledTheme(p.cfg.Theme) {
		p.printTasksStyled(w, title, tasks)
		return
	}

	// Use plain output (backward compatible)
	switch p.cfg.Format {
	case "table":
		p.printTaskTable(w, title, tasks)
	default:
		p.printTaskText(w, title, tasks)
	}
}

func (p *Printer) printTaskText(w io.Writer, title string, tasks []entity.Task) {
	fmt.Fprintf(w, "\n%s\n", title)
	for _, t := range tasks {
		// Format: - [Icon] Type: Message
		line := "- "
		if p.cfg.ShowIcon && t.Icon != "" {
			line += t.Icon + " "
		}
		if t.Type != "" {
			line += t.Type + ": "
		}
		line += t.Message
		fmt.Fprintln(w, line)
	}
	fmt.Fprintln(w, "------------------------------------------")
}

// printTaskTable outputs tasks in table format using lipgloss table.
func (p *Printer) printTaskTable(w io.Writer, title string, tasks []entity.Task) {
	fmt.Fprintf(w, "\n%s\n\n", title)

	// Determine border style
	var border lipgloss.Border
	switch p.cfg.Style {
	case "markdown":
		border = lipgloss.MarkdownBorder()
	default:
		border = lipgloss.NormalBorder()
	}

	// Create table
	t := table.New().
		Border(border).
		BorderStyle(lipgloss.NewStyle().Foreground(ui.AccentColor))

	// Setup headers
	headers := []string{}
	if p.cfg.ShowIcon {
		headers = append(headers, "Icon")
	}
	headers = append(headers, "Type", "Message")
	t.Headers(headers...)

	// Add rows
	for _, task := range tasks {
		row := []string{}
		if p.cfg.ShowIcon {
			row = append(row, task.Icon)
		}
		row = append(row, task.Type, task.Message)
		t.Row(row...)
	}

	// Apply style function
	t.StyleFunc(func(row, _ int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 1)

		if row == table.HeaderRow {
			return style.Bold(true).Foreground(ui.AccentColor).Align(lipgloss.Center)
		}

		if row%2 == 0 {
			return style.Foreground(lipgloss.Color("252"))
		}
		return style.Foreground(ui.MutedColor)
	})

	fmt.Fprintln(w, t.Render())
}

// printTasksStyled outputs tasks using lipgloss styling with the configured theme.
func (p *Printer) printTasksStyled(w io.Writer, title string, tasks []entity.Task) {
	// Render section header
	fmt.Fprintln(w)
	fmt.Fprintln(w, ui.SectionHeader(title, 80))

	// Render each task
	for _, t := range tasks {
		line := "  • "
		if p.cfg.ShowIcon && t.Icon != "" {
			line += t.Icon + " "
		}
		if t.Type != "" {
			line += t.Type + ": "
		}
		line += t.Message

		styledLine := ui.TaskStyle.Render(line)
		fmt.Fprintln(w, styledLine)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, ui.Separator(80, "━"))
}
