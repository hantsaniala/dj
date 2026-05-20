package cmd

import (
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the uvicorn dev server",
	Long:  `Start the uvicorn development server with hot-reload. Configured via serve.command in config.yaml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand("serve.command")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
