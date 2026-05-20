package cmd

import "github.com/spf13/cobra"

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker containers and services",
	Long:  `Manage Docker containers and services for the project. Subcommands: build, clean, shell, logs, migrate.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(
		buildCmd,
		cleanCmd,
		shellCmd,
		logsCmd,
		migrateDockerCmd,
	)
}
