package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull, install deps, migrate, and start server",
	Long:  `Full update workflow: pull latest code, install dependencies, run migrations, and start the dev server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pullCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err := depCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err := migrateCmd.RunE(cmd, args); err != nil {
			return err
		}
		return serveCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
