package cmd

import (
	"os/exec"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format [path...]",
	Short: "Format Python code with ruff or black",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		targets := args
		if len(targets) == 0 {
			targets = []string{"."}
		}
		if _, err := exec.LookPath("ruff"); err == nil {
			c := runner.Command("ruff", append([]string{"format"}, targets...)...)
			return runner.Run(c)
		}
		c := runner.Command("black", targets...)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
