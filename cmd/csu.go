package cmd

import (
	"github.com/spf13/cobra"
)

var csuCmd = &cobra.Command{
	Use:   "csu",
	Short: "Create a Django superuser",
	Long:  `Create a Django superuser via the interactive createsuperuser management command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage("createsuperuser")
	},
}

func init() {
	rootCmd.AddCommand(csuCmd)
}
