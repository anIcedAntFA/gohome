package viper

import "fmt"

// Validate validates the Config struct after loading from all sources (Viper, flags, file).
// This separates business logic validation from Viper's loading mechanism,
// following the Single Responsibility Principle.
//
// Validation rules:
//   - Time period values cannot be negative
//   - MaxDepth must be between 1 and 10 (prevent excessive recursion)
//   - Format must be either "text" or "table"
//   - Style must be either "normal" or "markdown"
//   - AllBranches and Branch cannot be used together (conflicting options)
//   - Path must not be empty
//
// Returns error with specific context if validation fails.
func (c *Config) Validate() error {
	// Time period validation
	if c.Hours < 0 {
		return fmt.Errorf("hours cannot be negative, got %d", c.Hours)
	}
	if c.Days < 0 {
		return fmt.Errorf("days cannot be negative, got %d", c.Days)
	}
	if c.Weeks < 0 {
		return fmt.Errorf("weeks cannot be negative, got %d", c.Weeks)
	}
	if c.Months < 0 {
		return fmt.Errorf("months cannot be negative, got %d", c.Months)
	}
	if c.Years < 0 {
		return fmt.Errorf("years cannot be negative, got %d", c.Years)
	}

	// MaxDepth validation (prevent excessive recursion and performance issues)
	if c.MaxDepth < 1 {
		return fmt.Errorf("max-depth must be at least 1, got %d", c.MaxDepth)
	}
	if c.MaxDepth > 10 {
		return fmt.Errorf("max-depth too large (%d), maximum is 10 to prevent performance issues", c.MaxDepth)
	}

	// Path validation
	if c.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Format validation (enum check)
	if c.Format != "text" && c.Format != "table" {
		return fmt.Errorf("invalid format %q, must be 'text' or 'table'", c.Format)
	}

	// Style validation (enum check)
	if c.Style != "normal" && c.Style != "markdown" {
		return fmt.Errorf("invalid style %q, must be 'normal' or 'markdown'", c.Style)
	}

	// Branch filtering conflict validation
	if c.AllBranches && c.Branch != "" {
		return fmt.Errorf("cannot use --all-branches and --branch together; choose one branch filtering method")
	}

	return nil
}
