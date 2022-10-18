package auth

import (
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/spf13/cobra"

	authSetCmd "github.com/CaleMacDonald/soccer-cli/pkg/cmd/auth/set"
	authStatusCmd "github.com/CaleMacDonald/soccer-cli/pkg/cmd/auth/status"
	authTokenCmd "github.com/CaleMacDonald/soccer-cli/pkg/cmd/auth/token"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate soccer-cli with football-data.org",
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(authSetCmd.NewCmdSet(f, nil))
	cmd.AddCommand(authStatusCmd.NewCmdStatus(f, nil))
	cmd.AddCommand(authTokenCmd.NewCmdToken(f, nil))

	return cmd
}
