package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// RepoCard represents a styled repository card.
type RepoCard struct {
	Name        string
	Branch      string
	CommitCount int
	TimeRange   string
	ShowBorder  bool
	MaxWidth    int
}

// Render renders the repository card with styling.
func (rc RepoCard) Render() string {
	// Build metadata line
	metadata := fmt.Sprintf("📍 %s", rc.Branch)
	if rc.CommitCount > 0 {
		metadata += fmt.Sprintf(" • iconhere %d commit", rc.CommitCount)
		if rc.CommitCount > 1 {
			metadata += "s"
		}
	}
	if rc.TimeRange != "" {
		metadata += fmt.Sprintf(" • iconhere %s", rc.TimeRange)
	}

	metadataStyled := MetadataStyle.Render(metadata)

	// Combine name and metadata
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		metadataStyled,
	)

	// Card with border
	if rc.ShowBorder {
		cardStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(0, 1)

		if rc.MaxWidth > 0 {
			cardStyle = cardStyle.Width(rc.MaxWidth - 4) // Account for borders
		}

		// Add repo name as title in border
		title := fmt.Sprintf("─ %s ", rc.Name)
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(title)+strings.Repeat("─", maxInt(0, rc.MaxWidth-len(rc.Name)-4)),
			cardStyle.Render(content),
		)
	}

	// Simple format without border
	nameStyled := RepoNameStyle.Render(rc.Name)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		nameStyled,
		metadataStyled,
	)
}

// CommitItem represents a styled commit item.
type CommitItem struct {
	Type    string
	Scope   string
	Message string
	Icon    string
}

// Render renders a single commit with styling.
func (ci CommitItem) Render() string {
	var parts []string

	// Start with bullet point
	parts = append(parts, "•")

	// Add icon if present
	if ci.Icon != "" {
		parts = append(parts, ci.Icon)
	}

	// Add type with color
	if ci.Type != "" {
		typeColor := GetCommitTypeColor(ci.Type)
		typeStyled := lipgloss.NewStyle().
			Bold(true).
			Foreground(typeColor).
			Render(ci.Type)
		parts = append(parts, typeStyled)
	}

	// Add scope if present
	if ci.Scope != "" {
		scopeStyled := lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(fmt.Sprintf("(%s)", ci.Scope))
		parts = append(parts, scopeStyled)
	}

	// Add message
	messageStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Render(ci.Message)
	parts = append(parts, messageStyled)

	// Join all parts with single space
	return "  " + strings.Join(parts, " ")
}

// RenderCommitList renders a list of commits for a repository with a bordered card layout.
func RenderCommitList(repoName string, commits []entity.Commit, showIcon, showScope bool, maxWidth int) string {
	if len(commits) == 0 {
		return ""
	}

	var content strings.Builder

	// Repository header with metadata
	headerLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		Render(fmt.Sprintf("📂 %s", repoName))

	metadata := fmt.Sprintf("%d commit", len(commits))
	if len(commits) > 1 {
		metadata += "s"
	}
	metadataLine := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render(metadata)

	content.WriteString(headerLine)
	content.WriteString("  ")
	content.WriteString(metadataLine)
	content.WriteString("\n\n")

	// Render each commit with compact styling
	for i, commit := range commits {
		item := CommitItem{
			Type:    commit.Type,
			Scope:   commit.Scope,
			Message: commit.Message,
			Icon:    commit.Icon,
		}

		if !showIcon {
			item.Icon = ""
		}
		if !showScope {
			item.Scope = ""
		}

		content.WriteString(item.Render())
		if i < len(commits)-1 {
			content.WriteString("\n")
		}
	}

	// Wrap everything in a bordered card
	cardStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor).
		Padding(1, 2)

	if maxWidth > 0 {
		cardStyle = cardStyle.Width(maxWidth - 6) // Account for borders and padding
	}

	return cardStyle.Render(content.String())
}

// SectionHeader renders a styled section header.
func SectionHeader(title string, width int) string {
	style := HeaderStyle
	if width > 0 {
		style = style.Width(width)
	}
	return style.Render(title)
}

// Separator renders a horizontal separator line.
func Separator(width int, char string) string {
	if char == "" {
		char = "━"
	}
	line := strings.Repeat(char, width)
	return lipgloss.NewStyle().
		Foreground(AccentColor).
		Render(line)
}

// maxInt returns the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
