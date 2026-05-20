package cmd

import (
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Pull, build Docker, and run migrations inside container",
	Long:  `Deploy workflow: pull latest code, build Docker images, and run migrations inside the container.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pullCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err := buildCmd.RunE(cmd, args); err != nil {
			return err
		}
		return migrateDockerCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
