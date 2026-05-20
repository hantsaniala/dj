package cmd

import (
	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Show all URL patterns",
	Long:  `List all registered URL patterns. Requires django-extensions (show_urls management command).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runManage("show_urls")
	},
}

func init() {
	rootCmd.AddCommand(routesCmd)
}
