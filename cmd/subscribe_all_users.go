package cmd

import (
	"github.com/spf13/cobra"
)

var subscribeAllUsersCmd = &cobra.Command{
	Use:   "subscribe-all-users",
	Short: "Subscribe all existing users to default notification topics",
	Long:  `Subscribe all existing users to default notification topics via the subscribe_all_users management command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage("subscribe_all_users")
	},
}

func init() {
	rootCmd.AddCommand(subscribeAllUsersCmd)
}
