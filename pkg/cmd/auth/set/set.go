package set

import (
	"fmt"
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"regexp"

	"github.com/MakeNowJust/heredoc"
)

type Options struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)

	ApiKey string
}

func NewCmdSet(f *cmdutil.Factory, runF func(opts *Options) error) *cobra.Command {
	opts := &Options{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "set",
		Args:  cobra.ExactArgs(1),
		Short: "Set the access token",
		Long: heredoc.Docf(`
			Set the access token used to access football-data.org

			To obtain a access token, create an account at https://www.football-data.org/client/register
			and the API key will be emailed to you
		`, "`"),
		Example: heredoc.Doc(`
			$ soccer-cli auth set ABC123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opts.ApiKey = args[0]
			}

			if runF != nil {
				return runF(opts)
			}

			return setRun(opts)
		},
	}

	return cmd
}

func setRun(opts *Options) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	r, err := regexp.Compile("[a-z1-9]+")
	if err != nil {
		return err
	}

	if !r.MatchString(opts.ApiKey) {
		return fmt.Errorf("the api key is not valid")
	}

	err = cfg.SetAuthToken(opts.ApiKey)
	if err == nil {
		fmt.Fprintln(opts.IO.Out, "Auth token is now set")
	}

	return err
}
