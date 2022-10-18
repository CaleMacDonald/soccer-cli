package root

import (
	"fmt"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/MakeNowJust/heredoc"
	"github.com/kr/text"
	"github.com/spf13/cobra"
	"io"
)

var HelpTopics = map[string]map[string]string{
	"environment": {
		"short": "Environment variables that can be used with soccer-cli",
		"long": heredoc.Doc(`
			NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output.
			CLICOLOR: set to "0" to disable printing ANSI colors in output.
			CLICOLOR_FORCE: set to a value other than "0" to keep ANSI colors in output
			even when the output is piped.
		`),
	},
	"exit-codes": {
		"short": "Exit codes used by soccer-cli",
		"long": heredoc.Doc(`
			soccer-cli follows normal conventions regarding exit codes.

			- If a command completes successfully, the exit code will be 0

			- If a command fails for any reason, the exit code will be 1

			- If a command is running but gets cancelled, the exit code will be 2

			- If a command encounters an authentication issue, the exit code will be 4

			NOTE: It is possible that a particular command may have more exit codes, so it is a good
			practice to check documentation for the command if you are relying on exit codes to
			control some behavior.
		`),
	},
}

func NewHelpTopic(ios *iostreams.IOStreams, topic string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     topic,
		Short:   HelpTopics[topic]["short"],
		Long:    HelpTopics[topic]["long"],
		Example: HelpTopics[topic]["example"],
		Hidden:  true,
		Annotations: map[string]string{
			"markdown:generate": "true",
			"markdown:basename": "soccer_cli_help_" + topic,
		},
	}

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		helpTopicHelpFunc(ios.Out, c, args)
	})
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return helpTopicUsageFunc(ios.ErrOut, c)
	})

	return cmd
}

func helpTopicHelpFunc(w io.Writer, command *cobra.Command, args []string) {
	fmt.Fprint(w, command.Long)
	if command.Example != "" {
		fmt.Fprintf(w, "\n\nEXAMPLES\n")
		fmt.Fprint(w, text.Indent(command.Example, "  "))
	}
}

func helpTopicUsageFunc(w io.Writer, command *cobra.Command) error {
	fmt.Fprintf(w, "Usage: soccer-cli help %s", command.Use)
	return nil
}
