package root

import (
	"github.com/CaleMacDonald/soccer-cli/pkg/cmd/auth"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmd/leagues"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmd/standings"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	versionCmd "github.com/CaleMacDonald/soccer-cli/pkg/cmd/version"
)

func NewCmdRoot(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soccer-cli <command> <subcommand> [flags]",
		Short: "Soccer CLI",
		Long:  "Soccer info directly in your terminal",

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ soccer-cli standings PL
		`),
		Annotations: map[string]string{
			"help:feedback": heredoc.Doc(`
				Open an issue at https://github.com/CaleMacDonald/soccer-cli
			`),
			"versionInfo": versionCmd.Format(version, buildDate),
		},
	}

	cmd.Flags().Bool("version", false, "Show soccer-cli version")
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		rootHelpFunc(f, c, args)
	})
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return rootUsageFunc(f.IOStreams.ErrOut, c)
	})
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))
	cmd.AddCommand(standings.NewCmdStandings(f))
	cmd.AddCommand(leagues.NewCmdLeagues(f, nil))
	cmd.AddCommand(auth.NewCmdAuth(f))

	cmd.AddCommand(NewHelpTopic(f.IOStreams, "environment"))
	cmd.AddCommand(NewHelpTopic(f.IOStreams, "exit-codes"))

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}
