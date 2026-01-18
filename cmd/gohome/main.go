// Package main is the entry point for the gohome CLI application.
// It aggregates git commit reports and provides formatting options.
package main

import (
	"github.com/anIcedAntFA/gohome/cmd/gohome/cmd"
)

func main() {
	cmd.Execute()
}
