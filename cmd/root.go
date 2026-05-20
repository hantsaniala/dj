package cmd

import (
	"fmt"
	"os"

	"github.com/hantsaniala/dj/internal/config"
	"github.com/spf13/cobra"
)

var env string

var rootCmd = &cobra.Command{
	Use:   "dj",
	Short: "Django project CLI helper",
	Long: `dj is a CLI helper for Django project management.
It wraps manage.py, docker, and common development workflows into simple commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "Environment profile (e.g. production)")
	cobra.OnInitialize(func() {
		config.Init(env)
	})
}
