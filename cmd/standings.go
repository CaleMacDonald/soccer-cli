package cmd

import (
	"github.com/CaleMacDonald/football-data"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"soccer-cli/pkg"
	"soccer-cli/pkg/completion"
)

type StandingsCmd struct {
	io  pkg.IOStreams
	api *footballdata.Api
}

func NewStandingsCommand(api *footballdata.Api, streams pkg.IOStreams) *StandingsCmd {
	return &StandingsCmd{
		api: api,
		io:  streams,
	}
}

func (c *StandingsCmd) CobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "standings",
		Short:             "Lists the standings for a particular league",
		RunE:              c.RunE,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.LeagueCodeCompletionFunc(),
	}
}

func (c *StandingsCmd) RunE(_ *cobra.Command, args []string) error {
	standings, err := c.api.Standings(args[0])
	if err != nil {
		return err
	}

	for _, standing := range standings.Standings {
		if standing.Type != "TOTAL" {
			continue
		}

		t := table.NewWriter()
		t.SetOutputMirror(c.io.Out)
		t.AppendHeader(table.Row{"#", "Team", "MP", "W", "D", "L", "GF", "GA", "GD", "Pts"})

		for _, record := range standing.Table {
			//t.AppendSeparator()
			t.AppendRow([]interface{}{
				record.Position,
				record.Team.Name,
				record.PlayedGames,
				record.Won,
				record.Draw,
				record.Lost,
				record.GoalsFor,
				record.GoalsAgainst,
				record.GoalDifference,
				record.Points,
			})
		}

		t.Render()
	}

	return nil
}
