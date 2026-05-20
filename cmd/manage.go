package cmd

import (
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:   "manage [args...]",
	Short: "Run a Django manage.py command",
	Long:  `Run any Django manage.py command. Escape hatch for commands not explicitly wrapped by dj.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage(args...)
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
