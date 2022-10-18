package standings

import (
	footballdata "github.com/CaleMacDonald/football-data"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/completion"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewCmdStandings(f *cmdutil.Factory) *cobra.Command {

	return &cobra.Command{
		Use:   "standings",
		Short: "Lists the standings for a particular league",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := f.Config()
			if err != nil {
				return err
			}

			authToken := config.AuthToken()
			api := footballdata.NewFootballData(authToken)

			standings, err := api.Standings(args[0])
			if err != nil {
				return err
			}

			for _, standing := range standings.Standings {
				if standing.Type != "TOTAL" {
					continue
				}

				t := table.NewWriter()
				t.SetOutputMirror(f.IOStreams.Out)
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
		},
		Annotations: map[string]string{
			"IsCore": "true",
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.LeagueCodeCompletionFunc(),
	}
}
