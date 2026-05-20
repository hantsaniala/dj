package cmd

import (
	"github.com/spf13/cobra"
)

var mmCmd = &cobra.Command{
	Use:   "mm [args...]",
	Short: "Make Django database migrations",
	Long:  `Create migration files based on model changes. Scope to specific apps by passing them as arguments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage(append([]string{"makemigrations"}, args...)...)
	},
}

func init() {
	rootCmd.AddCommand(mmCmd)
}
