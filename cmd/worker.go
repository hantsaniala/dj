package cmd

import (
	"fmt"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workerCmd = &cobra.Command{
	Use:   "worker [queue]",
	Short: "Start a Celery worker",
	Long:  `Start a Celery worker process. Optionally specify a queue name (default: celery). Configured via celery.app in config.yaml.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := viper.GetString("celery.app")
		if app == "" {
			return fmt.Errorf("celery.app not configured")
		}
		queue := "celery"
		if len(args) > 0 {
			queue = args[0]
		}
		c := runner.Command("uv", "run", "--env-file", ".env",
			"celery", "-A", app, "worker",
			"-l", "debug",
			"-Q", queue,
			"-E",
		)
		return runner.Run(c)
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
