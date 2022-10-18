package root

import (
	"bytes"
	"fmt"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/kr/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"sort"
	"strings"
)

func rootUsageFunc(w io.Writer, command *cobra.Command) error {
	_, _ = fmt.Fprintf(w, "Usage: %s", command.UseLine())

	subcommands := command.Commands()
	if len(subcommands) > 0 {
		_, _ = fmt.Fprintf(w, "\n\nAvailable Commands:\n")
		for _, c := range subcommands {
			if c.Hidden {
				continue
			}

			_, _ = fmt.Fprintf(w, "  %s\n", c.Name())
		}

		return nil
	}

	flagUsages := command.LocalFlags().FlagUsages()
	if flagUsages != "" {
		_, _ = fmt.Fprintln(w, "\n\nFlags:")
		_, _ = fmt.Fprintf(w, text.Indent(dedent(flagUsages), "  "))
	}

	return nil
}

// Display helpful error message in case subcommand name was mistyped.
// This matches Cobra's behavior for root command, which Cobra
// confusingly doesn't apply to nested commands.
func nestedSuggestFunc(w io.Writer, command *cobra.Command, arg string) {
	_, _ = fmt.Fprintf(w, "unknown command %q for %q\n", arg, command.CommandPath())

	var candidates []string
	if arg == "help" {
		candidates = []string{"--help"}
	} else {
		if command.SuggestionsMinimumDistance <= 0 {
			command.SuggestionsMinimumDistance = 2
		}
		candidates = command.SuggestionsFor(arg)
	}

	if len(candidates) > 0 {
		_, _ = fmt.Fprint(w, "\nDid you mean this?\n")
		for _, c := range candidates {
			_, _ = fmt.Fprintf(w, "\t%s\n", c)
		}
	}

	_, _ = fmt.Fprint(w, "\n")
	_ = rootUsageFunc(w, command)
}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootFlagErrorFunc(_ *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}

var hasFailed bool

// HasFailed signals that the main process should exit with non-zero status
func HasFailed() bool {
	return hasFailed
}

func rootHelpFunc(f *cmdutil.Factory, command *cobra.Command, args []string) {
	if isRootCmd(command) {
		if versionVal, err := command.Flags().GetBool("version"); err == nil && versionVal {
			_, _ = fmt.Fprint(f.IOStreams.Out, command.Annotations["versionInfo"])
			return
		} else if err != nil {
			_, _ = fmt.Fprintln(f.IOStreams.ErrOut, err)
			hasFailed = true
			return
		}
	}

	cs := f.IOStreams.ColorScheme()

	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(f.IOStreams.ErrOut, command, args[1])
		hasFailed = true
		return
	}

	namePadding := 12
	var coreCommands []string
	var actionsCommands []string
	var additionalCommands []string
	for _, c := range command.Commands() {
		if c.Short == "" {
			continue
		}
		if c.Hidden {
			continue
		}

		s := rpad(c.Name()+":", namePadding) + c.Short
		if _, ok := c.Annotations["IsCore"]; ok {
			coreCommands = append(coreCommands, s)
		} else if _, ok := c.Annotations["IsActions"]; ok {
			actionsCommands = append(actionsCommands, s)
		} else {
			additionalCommands = append(additionalCommands, s)
		}
	}

	// If there are no core commands, assume everything is a core command
	if len(coreCommands) == 0 {
		coreCommands = additionalCommands
		additionalCommands = []string{}
	}

	type helpEntry struct {
		Title string
		Body  string
	}

	longText := command.Long
	if longText == "" {
		longText = command.Short
	}
	if longText != "" && command.LocalFlags().Lookup("jq") != nil {
		longText = strings.TrimRight(longText, "\n") +
			"\n\nFor more information about output formatting flags, see `soccer-cli help formatting`."
	}

	var helpEntries []helpEntry
	if longText != "" {
		helpEntries = append(helpEntries, helpEntry{"", longText})
	}
	helpEntries = append(helpEntries, helpEntry{"USAGE", command.UseLine()})
	if len(coreCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"CORE COMMANDS", strings.Join(coreCommands, "\n")})
	}
	if len(actionsCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"ACTIONS COMMANDS", strings.Join(actionsCommands, "\n")})
	}
	if len(additionalCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"ADDITIONAL COMMANDS", strings.Join(additionalCommands, "\n")})
	}

	if isRootCmd(command) {
		var helpTopics []string
		if c := findCommand(command, "actions"); c != nil {
			helpTopics = append(helpTopics, rpad(c.Name()+":", namePadding)+c.Short)
		}
		for topic, params := range HelpTopics {
			helpTopics = append(helpTopics, rpad(topic+":", namePadding)+params["short"])
		}
		sort.Strings(helpTopics)
		helpEntries = append(helpEntries, helpEntry{"HELP TOPICS", strings.Join(helpTopics, "\n")})
	}

	flagUsages := command.LocalFlags().FlagUsages()
	if flagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{"FLAGS", dedent(flagUsages)})
	}
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()
	if inheritedFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{"INHERITED FLAGS", dedent(inheritedFlagUsages)})
	}
	if _, ok := command.Annotations["help:arguments"]; ok {
		helpEntries = append(helpEntries, helpEntry{"ARGUMENTS", command.Annotations["help:arguments"]})
	}
	if command.Example != "" {
		helpEntries = append(helpEntries, helpEntry{"EXAMPLES", command.Example})
	}
	if _, ok := command.Annotations["help:environment"]; ok {
		helpEntries = append(helpEntries, helpEntry{"ENVIRONMENT VARIABLES", command.Annotations["help:environment"]})
	}
	helpEntries = append(helpEntries, helpEntry{"LEARN MORE", `
Use 'soccer-cli <command> <subcommand> --help' for more information about a command.`})
	if _, ok := command.Annotations["help:feedback"]; ok {
		helpEntries = append(helpEntries, helpEntry{"FEEDBACK", command.Annotations["help:feedback"]})
	}

	out := f.IOStreams.Out
	for _, e := range helpEntries {
		if e.Title != "" {
			// If there is a title, add indentation to each line in the body
			_, _ = fmt.Fprintln(out, cs.Bold(e.Title))
			_, _ = fmt.Fprintln(out, text.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {
			// If there is no title print the body as is
			_, _ = fmt.Fprintln(out, e.Body)
		}
		_, _ = fmt.Fprintln(out)
	}
}

func findCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, c := range cmd.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		_, _ = fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
