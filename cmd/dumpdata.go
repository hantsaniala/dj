package cmd

import (
	"github.com/spf13/cobra"
)

var dumpdataCmd = &cobra.Command{
	Use:   "dumpdata [args...]",
	Short: "Export database data as JSON fixtures",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage(append([]string{"dumpdata"}, args...)...)
	},
}

func init() {
	rootCmd.AddCommand(dumpdataCmd)
}
