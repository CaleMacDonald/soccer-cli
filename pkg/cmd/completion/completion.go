package completion

import (
	"fmt"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdCompletion(io *iostreams.IOStreams) *cobra.Command {
	var shellType string

	cmd := &cobra.Command{
		Use:   "completion -s <shell>",
		Short: "Generate shell completion scripts",
		Long: heredoc.Docf(`
			Generate shell completion scripts for Soccer CLI commands.

			When installing Soccer CLI through a package manager, it's possible that
			no additional shell configuration is necessary to gain completion support. For
			Homebrew, see <https://docs.brew.sh/Shell-Completion>
			
			If you need to set up completions manually, follow the instructions below. The exact
			config file locations might vary based on your system. Make sure to restart your
			shell before testing whether completions are working.
			

			### bash
			First, ensure that you install %[1]sbash-completion%[1]s using your package manager.
			After, add this to your %[1]s~/.bash_profile%[1]s:
				eval "$(soccer-cli completion -s bash)"
			
			### zsh
			Generate a %[1]s_soccer-cli%[1]s completion script and put it somewhere in your %[1]s$fpath%[1]s:
				soccer-cli completion -s zsh > /usr/local/share/zsh/site-functions/_soccer-cli
			Ensure that the following is present in your %[1]s~/.zshrc%[1]s:
				autoload -U compinit
				compinit -i
			
			Zsh version 5.7 or later is recommended.
			

			### fish
			Generate a %[1]ssoccer-cli.fish%[1]s completion script:
				soccer-cli completion -s fish > ~/.config/fish/completions/soccer-cli.fish
			
			### PowerShell
			Open your profile script with:
				mkdir -Path (Split-Path -Parent $profile) -ErrorAction SilentlyContinue
				notepad $profile
			
			Add the line and save the file:
				Invoke-Expression -Command $(soccer-cli completion -s powershell | Out-String)
		`, "`"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if shellType == "" {
				if io.IsStdoutTTY() {
					return cmdutil.FlagErrorf("error: the value for `--shell` is required")
				}
				shellType = "bash"
			}

			w := io.Out
			rootCmd := cmd.Parent()

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletionV2(w, true)
			case "zsh":
				return rootCmd.GenZshCompletion(w)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(w)
			case "fish":
				return rootCmd.GenFishCompletion(w, true)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
		DisableFlagsInUseLine: true,
	}

	cmdutil.DisableAuthCheck(cmd)
	cmdutil.StringEnumFlag(cmd, &shellType, "shell", "s", "", []string{"bash", "zsh", "fish", "powershell"}, "Shell type")

	return cmd
}
