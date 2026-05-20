package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Stop and remove Docker containers",
	Long:  `Stop and remove Docker containers via "docker compose -f <file> down".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		composeFile := viper.GetString("docker.compose_file")
		if composeFile == "" {
			return fmt.Errorf("docker.compose_file not configured")
		}
		c := runner.Command("docker", "compose", "-f", composeFile, "down")
		return runner.Run(c)
	},
}
