package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateDockerCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run Django migrations inside Docker container",
	Long:  `Run Django database migrations inside a running Docker container via "docker exec <image> python manage.py migrate".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.Require("docker", "https://docs.docker.com/engine/install/"); err != nil {
			return err
		}
		image := viper.GetString("docker.image")
		if image == "" {
			return fmt.Errorf("docker.image not configured")
		}
		c := runner.Command("docker", "exec", "-it", image, "python", "manage.py", "migrate")
		return runner.Run(c)
	},
}
