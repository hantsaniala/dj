package cmd

import (
	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [command] [args...]",
	Short: "Run a command in the project environment (uv run --env-file .env ...)",
	Long:  `Run any command inside the project's uv-managed Python environment. Equivalent to "uv run --env-file .env <command>".`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := runner.Command("uv", append([]string{"run", "--env-file", ".env"}, args...)...)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
