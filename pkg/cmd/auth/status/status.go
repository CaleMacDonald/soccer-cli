package status

import (
	"fmt"
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type Options struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)

	ShowToken bool
}

func NewCmdStatus(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Args:  cobra.ExactArgs(0),
		Short: "View authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return statusRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ShowToken, "show-token", "t", false, "Display the auth token")

	return cmd
}

func statusRun(opts *Options) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	token := cfg.AuthToken()
	if len(token) == 0 {
		_, _ = fmt.Fprintln(opts.IO.ErrOut, "no authentication token is set")
		return cmdutil.SilentError
	}

	if opts.ShowToken {
		_, _ = fmt.Fprintf(opts.IO.Out, "authentication token is present and set to %s\n", token)
	} else {
		_, _ = fmt.Fprintln(opts.IO.Out, "authentication token is present")
	}

	return nil
}
