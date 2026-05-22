package cmd

import (
	"os/exec"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [path...]",
	Short: "Lint Python code with ruff or flake8",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.RequireOne("pip install ruff", "ruff", "flake8"); err != nil {
			return err
		}
		targets := args
		if len(targets) == 0 {
			targets = []string{"."}
		}
		if _, err := exec.LookPath("ruff"); err == nil {
			c := runner.Command("ruff", append([]string{"check"}, targets...)...)
			return runner.Run(c)
		}
		c := runner.Command("flake8", targets...)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
