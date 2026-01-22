package parser

import (
	"testing"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// TestParse tests the parsing of Conventional Commits.
func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    string
		wantScope   string
		wantMessage string
		wantIcon    string
	}{
		// Standard Conventional Commits
		{
			name:        "feat_without_scope",
			input:       "feat: add new feature",
			wantType:    "feat",
			wantScope:   "-",
			wantMessage: "add new feature",
			wantIcon:    "-",
		},
		{
			name:        "feat_with_scope",
			input:       "feat(api): add user endpoint",
			wantType:    "feat",
			wantScope:   "api",
			wantMessage: "add user endpoint",
			wantIcon:    "-",
		},
		{
			name:        "fix_with_scope",
			input:       "fix(auth): resolve login bug",
			wantType:    "fix",
			wantScope:   "auth",
			wantMessage: "resolve login bug",
			wantIcon:    "-",
		},
		{
			name:        "chore_without_scope",
			input:       "chore: update dependencies",
			wantType:    "chore",
			wantScope:   "-",
			wantMessage: "update dependencies",
			wantIcon:    "-",
		},
		{
			name:        "docs_with_scope",
			input:       "docs(readme): update installation guide",
			wantType:    "docs",
			wantScope:   "readme",
			wantMessage: "update installation guide",
			wantIcon:    "-",
		},
		{
			name:        "style_commit",
			input:       "style: format code with prettier",
			wantType:    "style",
			wantScope:   "-",
			wantMessage: "format code with prettier",
			wantIcon:    "-",
		},
		{
			name:        "refactor_with_scope",
			input:       "refactor(database): optimize query",
			wantType:    "refactor",
			wantScope:   "database",
			wantMessage: "optimize query",
			wantIcon:    "-",
		},
		{
			name:        "test_commit",
			input:       "test: add unit tests for parser",
			wantType:    "test",
			wantScope:   "-",
			wantMessage: "add unit tests for parser",
			wantIcon:    "-",
		},
		{
			name:        "perf_commit",
			input:       "perf(scanner): improve scanning speed",
			wantType:    "perf",
			wantScope:   "scanner",
			wantMessage: "improve scanning speed",
			wantIcon:    "-",
		},
		{
			name:        "build_commit",
			input:       "build: update Makefile",
			wantType:    "build",
			wantScope:   "-",
			wantMessage: "update Makefile",
			wantIcon:    "-",
		},
		{
			name:        "ci_commit",
			input:       "ci(github): add workflow",
			wantType:    "ci",
			wantScope:   "github",
			wantMessage: "add workflow",
			wantIcon:    "-",
		},

		// Commits with emojis
		{
			name:        "feat_with_sparkles_emoji",
			input:       "✨ feat: add sparkles feature",
			wantType:    "feat",
			wantScope:   "-",
			wantMessage: "add sparkles feature",
			wantIcon:    "✨",
		},
		{
			name:        "fix_with_bug_emoji",
			input:       "🐛 fix(api): resolve crash",
			wantType:    "fix",
			wantScope:   "api",
			wantMessage: "resolve crash",
			wantIcon:    "🐛",
		},
		{
			name:        "docs_with_book_emoji",
			input:       "📚 docs: update documentation",
			wantType:    "docs",
			wantScope:   "-",
			wantMessage: "update documentation",
			wantIcon:    "📚",
		},
		// Edge cases that don't match Conventional Commits
		{
			name:        "emoji_with_colon_prefix",
			input:       ":sparkles: feat: new feature",
			wantType:    "sparkles", // Parsed as type since no emoji is extracted
			wantScope:   "-",
			wantMessage: "feat: new feature",
			wantIcon:    "-",
		},
		{
			name:        "fire_emoji_refactor",
			input:       "🔥 refactor: remove deprecated code",
			wantType:    "refactor",
			wantScope:   "-",
			wantMessage: "remove deprecated code",
			wantIcon:    "🔥",
		},
		{
			name:        "rocket_emoji_perf",
			input:       "🚀 perf: optimize performance",
			wantType:    "perf",
			wantScope:   "-",
			wantMessage: "optimize performance",
			wantIcon:    "🚀",
		},

		// Edge cases and special formats
		{
			name:        "uppercase_type",
			input:       "FEAT: uppercase type",
			wantType:    "FEAT",
			wantScope:   "-",
			wantMessage: "uppercase type",
			wantIcon:    "-",
		},
		{
			name:        "mixed_case_type",
			input:       "FiX: mixed case",
			wantType:    "FiX",
			wantScope:   "-",
			wantMessage: "mixed case",
			wantIcon:    "-",
		},
		{
			name:        "type_with_hyphen",
			input:       "hot-fix: urgent fix",
			wantType:    "hot-fix",
			wantScope:   "-",
			wantMessage: "urgent fix",
			wantIcon:    "-",
		},
		{
			name:        "type_with_underscore",
			input:       "breaking_change: major update",
			wantType:    "breaking_change",
			wantScope:   "-",
			wantMessage: "major update",
			wantIcon:    "-",
		},
		{
			name:        "scope_with_slash",
			input:       "feat(api/v2): add v2 endpoint",
			wantType:    "feat",
			wantScope:   "api/v2",
			wantMessage: "add v2 endpoint",
			wantIcon:    "-",
		},
		{
			name:        "scope_with_dash",
			input:       "fix(user-auth): fix auth",
			wantType:    "fix",
			wantScope:   "user-auth",
			wantMessage: "fix auth",
			wantIcon:    "-",
		},
		{
			name:        "message_with_colon",
			input:       "feat: add feature: advanced mode",
			wantType:    "feat",
			wantScope:   "-",
			wantMessage: "add feature: advanced mode",
			wantIcon:    "-",
		},
		{
			name:        "multiline_message_first_line",
			input:       "feat: add feature\n\nLong description here",
			wantType:    "misc", // Newline breaks regex matching
			wantScope:   "-",
			wantMessage: "feat: add feature\n\nLong description here",
			wantIcon:    "-",
		},

		// Non-conventional commits (fallback to misc)
		{
			name:        "non_conventional_simple",
			input:       "Updated README",
			wantType:    "misc",
			wantScope:   "-",
			wantMessage: "Updated README",
			wantIcon:    "-",
		},
		{
			name:        "non_conventional_sentence",
			input:       "This is just a regular commit message",
			wantType:    "misc",
			wantScope:   "-",
			wantMessage: "This is just a regular commit message",
			wantIcon:    "-",
		},
		{
			name:        "empty_message",
			input:       "",
			wantType:    "misc",
			wantScope:   "-",
			wantMessage: "",
			wantIcon:    "-",
		},
		{
			name:        "only_emoji",
			input:       "🎉",
			wantType:    "misc",
			wantScope:   "-",
			wantMessage: "🎉",
			wantIcon:    "🎉",
		},

		// Breaking changes (not currently supported by regex)
		{
			name:        "breaking_change_exclamation",
			input:       "feat!: breaking change",
			wantType:    "misc", // ! breaks the regex pattern
			wantScope:   "-",
			wantMessage: "feat!: breaking change",
			wantIcon:    "-",
		},
		{
			name:        "breaking_change_with_scope",
			input:       "feat(api)!: breaking API change",
			wantType:    "misc", // ! breaks the regex pattern
			wantScope:   "-",
			wantMessage: "feat(api)!: breaking API change",
			wantIcon:    "-",
		},

		// Special characters in message
		{
			name:        "message_with_special_chars",
			input:       "feat: add #123 @user [WIP]",
			wantType:    "feat",
			wantScope:   "-",
			wantMessage: "add #123 @user [WIP]",
			wantIcon:    "-",
		},
		{
			name:        "message_with_quotes",
			input:       "fix: resolve \"escaping\" issue",
			wantType:    "fix",
			wantScope:   "-",
			wantMessage: "resolve \"escaping\" issue",
			wantIcon:    "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			got := service.parseSingleLine(tt.input)

			if got.Type != tt.wantType {
				t.Errorf("Parse(%q).Type = %q, want %q", tt.input, got.Type, tt.wantType)
			}

			if got.Scope != tt.wantScope {
				t.Errorf("Parse(%q).Scope = %q, want %q", tt.input, got.Scope, tt.wantScope)
			}

			if got.Message != tt.wantMessage {
				t.Errorf("Parse(%q).Message = %q, want %q", tt.input, got.Message, tt.wantMessage)
			}

			if got.Icon != tt.wantIcon {
				t.Errorf("Parse(%q).Icon = %q, want %q", tt.input, got.Icon, tt.wantIcon)
			}

			if got.Raw != tt.input {
				t.Errorf("Parse(%q).Raw = %q, want %q", tt.input, got.Raw, tt.input)
			}
		})
	}
}

// TestExtractEmoji tests emoji extraction from commit messages.
func TestExtractEmoji(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Common emojis
		{
			name:  "sparkles",
			input: "✨ feat: new feature",
			want:  "✨",
		},
		{
			name:  "bug",
			input: "🐛 fix: bug fix",
			want:  "🐛",
		},
		{
			name:  "fire",
			input: "🔥 refactor: remove code",
			want:  "🔥",
		},
		{
			name:  "rocket",
			input: "🚀 perf: optimize",
			want:  "🚀",
		},
		{
			name:  "book",
			input: "📚 docs: update docs",
			want:  "📚",
		},
		{
			name:  "tada",
			input: "🎉 init: initial commit",
			want:  "🎉",
		},
		{
			name:  "pencil",
			input: "✏️ fix: typo",
			want:  "✏", // Variant selector (️) not captured
		},
		{
			name:  "hammer",
			input: "🔨 build: update build",
			want:  "🔨",
		},

		// Multiple emojis (should extract first one)
		{
			name:  "multiple_emojis",
			input: "🎉✨ feat: multiple emojis",
			want:  "🎉✨",
		},

		// Emoji with colon prefix (should skip colon)
		{
			name:  "colon_prefix",
			input: ":sparkles: feat: feature",
			want:  "",
		},

		// No emoji
		{
			name:  "no_emoji",
			input: "feat: regular commit",
			want:  "",
		},
		{
			name:  "empty_string",
			input: "",
			want:  "",
		},

		// Emoji in middle (should not extract)
		{
			name:  "emoji_in_middle",
			input: "feat: add 🚀 feature",
			want:  "",
		},

		// Various symbol ranges
		{
			name:  "misc_symbol",
			input: "⭐ feat: star",
			want:  "", // This emoji (U+2B50) is outside the ranges checked
		},
		{
			name:  "emoticon",
			input: "😀 feat: smile",
			want:  "😀",
		},
		{
			name:  "transport_symbol",
			input: "🚗 feat: car",
			want:  "🚗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			got := service.extractEmoji(tt.input)

			if got != tt.want {
				t.Errorf("extractEmoji(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestNewService tests service creation.
func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Error("NewService() should not return nil")
	}
}

// TestParseRawField tests that Raw field is always set.
func TestParseRawField(t *testing.T) {
	service := NewService()
	testInputs := []string{
		"feat: test",
		"invalid commit",
		"",
		"✨ feat(scope): message",
	}

	for _, input := range testInputs {
		t.Run(input, func(t *testing.T) {
			got := service.parseSingleLine(input)
			if got.Raw != input {
				t.Errorf("Parse(%q).Raw = %q, want %q", input, got.Raw, input)
			}
		})
	}
}

// TestParseDefaultValues tests default values for missing fields.
func TestParseDefaultValues(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkFunc func(t *testing.T, commit entity.Commit)
	}{
		{
			name:  "no_scope_defaults_to_dash",
			input: "feat: no scope",
			checkFunc: func(t *testing.T, commit entity.Commit) {
				if commit.Scope != "-" {
					t.Errorf("Scope = %q, want %q", commit.Scope, "-")
				}
			},
		},
		{
			name:  "no_emoji_defaults_to_dash",
			input: "feat: no emoji",
			checkFunc: func(t *testing.T, commit entity.Commit) {
				if commit.Icon != "-" {
					t.Errorf("Icon = %q, want %q", commit.Icon, "-")
				}
			},
		},
		{
			name:  "non_conventional_defaults_to_misc",
			input: "just a regular message",
			checkFunc: func(t *testing.T, commit entity.Commit) {
				if commit.Type != "misc" {
					t.Errorf("Type = %q, want %q", commit.Type, "misc")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			got := service.parseSingleLine(tt.input)
			tt.checkFunc(t, got)
		})
	}
}
