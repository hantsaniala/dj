package cmd

import (
	"fmt"
	"strings"

	"github.com/hantsaniala/dj/internal/runner"
	"github.com/spf13/viper"
)

func runManage(args ...string) error {
	manager := viper.GetString("manager")
	if manager == "" {
		return fmt.Errorf("manager command not configured")
	}
	parts := strings.Fields(manager)
	allArgs := append(parts[1:], args...)
	c := runner.Command(parts[0], allArgs...)
	return runner.Run(c)
}

func runCommand(cmdKey string) error {
	cmdStr := viper.GetString(cmdKey)
	if cmdStr == "" {
		return fmt.Errorf("%s not configured", cmdKey)
	}
	parts := strings.Fields(cmdStr)
	c := runner.Command(parts[0], parts[1:]...)
	return runner.Run(c)
}
