package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and start Docker containers",
	Long:  `Build Docker images and start containers via "docker compose -f <file> up --build -d".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.Require("docker", "https://docs.docker.com/engine/install/"); err != nil {
			return err
		}
		composeFile := viper.GetString("docker.compose_file")
		if composeFile == "" {
			return fmt.Errorf("docker.compose_file not configured")
		}
		c := runner.Command("docker", "compose", "-f", composeFile, "up", "--build", "-d")
		return runner.Run(c)
	},
}
