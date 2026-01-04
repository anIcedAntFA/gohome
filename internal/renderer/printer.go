package renderer

import (
	"fmt"

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

func (p *Printer) Print(repoName string, commits []entity.Commit) {
	if len(commits) == 0 {
		return
	}

	p.printFilterInfo()

	if p.cfg.Format == "table" {
		fmt.Println("TODO: update later")
	} else {
		p.printText(repoName, commits)
	}
}

func (p *Printer) printFilterInfo() {
	// Print filter options
	fmt.Printf("\n🔍 Filters: format=%s", p.cfg.Format)

	if p.cfg.ShowIcon {
		fmt.Printf(", show_icon=true")
	}

	if p.cfg.ShowScope {
		fmt.Printf(", show_scope=true")
	}

	fmt.Println("\n------------------------------------------")
}

func (p *Printer) printText(repoName string, commits []entity.Commit) {
	fmt.Printf("\n📁 Repository: %s\n", repoName)
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

		fmt.Println(line)
	}
	fmt.Println("------------------------------------------")
}
