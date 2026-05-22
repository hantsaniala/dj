package cmd

import (
	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
)

var startprojectCmd = &cobra.Command{
	Use:   "startproject <name>",
	Short: "Create a new Django project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.Require("uv", "curl -LsSf https://astral.sh/uv/install.sh | sh"); err != nil {
			return err
		}
		c := runner.Command("uv", "run", "--env-file", ".env", "django-admin", "startproject", args[0])
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(startprojectCmd)
}
