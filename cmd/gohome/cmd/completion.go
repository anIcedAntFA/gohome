package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command.
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for gohome.

Completions provide tab-completion for commands, subcommands, and flag values.

To load completions:

Bash:

  $ source <(gohome completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gohome completion bash > /etc/bash_completion.d/gohome
  # macOS:
  $ gohome completion bash > $(brew --prefix)/etc/bash_completion.d/gohome

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ gohome completion zsh > "${fpath[1]}/_gohome"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ gohome completion fish | source

  # To load completions for each session, execute once:
  $ gohome completion fish > ~/.config/fish/completions/gohome.fish

PowerShell:

  PS> gohome completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gohome completion powershell > gohome.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
