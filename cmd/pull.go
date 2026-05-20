package cmd

import (
	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest changes from git",
	Long:  `Pull latest changes from the configured git branch (git.branch, default: develop). Equivalent to "git pull origin <branch>".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := viper.GetString("git.branch")
		c := runner.Command("git", "pull", "origin", branch)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
