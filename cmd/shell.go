package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shellCmd = &cobra.Command{
	Use:   "shell [service]",
	Short: "Open a bash shell in a Docker container",
	Long:  `Open an interactive bash shell inside a running Docker container. Optionally specify a service name (default: docker.image config).`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.Require("docker", "https://docs.docker.com/engine/install/"); err != nil {
			return err
		}
		image := viper.GetString("docker.image")
		if image == "" {
			return fmt.Errorf("docker.image not configured")
		}
		name := image
		if len(args) > 0 {
			name = args[0]
		}
		c := runner.Command("docker", "exec", "-it", name, "bash")
		return runner.Run(c)
	},
}
