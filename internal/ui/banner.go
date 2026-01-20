package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Banner represents the application banner configuration.
type Banner struct {
	Title    string
	Subtitle string
	Author   string
	Version  string
	Date     time.Time
}

// RenderBanner renders the application banner with ASCII art logo.
func RenderBanner(b Banner) string {
	// Simpler ASCII art logo for GOHOME
	logo := `
   ____       _   _
  / ___| ___ | | | | ___  _ __ ___   ___
 | |  _ / _ \| |_| |/ _ \| '_ ' _ \ / _ \
 | |_| | (_) |  _  | (_) | | | | | |  __/
  \____|\___/|_| |_|\___/|_| |_| |_|\___|`

	// Style the logo
	styledLogo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render(logo)

	// Build subtitle with date and author
	var subtitle strings.Builder
	if b.Subtitle != "" {
		subtitle.WriteString(b.Subtitle)
	} else {
		subtitle.WriteString("Daily Standup Report")
	}

	if !b.Date.IsZero() {
		dateStr := b.Date.Format("2006-01-02")
		subtitle.WriteString(fmt.Sprintf("\n%s", dateStr))
	}

	if b.Author != "" {
		subtitle.WriteString(fmt.Sprintf(" • @%s", b.Author))
	}

	// Style the subtitle
	styledSubtitle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Align(lipgloss.Center).
		Render(subtitle.String())

	// Combine logo and subtitle
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styledLogo,
		styledSubtitle,
	)

	// Wrap in a border with reduced padding
	return BannerStyle.
		Align(lipgloss.Center).
		Render(content)
}

// RenderSimpleBanner renders a simpler banner without ASCII art.
func RenderSimpleBanner(title, subtitle string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title),
		"",
		subtitleStyle.Render(subtitle),
	)

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center).
		Render(content)
}

// RenderMinimalBanner renders a minimal single-line banner.
func RenderMinimalBanner(title string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Render("▶ " + title)
}
