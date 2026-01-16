package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anIcedAntFA/gohome/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, commit hash, and build date of gohome.`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("gohome version %s\n", version.Version)
	if version.Commit != "" {
		fmt.Printf("commit: %s\n", version.Commit)
	}
	if version.Date != "" {
		fmt.Printf("built: %s\n", version.Date)
	}
}
