package cmd

import (
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [args...]",
	Short: "Run Django database migrations",
	Long:  `Apply database migrations. Pass app names as arguments for targeted migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage(append([]string{"migrate"}, args...)...)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
