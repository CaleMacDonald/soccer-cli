package token

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
}

func NewCmdToken(f *cmdutil.Factory, runF func(options *Options) error) *cobra.Command {
	opts := &Options{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print the auth token soccer-cli is configured to use",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return tokenRun(opts)
		},
	}
	return cmd
}

func tokenRun(opts *Options) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	token := cfg.AuthToken()

	if len(token) > 0 {
		_, _ = fmt.Fprintf(opts.IO.Out, "%s\n", token)
	} else {
		_, _ = fmt.Fprintf(opts.IO.Out, "no authentication token\n")
	}

	return nil
}
