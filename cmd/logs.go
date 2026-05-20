package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Follow Docker container logs",
	Long:  `Tail logs from a running Docker container. Optionally specify a service name (default: docker.image config).`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := viper.GetString("docker.image")
		if image == "" {
			return fmt.Errorf("docker.image not configured")
		}
		name := image
		if len(args) > 0 {
			name = args[0]
		}
		c := runner.Command("docker", "logs", name, "--follow")
		return runner.Run(c)
	},
}
