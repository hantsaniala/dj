package cmd

import (
	"github.com/spf13/cobra"
)

var startappCmd = &cobra.Command{
	Use:   "startapp <name>",
	Short: "Create a new Django app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage(append([]string{"startapp"}, args...)...)
	},
}

func init() {
	rootCmd.AddCommand(startappCmd)
}
