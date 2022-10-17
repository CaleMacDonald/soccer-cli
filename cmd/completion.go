package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	completionLong = `
		Output shell completion code for the specified shell (bash, zsh, fish, or powershell).
		The shell code must be evaluated to provide interactive
		completion of soccer-cli commands.  This can be done by sourcing it from
		the .bash_profile.
		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2.`

	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"fish":       runCompletionFish,
		"powershell": runCompletionPowershell,
	}
)

type CompletionCmd struct {
}

func NewCompletionCmd() *CompletionCmd {
	return &CompletionCmd{}
}

func (c *CompletionCmd) CobraCommand() *cobra.Command {
	var shells []string
	for s := range completionShells {
		shells = append(shells, s)
	}

	return &cobra.Command{
		Use:                   "completion SHELL",
		DisableFlagsInUseLine: true,
		Short:                 "Output shell completion code for the specified shell (bash, zsh, fish, or powershell)",
		Long:                  completionLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCompletion(os.Stdout, cmd, args)
		},
		ValidArgs: shells,
	}
}

// RunCompletion checks given arguments and executes command
func RunCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("shell not specified")
	}

	if len(args) > 1 {
		return fmt.Errorf("too many arguments - expected only the shell type")
	}

	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("unsupported shell type %q", args[0])
	}

	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, command *cobra.Command) error {
	return command.GenBashCompletionV2(out, false) // TODO: Upgrade to Cobra 1.3.0 or later before including descriptions (See https://github.com/spf13/cobra/pull/1509)
}

func runCompletionZsh(out io.Writer, command *cobra.Command) error {
	zshHead := fmt.Sprintf("#compdef %[1]s\ncompdef _%[1]s %[1]s\n", command.Name())
	_, err := out.Write([]byte(zshHead))
	if err != nil {
		return err
	}

	return command.GenZshCompletion(out)
}

func runCompletionFish(out io.Writer, command *cobra.Command) error {
	return command.GenFishCompletion(out, true)
}

func runCompletionPowershell(out io.Writer, command *cobra.Command) error {
	return command.GenPowerShellCompletionWithDesc(out)
}
