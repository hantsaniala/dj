package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hantsaniala/dj/internal/config"
	"github.com/spf13/cobra"
)

var (
	env     string
	Version = "0.0.0-dev" // overridden by -ldflags at build
)

const asciiArt = `░█▀▄░▀▀█
░█░█░░░█
░▀▀░░▀▀░`

func version() string {
	if Version != "0.0.0-dev" {
		return Version
	}
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return Version
}

var rootCmd = &cobra.Command{
	Use:   "dj",
	Short: "Django project CLI helper",
	Long: `dj is a CLI helper for Django project management.
It wraps manage.py, docker, and common development workflows into simple commands.`,
	Args:              cobra.ArbitraryArgs,
	SilenceErrors:     true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			ext := "dj-" + args[0]
			if _, err := exec.LookPath(ext); err == nil {
				c := exec.Command(ext, args[1:]...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin
				return c.Run()
			}
			return fmt.Errorf("unknown command: %s\nSee 'dj --help'", args[0])
		}
		fmt.Printf("%s %s\n\n", asciiArt, version())
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
	rootCmd.Version = version()
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "Environment profile (e.g. production)")
	rootCmd.SetVersionTemplate(asciiArt + " {{.Version}}\n")
	cobra.OnInitialize(func() {
		config.Init(env)
	})
}
