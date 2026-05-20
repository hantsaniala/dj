package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hantsaniala/dj/internal/config"
	"github.com/spf13/cobra"
)

type initAnswers struct {
	Manager           string
	ServeCommand      string
	DockerImage       string
	DockerComposeFile string
	CeleryApp         string
	GitBranch         string
	DepCommand        string
	EnvFile           string
}

const tomlTemplate = `manager = "{{.Manager}}"
env_file = "{{.EnvFile}}"

[serve]
command = "{{.ServeCommand}}"

[docker]
image = "{{.DockerImage}}"
compose_file = "{{.DockerComposeFile}}"

[celery]
app = "{{.CeleryApp}}"

[git]
branch = "{{.GitBranch}}"

[dep]
command = "{{.DepCommand}}"
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dj configuration interactively",
	Long:  `Create or overwrite dj.toml with interactive prompts, copy .env.example if needed, and generate SECRET_KEY.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("dj.toml"); err == nil {
			overwrite := false
			prompt := &survey.Confirm{
				Message: "dj.toml already exists. Overwrite?",
				Default: false,
			}
			survey.AskOne(prompt, &overwrite)
			if !overwrite {
				fmt.Println("Aborted.")
				return nil
			}
		}

		answers := initAnswers{
			Manager:           "uv run --env-file .env manage.py",
			ServeCommand:      "uv run --env-file .env uvicorn --host 0.0.0.0 --port 8000 core.asgi:application --reload",
			DockerImage:       "",
			DockerComposeFile: "docker-compose.yml",
			CeleryApp:         "core",
			GitBranch:         "develop",
			DepCommand:        "uv sync",
			EnvFile:           ".env",
		}

		qs := []*survey.Question{
			{
				Name: "Manager",
				Prompt: &survey.Input{
					Message: "Manager command:",
					Default: answers.Manager,
				},
			},
			{
				Name: "ServeCommand",
				Prompt: &survey.Input{
					Message: "Dev server command:",
					Default: answers.ServeCommand,
				},
			},
			{
				Name: "DockerImage",
				Prompt: &survey.Input{
					Message: "Docker image name (leave empty to skip docker):",
					Default: answers.DockerImage,
				},
			},
			{
				Name: "DockerComposeFile",
				Prompt: &survey.Input{
					Message: "Docker compose file:",
					Default: answers.DockerComposeFile,
				},
			},
			{
				Name: "CeleryApp",
				Prompt: &survey.Input{
					Message: "Celery app module:",
					Default: answers.CeleryApp,
				},
			},
			{
				Name: "GitBranch",
				Prompt: &survey.Input{
					Message: "Git branch for pull:",
					Default: answers.GitBranch,
				},
			},
			{
				Name: "DepCommand",
				Prompt: &survey.Input{
					Message: "Dependency install command:",
					Default: answers.DepCommand,
				},
			},
			{
				Name: "EnvFile",
				Prompt: &survey.Input{
					Message: "Default env file:",
					Default: answers.EnvFile,
				},
			},
		}

		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		f, err := os.Create("dj.toml")
		if err != nil {
			return fmt.Errorf("failed to create dj.toml: %w", err)
		}
		defer f.Close()

		tmpl := template.Must(template.New("dj").Parse(tomlTemplate))
		if err := tmpl.Execute(f, answers); err != nil {
			return err
		}

		fmt.Println("Created dj.toml")

		if _, err := os.Stat(answers.EnvFile); os.IsNotExist(err) {
			if _, err := os.Stat(".env.example"); err == nil {
				copyEnv := true
				survey.AskOne(&survey.Confirm{
					Message: fmt.Sprintf("%s not found. Copy from .env.example?", answers.EnvFile),
					Default: true,
				}, &copyEnv)
				if copyEnv {
					data, err := os.ReadFile(".env.example")
					if err != nil {
						return err
					}
					if err := os.WriteFile(answers.EnvFile, data, 0644); err != nil {
						return err
					}
					fmt.Printf("Copied .env.example to %s\n", answers.EnvFile)
				}
			} else {
				fmt.Println("Warning: .env.example not found. Create .env manually.")
			}
		}

		if data, err := os.ReadFile(answers.EnvFile); err == nil {
			content := string(data)
			if !stringsContainsAny(content, "SECRET_KEY=", "DJANGO_SECRET_KEY=") {
				genKey := true
				survey.AskOne(&survey.Confirm{
					Message: "SECRET_KEY not found in env file. Generate one?",
					Default: true,
				}, &genKey)
				if genKey {
					key, err := config.GenerateKey()
					if err != nil {
						return fmt.Errorf("failed to generate key: %w", err)
					}
					f, err := os.OpenFile(answers.EnvFile, os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer f.Close()
					fmt.Fprintf(f, "\nSECRET_KEY='%s'\n", key)
					fmt.Println("SECRET_KEY generated and appended to", answers.EnvFile)
				}
			}
		}

		fmt.Println("Done. Run 'dj --help' to see available commands.")
		return nil
	},
}

func stringsContainsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(initCmd)
}
