package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hantsaniala/dj/internal/config"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Generate and set a new Django SECRET_KEY in .env",
	Long:  `Generate a cryptographically random 50-character SECRET_KEY and write it into .env. Replaces existing key if present.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := config.GenerateKey()
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}
		return updateEnvKey("SECRET_KEY", key)
	},
}

func updateEnvKey(key, value string) error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	found := false
	newLine := fmt.Sprintf("%s='%s'", key, value)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+"=") || strings.HasPrefix(trimmed, "# "+key+"=") {
			lines[i] = newLine
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, newLine)
	}
	return os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
