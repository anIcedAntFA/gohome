// Package ui provides custom help templates and formatters for Cobra commands.
package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// CustomHelpTemplate returns a themed help template for Cobra commands.
// It uses lipgloss styling when a styled theme is active, otherwise falls back to plain text.
func CustomHelpTemplate(theme string) string {
	if !IsStyledTheme(theme) {
		// Return empty to use default Cobra help template
		return ""
	}

	// Custom themed help template
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// CustomUsageTemplate returns a themed usage template for Cobra commands.
func CustomUsageTemplate(theme string) string {
	if !IsStyledTheme(theme) {
		return ""
	}

	return `{{if .HasAvailableSubCommands}}{{"\n"}}{{bold "COMMANDS"}}
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{"\n"}}{{bold "FLAGS"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{"\n"}}{{bold "GLOBAL FLAGS"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasExample}}{{"\n"}}{{bold "EXAMPLES"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{"\n"}}Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// FormatHelp formats a Cobra command's help output with themed styling.
func FormatHelp(cmd *cobra.Command, theme string) string {
	if !IsStyledTheme(theme) {
		// Return default help without styling
		return cmd.UsageString()
	}

	ApplyTheme(theme)

	var b strings.Builder

	// Only show banner for root command, not subcommands
	if cmd.Parent() == nil {
		banner := Banner{
			Subtitle: "🏠 Git Activity Aggregator & Standup Report Generator",
			Author:   "",
			Date:     time.Time{},
			Version:  cmd.Root().Version,
		}
		b.WriteString(RenderBanner(banner))
		b.WriteString("\n\n")
	} else if cmd.Short != "" {
		// For subcommands, show Short description
		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Render(cmd.Short)
		b.WriteString(desc)
		b.WriteString("\n\n")
	}

	// Usage
	b.WriteString("\n")
	b.WriteString(SectionHeader("USAGE", 80))
	b.WriteString("\n\n")
	usage := lipgloss.NewStyle().
		Foreground(MutedColor).
		PaddingLeft(2).
		Render(fmt.Sprintf("%s %s", cmd.CommandPath(), cmd.UseLine()))
	b.WriteString(usage)
	b.WriteString("\n")

	// Subcommands
	if cmd.HasAvailableSubCommands() {
		b.WriteString("\n")
		b.WriteString(SectionHeader("COMMANDS", 80))
		b.WriteString("\n\n")
		for _, subCmd := range cmd.Commands() {
			if !subCmd.IsAvailableCommand() && subCmd.Name() != "help" {
				continue
			}
			cmdLine := lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true).
				Render(fmt.Sprintf("  %-15s", subCmd.Name()))
			cmdDesc := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(subCmd.Short)
			b.WriteString(cmdLine)
			b.WriteString(cmdDesc)
			b.WriteString("\n")
		}
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		b.WriteString("\n")
		b.WriteString(SectionHeader("FLAGS", 80))
		b.WriteString("\n\n")
		flagUsages := cmd.LocalFlags().FlagUsages()
		styledFlags := styleFlagUsages(flagUsages)
		b.WriteString(styledFlags)
	}

	// Global Flags
	if cmd.HasAvailableInheritedFlags() {
		b.WriteString("\n")
		b.WriteString(SectionHeader("GLOBAL FLAGS", 80))
		b.WriteString("\n\n")
		flagUsages := cmd.InheritedFlags().FlagUsages()
		styledFlags := styleFlagUsages(flagUsages)
		b.WriteString(styledFlags)
	}

	// Examples
	if cmd.Example != "" {
		b.WriteString(SectionHeader("EXAMPLES", 80))
		b.WriteString("\n")
		examples := lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingLeft(2).
			Render(strings.TrimSpace(cmd.Example))
		b.WriteString(examples)
		b.WriteString("\n\n")
	}

	// Footer
	if cmd.HasAvailableSubCommands() {
		b.WriteString("\n")
		footer := lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Render(fmt.Sprintf("Use \"%s [command] --help\" for more information about a command.", cmd.CommandPath()))
		b.WriteString(footer)
	}

	return b.String()
}

// styleFlagUsages styles flag usage text with colors and formatting.
func styleFlagUsages(usages string) string {
	if usages == "" {
		return ""
	}

	var b strings.Builder
	lines := strings.Split(usages, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse flag line: "  -f, --flag string   Description"
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)

		// Split into flag part and description
		parts := strings.SplitN(trimmed, "  ", 2)
		if len(parts) < 2 {
			b.WriteString(line)
			b.WriteString("\n")
			continue
		}

		flagPart := parts[0]
		descPart := strings.TrimSpace(parts[1])

		// Style flag name (bold, accent color)
		styledFlag := lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Render(flagPart)

		// Style description (muted)
		styledDesc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Render(descPart)

		// Reconstruct line with proper spacing
		b.WriteString(strings.Repeat(" ", indent))
		b.WriteString(styledFlag)
		b.WriteString("  ")
		b.WriteString(styledDesc)
		b.WriteString("\n")
	}

	return b.String()
}

// HelpFunc returns a custom help function for Cobra commands with themed styling.
func HelpFunc(theme string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, _ []string) {
		help := FormatHelp(cmd, theme)
		fmt.Fprint(cmd.OutOrStdout(), help)
	}
}

// UsageFunc returns a custom usage function for Cobra commands with themed styling.
func UsageFunc(theme string) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		help := FormatHelp(cmd, theme)
		fmt.Fprint(cmd.ErrOrStderr(), help)
		return nil
	}
}
