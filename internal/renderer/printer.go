package renderer

import (
	"fmt"
	"io"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

type Config struct {
	Format    string // "text" or "table"
	ShowIcon  bool
	ShowScope bool
}

type Printer struct {
	cfg Config
}

func NewPrinter(cfg Config) *Printer {
	return &Printer{cfg: cfg}
}

func (p *Printer) Print(w io.Writer, repoName string, commits []entity.Commit) {
	if len(commits) == 0 {
		return
	}

	if p.cfg.Format == "table" {
		fmt.Println("TODO: update later")
	} else {
		p.printText(w, repoName, commits)
	}
}

func (p *Printer) printText(w io.Writer, repoName string, commits []entity.Commit) {
	fmt.Fprintf(w, "\n📁 Repository: %s\n", repoName)

	for _, c := range commits {
		line := "- "

		if p.cfg.ShowIcon {
			line += c.Icon + " "
		}

		line += c.Type

		if p.cfg.ShowScope && c.Scope != "" {
			line += "(" + c.Scope + ")"
		}

		line += ": " + c.Message

		fmt.Fprintln(w, line)
	}

	fmt.Fprintln(w, "------------------------------------------")
}
