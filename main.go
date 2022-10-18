package main

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/CaleMacDonald/soccer-cli/internal/build"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmd/factory"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmd/root"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	buildDate := build.Date
	buildVersion := build.Version

	cmdFactory := factory.New()

	stderr := cmdFactory.IOStreams.ErrOut

	rootCmd := root.NewCmdRoot(cmdFactory, buildVersion, buildDate)

	cfg, err := cmdFactory.Config()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to read configuration:  %s\n", err)
		return exitError
	}

	authError := errors.New("authError")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// require that the user is authenticated before running most commands
		if cmdutil.IsAuthCheckEnabled(cmd) && !cmdutil.CheckAuth(cfg) {
			_, _ = fmt.Fprintf(stderr, authHelp())
			return authError
		}

		return nil
	}

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		var pagerPipeError *iostreams.ErrClosedPagerPipe
		var noResultsError cmdutil.NoResultsError

		if err == cmdutil.SilentError {
			return exitError
		} else if cmdutil.IsUserCancellation(err) {
			if errors.Is(err, terminal.InterruptErr) {
				// ensure the next shell prompt has its own line
				_, _ = fmt.Fprintf(stderr, "\n")
			}
			return exitCancel
		} else if errors.Is(err, authError) {
			return exitAuth
		} else if errors.As(err, &pagerPipeError) {
			// ignore the error raised when piping to a closed pager
			return exitOK
		} else if errors.As(err, &noResultsError) {
			if cmdFactory.IOStreams.IsStdoutTTY() {
				_, _ = fmt.Fprintln(stderr, noResultsError.Error())
			}
			// no results are not treated as a command failure
			return exitOK
		}

		printError(stderr, err, cmd)
	}

	if root.HasFailed() {
		return exitError
	}

	return exitOK
}

func authHelp() string {
	return heredoc.Doc(`
		To get started with Soccer CLI, please run:  soccer-cli auth set <token>
		Alternatively, populate the SOCCER_CLI_AUTH_TOKEN environment variable with a Football Data API authentication token.
	`)
}

func printError(out io.Writer, err error, cmd *cobra.Command) {
	_, _ = fmt.Fprintln(out, err)

	var flagError *cmdutil.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			_, _ = fmt.Fprintln(out)
		}
		_, _ = fmt.Fprintln(out, cmd.UsageString())
	}
}
