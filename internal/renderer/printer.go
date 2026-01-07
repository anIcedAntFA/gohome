// Package renderer handles formatting and displaying commit data in various output formats.
package renderer

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// Config holds printer configuration options.
type Config struct {
	Format    string // "text" or "table"
	Style     string // "normal" or "markdown"
	ShowIcon  bool
	ShowScope bool
}

// Printer formats and outputs commit data according to configuration.
type Printer struct {
	cfg Config
}

// NewPrinter creates a new Printer instance with the given configuration.
func NewPrinter(cfg Config) *Printer {
	return &Printer{cfg: cfg}
}

// Print outputs formatted commit data to the provided writer.
func (p *Printer) Print(w io.Writer, repoName string, commits []entity.Commit) {
	if len(commits) == 0 {
		return
	}

	if p.cfg.Format == "table" {
		p.printTable(w, repoName, commits)
	} else {
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

// printTable outputs commits in table format.
func (p *Printer) printTable(w io.Writer, repoName string, commits []entity.Commit) {
	fmt.Fprintf(w, "\n📁 Repository: %s\n", repoName)

	// Initialize table with Options
	table := p.createTable(w, p.cfg.Style)

	// 1. Headers
	headers := []string{}
	if p.cfg.ShowIcon {
		headers = append(headers, "Icon")
	}
	headers = append(headers, "Type")
	if p.cfg.ShowScope {
		headers = append(headers, "Scope")
	}
	headers = append(headers, "Message")
	table.Header(headers)

	// 2. Data Rows
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
		_ = table.Append(row)
	}

	// 3. Render
	_ = table.Render()
	fmt.Fprintln(w)
}

// createTable initializes tablewriter.Table with Style configuration Options.
func (p *Printer) createTable(w io.Writer, style string) *tablewriter.Table {
	var options []tablewriter.Option

	// A. Configure Renderer (Interface)
	switch style {
	case "markdown":
		options = append(options, tablewriter.WithRenderer(renderer.NewMarkdown()))
	default: // "normal"
		// Use default Blueprint
		options = append(options, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{})))
	}

	// B. Configure Config (Alignment, Padding...)
	// We use functional option WithConfig to override default config
	conf := tablewriter.Config{
		Header: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignCenter},
		},
		Row: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		},
	}
	options = append(options, tablewriter.WithConfig(conf))

	// Create new table with writer and built options
	return tablewriter.NewTable(w, options...)
}

// PrintTasks outputs formatted task data to the provided writer.
func (p *Printer) PrintTasks(w io.Writer, tasks []entity.Task) {
	if len(tasks) == 0 {
		return
	}

	title := "📝 Additional Tasks"

	if p.cfg.Format == "table" {
		p.printTaskTable(w, title, tasks)
	} else {
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

func (p *Printer) printTaskTable(w io.Writer, title string, tasks []entity.Task) {
	fmt.Fprintf(w, "\n%s\n", title)

	// Tái sử dụng hàm createTable có sẵn
	table := p.createTable(w, p.cfg.Style)

	// Setup Header tương tự Commit Table
	headers := []string{}
	if p.cfg.ShowIcon {
		headers = append(headers, "Icon")
	}
	headers = append(headers, "Type", "Message")
	table.Header(headers)

	for _, t := range tasks {
		row := []string{}
		if p.cfg.ShowIcon {
			row = append(row, t.Icon)
		}
		row = append(row, t.Type, t.Message)
		_ = table.Append(row)
	}

	_ = table.Render()
	fmt.Fprintln(w)
}
