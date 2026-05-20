package cmd

import (
	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update dj to the latest version",
	Long:  `Reinstall dj from source using go install. Requires Go to be installed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := runner.Command("go", "install", "github.com/hantsaniala/dj@latest")
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}
