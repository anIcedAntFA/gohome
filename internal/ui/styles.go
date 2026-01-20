// Package ui provides styled components for terminal output using lipgloss.
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette using adaptive colors for light/dark terminal support.
var (
	// PrimaryColor is the primary brand color used for headers and important elements.
	PrimaryColor = lipgloss.AdaptiveColor{
		Light: "#5A67D8", // Indigo-600
		Dark:  "#818CF8", // Indigo-400
	}

	// AccentColor is used for highlights and success states.
	AccentColor = lipgloss.AdaptiveColor{
		Light: "#059669", // Emerald-600
		Dark:  "#34D399", // Emerald-400
	}

	// MutedColor is used for secondary text and metadata.
	MutedColor = lipgloss.AdaptiveColor{
		Light: "#6B7280", // Gray-500
		Dark:  "#9CA3AF", // Gray-400
	}

	// DangerColor is used for errors and warnings.
	DangerColor = lipgloss.AdaptiveColor{
		Light: "#DC2626", // Red-600
		Dark:  "#F87171", // Red-400
	}
)

// CommitTypeColors maps commit types to their corresponding colors.
var CommitTypeColors = map[string]lipgloss.Color{
	"feat":     lipgloss.Color("#10B981"), // Green
	"fix":      lipgloss.Color("#EF4444"), // Red
	"docs":     lipgloss.Color("#3B82F6"), // Blue
	"style":    lipgloss.Color("#8B5CF6"), // Purple
	"refactor": lipgloss.Color("#F59E0B"), // Amber
	"perf":     lipgloss.Color("#EC4899"), // Pink
	"test":     lipgloss.Color("#06B6D4"), // Cyan
	"chore":    lipgloss.Color("#6B7280"), // Gray
	"ci":       lipgloss.Color("#14B8A6"), // Teal
	"build":    lipgloss.Color("#F97316"), // Orange
}

// Style definitions for various UI components.
var (
	// BannerStyle is used for the main banner/logo display.
	BannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			MarginBottom(1)

	// HeaderStyle is used for section headers (e.g., "📂 Repositories").
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true).
			BorderForeground(AccentColor).
			MarginTop(1).
			MarginBottom(1)

	// RepoNameStyle is used for repository names.
	RepoNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor).
			PaddingLeft(2)

	// CommitMessageStyle is used for commit messages.
	CommitMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")). // Light gray
				PaddingLeft(4)

	// CommitTypeStyle is used for commit type labels (feat, fix, etc.).
	CommitTypeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	// MetadataStyle is used for branch names, timestamps, and author info.
	MetadataStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	// TaskStyle is used for task items.
	TaskStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FBBF24")). // Amber
			PaddingLeft(2)

	// EmptyStateStyle is used when there are no results.
	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Align(lipgloss.Center)
)

// Theme represents a color theme configuration.
type Theme struct {
	Name    string
	Primary lipgloss.AdaptiveColor
	Accent  lipgloss.AdaptiveColor
	Muted   lipgloss.AdaptiveColor
	Danger  lipgloss.AdaptiveColor
}

// Predefined themes.
var (
	// DefaultTheme is the default classic theme (plain text, no styling).
	// This maintains backward compatibility with the original output.
	DefaultTheme = Theme{
		Name: "default",
		Primary: lipgloss.AdaptiveColor{
			Light: "#000000", // Black
			Dark:  "#FFFFFF", // White
		},
		Accent: lipgloss.AdaptiveColor{
			Light: "#000000", // Black
			Dark:  "#FFFFFF", // White
		},
		Muted: lipgloss.AdaptiveColor{
			Light: "#000000", // Black
			Dark:  "#FFFFFF", // White
		},
		Danger: lipgloss.AdaptiveColor{
			Light: "#000000", // Black
			Dark:  "#FFFFFF", // White
		},
	}

	// DraculaTheme is based on the popular Dracula color scheme.
	// https://draculatheme.com/
	DraculaTheme = Theme{
		Name: "dracula",
		Primary: lipgloss.AdaptiveColor{
			Light: "#6272a4", // Comment
			Dark:  "#bd93f9", // Purple
		},
		Accent: lipgloss.AdaptiveColor{
			Light: "#8be9fd", // Cyan
			Dark:  "#50fa7b", // Green
		},
		Muted: lipgloss.AdaptiveColor{
			Light: "#44475a", // Selection
			Dark:  "#6272a4", // Comment
		},
		Danger: lipgloss.AdaptiveColor{
			Light: "#ff5555", // Red
			Dark:  "#ff5555", // Red
		},
	}

	// CatppuccinLatteTheme is based on Catppuccin Latte (light theme).
	// https://github.com/catppuccin/catppuccin
	CatppuccinLatteTheme = Theme{
		Name: "catppuccin-latte",
		Primary: lipgloss.AdaptiveColor{
			Light: "#1e66f5", // Blue
			Dark:  "#1e66f5", // Blue
		},
		Accent: lipgloss.AdaptiveColor{
			Light: "#40a02b", // Green
			Dark:  "#40a02b", // Green
		},
		Muted: lipgloss.AdaptiveColor{
			Light: "#6c6f85", // Subtext0
			Dark:  "#6c6f85", // Subtext0
		},
		Danger: lipgloss.AdaptiveColor{
			Light: "#d20f39", // Red
			Dark:  "#d20f39", // Red
		},
	}

	// CatppuccinMochaTheme is based on Catppuccin Mocha (dark theme).
	// https://github.com/catppuccin/catppuccin
	CatppuccinMochaTheme = Theme{
		Name: "catppuccin-mocha",
		Primary: lipgloss.AdaptiveColor{
			Light: "#89b4fa", // Blue
			Dark:  "#89b4fa", // Blue
		},
		Accent: lipgloss.AdaptiveColor{
			Light: "#a6e3a1", // Green
			Dark:  "#a6e3a1", // Green
		},
		Muted: lipgloss.AdaptiveColor{
			Light: "#6c7086", // Overlay0
			Dark:  "#6c7086", // Overlay0
		},
		Danger: lipgloss.AdaptiveColor{
			Light: "#f38ba8", // Red
			Dark:  "#f38ba8", // Red
		},
	}
)

// Themes maps theme names to their configurations.
var Themes = map[string]Theme{
	"default":          DefaultTheme,
	"dracula":          DraculaTheme,
	"catppuccin-latte": CatppuccinLatteTheme,
	"catppuccin-mocha": CatppuccinMochaTheme,
}

// ApplyTheme applies the specified theme to the global styles.
func ApplyTheme(themeName string) {
	theme, ok := Themes[themeName]
	if !ok {
		theme = DefaultTheme
	}

	// Update global color variables
	PrimaryColor = theme.Primary
	AccentColor = theme.Accent
	MutedColor = theme.Muted
	DangerColor = theme.Danger

	// Rebuild styles with new colors
	BannerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor).
		Padding(1, 2).
		MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottom(true).
		BorderForeground(AccentColor).
		MarginTop(1).
		MarginBottom(1)

	RepoNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		PaddingLeft(2)

	CommitMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		PaddingLeft(4)

	CommitTypeStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	MetadataStyle = lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true)

	TaskStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FBBF24")).
		PaddingLeft(2)

	EmptyStateStyle = lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		Align(lipgloss.Center)
}

// GetCommitTypeColor returns the color for a given commit type.
func GetCommitTypeColor(commitType string) lipgloss.Color {
	if color, ok := CommitTypeColors[commitType]; ok {
		return color
	}
	return lipgloss.Color("#6B7280") // Default to gray
}

// IsStyledTheme returns true if the theme uses lipgloss styling.
// Default theme returns false to maintain backward compatibility.
func IsStyledTheme(themeName string) bool {
	return themeName != "" && themeName != "default"
}
