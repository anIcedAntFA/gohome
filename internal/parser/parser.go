package parser

import (
	"regexp"
	"strings"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// Regex to parse Conventional Commits
var commitRegex = regexp.MustCompile(`(?i)^.*?([a-zA-Z0-9_-]+)(?:\(([^)]+)\))?:\s*(.+)$`)

// Service handles parsing logic
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Parse converts a raw log line into a Commit entity
func (s *Service) Parse(rawLine string) entity.Commit {
	emoji := s.extractEmoji(rawLine)
	if emoji == "" {
		emoji = "-"
	}

	matches := commitRegex.FindStringSubmatch(rawLine)
	commit := entity.Commit{
		Raw:  rawLine,
		Icon: emoji,
	}

	if len(matches) == 4 {
		commit.Type = matches[1]
		commit.Scope = matches[2]
		commit.Message = matches[3]
	} else {
		commit.Type = "misc"
		commit.Scope = "-"
		commit.Message = rawLine
	}

	if commit.Scope == "" {
		commit.Scope = "-"
	}

	return commit
}

func (s *Service) extractEmoji(input string) string {
	var emoji strings.Builder

	for _, r := range input {
		if (r >= 0x1F300 && r <= 0x1F9FF) || // Misc Symbols, Emoticons, Transport
			(r >= 0x2600 && r <= 0x27BF) || // Misc symbols, Dingbats
			(r >= 0x1F000 && r <= 0x1F2FF) { // Additional symbols
			emoji.WriteRune(r)
		} else if r == ' ' || r == ':' {
			continue
		} else {
			break
		}
	}
	return emoji.String()
}
