package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
)

type AuthCmd struct {
}

func NewAuthCommand() *AuthCmd {
	return &AuthCmd{}
}

func (c *AuthCmd) CobraCommand() *cobra.Command {
	setCommand := &cobra.Command{
		Use:   "set",
		Short: "Set the api key used to access football-data.org",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := args[0]
			r, err := regexp.Compile("[a-z1-9]+")
			if err != nil {
				return err
			}

			if !r.MatchString(apiKey) {
				return fmt.Errorf("the api key is not valid")
			}

			viper.Set("footballdata.apikey", apiKey)
			return viper.WriteConfig()
		},
	}

	clearCommand := &cobra.Command{
		Use:   "clear",
		Short: "Removes the api key used to access football-data.org",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Set("footballdata.apikey", "")
			return viper.WriteConfig()
		},
	}

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Set the api key used to access football-data.org",
	}
	cmd.AddCommand(setCommand, clearCommand)
	return cmd
}
