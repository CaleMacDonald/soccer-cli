package cmd

import (
	"fmt"
	footballdata "github.com/CaleMacDonald/football-data"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"soccer-cli/pkg"
)

type RootCmd struct {
	api *footballdata.Api
}

func NewRootCommand() *RootCmd {
	initConfig()

	apiKey := viper.GetString("footballdata.apikey")
	api := footballdata.NewFootballData(apiKey)

	return &RootCmd{
		api: api,
	}
}

func (c *RootCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "soccer-cli",
	}

	streams := pkg.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	authCmd := NewAuthCommand()
	cmd.AddCommand(authCmd.CobraCommand())

	completionCommand := NewCompletionCmd()
	cmd.AddCommand(completionCommand.CobraCommand())

	competitionsCommand := NewLeaguesCommand(c.api)
	cmd.AddCommand(competitionsCommand.CobraCommand())

	standingsCmd := NewStandingsCommand(c.api, streams)
	cmd.AddCommand(standingsCmd.CobraCommand())

	return cmd
}

func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".football-data")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Can't write config:", err)
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Cannot read configuration: ", err)
			os.Exit(1)
		}
	}

}
