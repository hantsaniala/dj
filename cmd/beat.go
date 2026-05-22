package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var beatCmd = &cobra.Command{
	Use:   "beat",
	Short: "Start Celery Beat scheduler",
	Long:  `Start Celery Beat with django-celery-beat DatabaseScheduler for periodic tasks. Configured via celery.app in config.yaml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.Require("uv", "curl -LsSf https://astral.sh/uv/install.sh | sh"); err != nil {
			return err
		}
		app := viper.GetString("celery.app")
		if app == "" {
			return fmt.Errorf("celery.app not configured")
		}
		c := runner.Command("uv", "run", "--env-file", ".env",
			"celery", "-A", app, "beat",
			"-l", "debug",
			"--scheduler", "django_celery_beat.schedulers:DatabaseScheduler",
		)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(beatCmd)
}
