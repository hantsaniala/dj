package cmd

import (
	"fmt"
	"strings"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Install Python dependencies",
	Long:  `Install Python dependencies. Command configured via dep.command in config.yaml (default: uv sync).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		depCmd := viper.GetString("dep.command")
		if depCmd == "" {
			return fmt.Errorf("dep.command not configured")
		}
		parts := strings.Fields(depCmd)
		if err := runner.Require(parts[0], "curl -LsSf https://astral.sh/uv/install.sh | sh"); err != nil {
			return err
		}
		c := runner.Command(parts[0], parts[1:]...)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(depCmd)
}
