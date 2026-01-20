package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// TestNewPrinter verifies printer creation with configuration.
func TestNewPrinter(t *testing.T) {
	cfg := Config{
		Format:    "text",
		Style:     "normal",
		ShowIcon:  true,
		ShowScope: true,
	}

	printer := NewPrinter(cfg)

	if printer == nil {
		t.Fatal("NewPrinter() returned nil")
	}

	if printer.cfg.Format != "text" {
		t.Errorf("printer.cfg.Format = %q, want %q", printer.cfg.Format, "text")
	}

	if printer.cfg.Style != "normal" {
		t.Errorf("printer.cfg.Style = %q, want %q", printer.cfg.Style, "normal")
	}

	if !printer.cfg.ShowIcon {
		t.Error("printer.cfg.ShowIcon = false, want true")
	}

	if !printer.cfg.ShowScope {
		t.Error("printer.cfg.ShowScope = false, want true")
	}
}

// TestPrint verifies printing commits with different formats.
func TestPrint(t *testing.T) {
	commits := []entity.Commit{
		{
			Raw:     "feat(api): add user endpoint",
			Type:    "feat",
			Scope:   "api",
			Message: "add user endpoint",
			Icon:    "✨",
		},
		{
			Raw:     "fix: bug fix",
			Type:    "fix",
			Scope:   "-",
			Message: "bug fix",
			Icon:    "🐛",
		},
	}

	tests := []struct {
		name          string
		format        string
		style         string
		showIcon      bool
		showScope     bool
		repoName      string
		commits       []entity.Commit
		expectStrings []string
	}{
		{
			name:      "text_format_with_all_options",
			format:    "text",
			style:     "normal",
			showIcon:  true,
			showScope: true,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 Repository: test-repo",
				"- ✨ feat(api): add user endpoint",
				"- 🐛 fix(-): bug fix", // Scope "-" is included when ShowScope=true
				"------------------------------------------",
			},
		},
		{
			name:      "text_format_without_icon",
			format:    "text",
			style:     "normal",
			showIcon:  false,
			showScope: true,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 Repository: test-repo",
				"- feat(api): add user endpoint",
				"- fix(-): bug fix", // Scope "-" shown
			},
		},
		{
			name:      "text_format_without_scope",
			format:    "text",
			style:     "normal",
			showIcon:  true,
			showScope: false,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 Repository: test-repo",
				"- ✨ feat: add user endpoint",
				"- 🐛 fix: bug fix",
			},
		},
		{
			name:      "text_format_minimal",
			format:    "text",
			style:     "normal",
			showIcon:  false,
			showScope: false,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 Repository: test-repo",
				"- feat: add user endpoint",
				"- fix: bug fix",
			},
		},
		{
			name:      "table_format_with_all_options",
			format:    "table",
			style:     "normal",
			showIcon:  true,
			showScope: true,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 test-repo",
				"Icon",
				"Type",
				"Scope",
				"Message",
				"✨",
				"feat",
				"api",
				"add user endpoint",
			},
		},
		{
			name:      "table_format_markdown_style",
			format:    "table",
			style:     "markdown",
			showIcon:  true,
			showScope: true,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"📁 test-repo",
				"Icon",
				"Type",
				"Scope",
				"Message",
				"|------",
				"✨",
				"feat",
				"api",
				"add user endpoint",
			},
		},
		{
			name:      "table_format_without_icon",
			format:    "table",
			style:     "normal",
			showIcon:  false,
			showScope: true,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"Type",
				"Scope",
				"Message",
				"feat",
				"api",
				"add user endpoint",
			},
		},
		{
			name:      "table_format_without_scope",
			format:    "table",
			style:     "normal",
			showIcon:  true,
			showScope: false,
			repoName:  "test-repo",
			commits:   commits,
			expectStrings: []string{
				"Icon",
				"Type",
				"Message",
				"✨",
				"feat",
				"add user endpoint",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Format:    tt.format,
				Style:     tt.style,
				ShowIcon:  tt.showIcon,
				ShowScope: tt.showScope,
			}
			printer := NewPrinter(cfg)

			var buf bytes.Buffer
			printer.Print(&buf, tt.repoName, tt.commits)

			output := buf.String()

			for _, expectStr := range tt.expectStrings {
				if !strings.Contains(output, expectStr) {
					t.Errorf("Print() output missing expected string %q\nGot output:\n%s", expectStr, output)
				}
			}
		})
	}
}

// TestPrintEmptyCommits verifies that empty commit list produces no output.
func TestPrintEmptyCommits(t *testing.T) {
	cfg := Config{
		Format:    "text",
		Style:     "normal",
		ShowIcon:  true,
		ShowScope: true,
	}
	printer := NewPrinter(cfg)

	var buf bytes.Buffer
	printer.Print(&buf, "test-repo", []entity.Commit{})

	output := buf.String()
	if output != "" {
		t.Errorf("Print() with empty commits = %q, want empty string", output)
	}
}

// TestPrintTextFormat verifies text formatting logic.
func TestPrintTextFormat(t *testing.T) {
	commits := []entity.Commit{
		{
			Raw:     "chore: update dependencies",
			Type:    "chore",
			Scope:   "-",
			Message: "update dependencies",
			Icon:    "-",
		},
		{
			Raw:     "docs(readme): improve setup instructions",
			Type:    "docs",
			Scope:   "readme",
			Message: "improve setup instructions",
			Icon:    "📚",
		},
	}

	tests := []struct {
		name      string
		showIcon  bool
		showScope bool
		wantLines []string
	}{
		{
			name:      "all_options_enabled",
			showIcon:  true,
			showScope: true,
			wantLines: []string{
				"- - chore(-): update dependencies", // Scope (-) shown even when dash
				"- 📚 docs(readme): improve setup instructions",
			},
		},
		{
			name:      "icon_disabled",
			showIcon:  false,
			showScope: true,
			wantLines: []string{
				"- chore(-): update dependencies", // Scope (-) shown
				"- docs(readme): improve setup instructions",
			},
		},
		{
			name:      "scope_disabled",
			showIcon:  true,
			showScope: false,
			wantLines: []string{
				"- - chore: update dependencies",
				"- 📚 docs: improve setup instructions",
			},
		},
		{
			name:      "both_disabled",
			showIcon:  false,
			showScope: false,
			wantLines: []string{
				"- chore: update dependencies",
				"- docs: improve setup instructions",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Format:    "text",
				Style:     "normal",
				ShowIcon:  tt.showIcon,
				ShowScope: tt.showScope,
			}
			printer := NewPrinter(cfg)

			var buf bytes.Buffer
			printer.Print(&buf, "test-repo", commits)

			output := buf.String()

			for _, wantLine := range tt.wantLines {
				if !strings.Contains(output, wantLine) {
					t.Errorf("Print() output missing line %q\nGot:\n%s", wantLine, output)
				}
			}
		})
	}
}

// TestPrintTableFormat verifies table formatting with different styles.
func TestPrintTableFormat(t *testing.T) {
	commits := []entity.Commit{
		{
			Raw:     "feat: new feature",
			Type:    "feat",
			Scope:   "-",
			Message: "new feature",
			Icon:    "✨",
		},
	}

	tests := []struct {
		name         string
		style        string
		showIcon     bool
		showScope    bool
		wantInOutput []string
	}{
		{
			name:      "normal_style",
			style:     "normal",
			showIcon:  true,
			showScope: true,
			wantInOutput: []string{
				"Icon",
				"Type",
				"Scope",
				"Message",
				"✨",
				"feat",
			},
		},
		{
			name:      "markdown_style",
			style:     "markdown",
			showIcon:  true,
			showScope: true,
			wantInOutput: []string{
				"Icon",
				"Type",
				"Scope",
				"Message",
				"|------", // Markdown separator row
				"✨",
				"feat",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Format:    "table",
				Style:     tt.style,
				ShowIcon:  tt.showIcon,
				ShowScope: tt.showScope,
			}
			printer := NewPrinter(cfg)

			var buf bytes.Buffer
			printer.Print(&buf, "test-repo", commits)

			output := buf.String()

			for _, want := range tt.wantInOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Print() output missing %q\nGot:\n%s", want, output)
				}
			}
		})
	}
}

// TestPrintSpecialCharactersInMessage verifies handling of special characters.
func TestPrintSpecialCharactersInMessage(t *testing.T) {
	commits := []entity.Commit{
		{
			Raw:     "feat: add \"quotes\" and 'apostrophes'",
			Type:    "feat",
			Scope:   "-",
			Message: "add \"quotes\" and 'apostrophes'",
			Icon:    "✨",
		},
		{
			Raw:     "fix: handle & and < and >",
			Type:    "fix",
			Scope:   "-",
			Message: "handle & and < and >",
			Icon:    "🐛",
		},
		{
			Raw:     "chore: update \\path\\to\\file",
			Type:    "chore",
			Scope:   "-",
			Message: "update \\path\\to\\file",
			Icon:    "-",
		},
	}

	cfg := Config{
		Format:    "text",
		Style:     "normal",
		ShowIcon:  true,
		ShowScope: false,
	}
	printer := NewPrinter(cfg)

	var buf bytes.Buffer
	printer.Print(&buf, "test-repo", commits)

	output := buf.String()

	expectedStrings := []string{
		"add \"quotes\" and 'apostrophes'",
		"handle & and < and >",
		"update \\path\\to\\file",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Print() output missing special chars string %q\nGot:\n%s", expected, output)
		}
	}
}

// TestPrintLongMessages verifies handling of long commit messages.
func TestPrintLongMessages(t *testing.T) {
	commits := []entity.Commit{
		{
			Raw:     "feat: this is a very long commit message that exceeds normal line length and should still be displayed correctly without truncation",
			Type:    "feat",
			Scope:   "-",
			Message: "this is a very long commit message that exceeds normal line length and should still be displayed correctly without truncation",
			Icon:    "✨",
		},
	}

	cfg := Config{
		Format:    "text",
		Style:     "normal",
		ShowIcon:  false,
		ShowScope: false,
	}
	printer := NewPrinter(cfg)

	var buf bytes.Buffer
	printer.Print(&buf, "test-repo", commits)

	output := buf.String()

	if !strings.Contains(output, "this is a very long commit message") {
		t.Error("Print() did not handle long message correctly")
	}
}

// TestPrintTasks verifies task printing with different formats.
func TestPrintTasks(t *testing.T) {
	tasks := []entity.Task{
		{
			Type:    "todo",
			Message: "Review pull requests",
			Icon:    "📝",
			Enabled: true,
		},
		{
			Type:    "meeting",
			Message: "Team standup",
			Icon:    "👥",
			Enabled: true,
		},
	}

	tests := []struct {
		name          string
		format        string
		style         string
		showIcon      bool
		tasks         []entity.Task
		expectStrings []string
	}{
		{
			name:     "text_format_with_icon",
			format:   "text",
			style:    "normal",
			showIcon: true,
			tasks:    tasks,
			expectStrings: []string{
				"📝 Additional Tasks",
				"- 📝 todo: Review pull requests",
				"- 👥 meeting: Team standup",
				"------------------------------------------",
			},
		},
		{
			name:     "text_format_without_icon",
			format:   "text",
			style:    "normal",
			showIcon: false,
			tasks:    tasks,
			expectStrings: []string{
				"📝 Additional Tasks",
				"- todo: Review pull requests",
				"- meeting: Team standup",
			},
		},
		{
			name:     "table_format_with_icon",
			format:   "table",
			style:    "normal",
			showIcon: true,
			tasks:    tasks,
			expectStrings: []string{
				"📝 Additional Tasks",
				"Icon",
				"Type",
				"Message",
				"📝",
				"todo",
				"Review pull requests",
			},
		},
		{
			name:     "table_format_markdown",
			format:   "table",
			style:    "markdown",
			showIcon: true,
			tasks:    tasks,
			expectStrings: []string{
				"📝 Additional Tasks",
				"Icon",
				"Type",
				"Message",
				"|------",
				"📝",
				"todo",
				"Review pull requests",
			},
		},
		{
			name:     "table_format_without_icon",
			format:   "table",
			style:    "normal",
			showIcon: false,
			tasks:    tasks,
			expectStrings: []string{
				"Type",
				"Message",
				"todo",
				"Review pull requests",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Format:   tt.format,
				Style:    tt.style,
				ShowIcon: tt.showIcon,
			}
			printer := NewPrinter(cfg)

			var buf bytes.Buffer
			printer.PrintTasks(&buf, tt.tasks)

			output := buf.String()

			for _, expectStr := range tt.expectStrings {
				if !strings.Contains(output, expectStr) {
					t.Errorf("PrintTasks() output missing expected string %q\nGot output:\n%s", expectStr, output)
				}
			}
		})
	}
}

// TestPrintTasksEmpty verifies that empty task list produces no output.
func TestPrintTasksEmpty(t *testing.T) {
	cfg := Config{
		Format:   "text",
		Style:    "normal",
		ShowIcon: true,
	}
	printer := NewPrinter(cfg)

	var buf bytes.Buffer
	printer.PrintTasks(&buf, []entity.Task{})

	output := buf.String()
	if output != "" {
		t.Errorf("PrintTasks() with empty tasks = %q, want empty string", output)
	}
}

// TestPrintTasksWithEmptyFields verifies handling of tasks with missing fields.
func TestPrintTasksWithEmptyFields(t *testing.T) {
	tasks := []entity.Task{
		{
			Type:    "",
			Message: "Task without type",
			Icon:    "",
			Enabled: true,
		},
		{
			Type:    "reminder",
			Message: "Task without icon",
			Icon:    "",
			Enabled: true,
		},
	}

	cfg := Config{
		Format:   "text",
		Style:    "normal",
		ShowIcon: true,
	}
	printer := NewPrinter(cfg)

	var buf bytes.Buffer
	printer.PrintTasks(&buf, tasks)

	output := buf.String()

	// Should handle empty fields gracefully
	if !strings.Contains(output, "Task without type") {
		t.Error("PrintTasks() did not include task without type")
	}

	if !strings.Contains(output, "reminder: Task without icon") {
		t.Error("PrintTasks() did not include task without icon")
	}
}

// TestCreateTable verifies table creation with different styles.
// NOTE: This test was removed because createTable() was replaced with direct lipgloss table usage.

// TestPrintMultipleRepos verifies printing commits from multiple repositories.
func TestPrintMultipleRepos(t *testing.T) {
	cfg := Config{
		Format:    "text",
		Style:     "normal",
		ShowIcon:  false,
		ShowScope: false,
	}
	printer := NewPrinter(cfg)

	// Simulate multiple repo outputs
	commits1 := []entity.Commit{
		{Type: "feat", Scope: "-", Message: "repo1 feature", Icon: "✨"},
	}
	commits2 := []entity.Commit{
		{Type: "fix", Scope: "-", Message: "repo2 fix", Icon: "🐛"},
	}

	var buf bytes.Buffer
	printer.Print(&buf, "repo1", commits1)
	printer.Print(&buf, "repo2", commits2)

	output := buf.String()

	expectedStrings := []string{
		"📁 Repository: repo1",
		"- feat: repo1 feature",
		"📁 Repository: repo2",
		"- fix: repo2 fix",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Print() multiple repos output missing %q", expected)
		}
	}
}

// TestConfigVariations verifies all configuration combinations.
func TestConfigVariations(t *testing.T) {
	commits := []entity.Commit{
		{Type: "test", Scope: "unit", Message: "test message", Icon: "🧪"},
	}

	formats := []string{"text", "table"}
	styles := []string{"normal", "markdown"}
	boolOptions := []bool{true, false}

	for _, format := range formats {
		for _, style := range styles {
			for _, showIcon := range boolOptions {
				for _, showScope := range boolOptions {
					name := format + "_" + style + "_icon" + boolStr(showIcon) + "_scope" + boolStr(showScope)
					t.Run(name, func(t *testing.T) {
						cfg := Config{
							Format:    format,
							Style:     style,
							ShowIcon:  showIcon,
							ShowScope: showScope,
						}
						printer := NewPrinter(cfg)

						var buf bytes.Buffer
						printer.Print(&buf, "test-repo", commits)

						// Just verify it doesn't panic and produces output
						if buf.Len() == 0 {
							t.Error("Print() produced no output")
						}
					})
				}
			}
		}
	}
}

// boolStr converts bool to string for test names.
func boolStr(b bool) string {
	if b {
		return "T"
	}
	return "F"
}
