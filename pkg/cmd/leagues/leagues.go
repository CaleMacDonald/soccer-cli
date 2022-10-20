package leagues

import (
	"fmt"
	footballdata "github.com/CaleMacDonald/football-data"
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/internal/statistics"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type Options struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)
}

func NewCmdLeagues(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "leagues",
		Short: "Shows a list of available leagues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return leaguesRun(opts)
		},
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	return cmd
}

func leaguesRun(opts *Options) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	statistics.TrackLeaguesRequest()

	authToken := cfg.AuthToken()
	api := footballdata.NewFootballData(authToken)

	competitions, err := api.Competitions()
	if err != nil {
		return err
	}

	err = cacheCompetitions(competitions, cfg)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(opts.IO.Out)
	t.AppendHeader(table.Row{"League", "League Code"})

	for _, competition := range competitions.Competitions {
		t.AppendRow([]interface{}{
			competition.Name,
			competition.Code,
		})
	}

	t.Render()
	return nil
}

func cacheCompetitions(competitions footballdata.CompetitionsResponse, cfg config.Config) error {
	values := make([]string, competitions.Count)
	for i, competition := range competitions.Competitions {
		values[i] = fmt.Sprintf("%s|%s", competition.Name, competition.Code)
	}

	return cfg.SetLeagues(values)
}
