package completion

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func LeagueCodeCompletionFunc() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		availableCompetitions := viper.GetStringSlice("footballdata.competitions")

		var comps []string
		if len(args) == 0 {
			for _, competition := range availableCompetitions {
				split := strings.Split(competition, "|")
				comps = append(comps, split[1])
			}
		} else if len(args) == 1 {
			for _, competition := range availableCompetitions {
				split := strings.Split(competition, "|")
				if strings.Contains(split[1], args[1]) {
					comps = append(comps, split[1])
				}
			}
		}

		return comps, cobra.ShellCompDirectiveNoFileComp
	}
}
