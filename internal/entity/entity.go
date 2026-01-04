// Package entity defines core data structures used throughout the application.
package entity

// Commit represents a parsed git log entry
type Commit struct {
	Raw     string
	Type    string
	Scope   string
	Message string
	Icon    string
}

// Task represents a manual or recurring task
type Task struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Icon    string `json:"icon"`
	Enabled bool   `json:"enabled"`
}
