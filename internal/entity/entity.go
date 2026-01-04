package entity

// Commit represents a parsed git log entry
type Commit struct {
	Raw     string
	Type    string
	Scope   string
	Message string
	Icon    string
}
