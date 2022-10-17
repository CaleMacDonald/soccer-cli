package cmd

import (
	"fmt"
	footballdata "github.com/CaleMacDonald/football-data"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type LeaguesCmd struct {
	api *footballdata.Api
}

func NewLeaguesCommand(api *footballdata.Api) *LeaguesCmd {
	return &LeaguesCmd{api: api}
}

func (c *LeaguesCmd) CobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "leagues",
		Short: "Shows league information",
		RunE:  c.RunE,
		Args:  cobra.NoArgs,
	}
}

func (c *LeaguesCmd) RunE(_ *cobra.Command, _ []string) error {
	competitions, err := c.api.Competitions()
	if err != nil {
		return err
	}

	err = cacheCompetitions(competitions)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
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

func cacheCompetitions(competitions footballdata.CompetitionsResponse) error {
	values := make([]string, competitions.Count)
	for i, competition := range competitions.Competitions {
		values[i] = fmt.Sprintf("%s|%s", competition.Name, competition.Code)
	}

	viper.Set("footballdata.competitions", values)
	return viper.WriteConfig()
}
