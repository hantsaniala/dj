package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Init(env string) {
	viper.SetDefault("manager", "uv run --env-file .env manage.py")
	viper.SetDefault("serve.command", "uv run --env-file .env uvicorn --host 0.0.0.0 --port 8000 core.asgi:application --reload")
	viper.SetDefault("docker.image", "")
	viper.SetDefault("docker.compose_file", "docker-compose.yml")
	viper.SetDefault("celery.app", "core")
	viper.SetDefault("git.branch", "develop")
	viper.SetDefault("dep.command", "uv sync")
	viper.SetDefault("env_file", ".env")

	viper.SetConfigName("dj")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	if env != "" {
		override := viper.New()
		override.SetConfigName("dj." + env)
		override.SetConfigType("toml")
		override.AddConfigPath(".")
		if err := override.ReadInConfig(); err == nil {
			for _, key := range override.AllKeys() {
				viper.Set(key, override.Get(key))
			}
		}
	}

	envFile := viper.GetString("env_file")
	loadDotEnv(envFile)

	viper.AutomaticEnv()
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	env := make(map[string]string)
	order := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "'\"")
		env[key] = val
		order = append(order, key)
	}

	for round := 0; round < 5; round++ {
		changed := false
		for key, val := range env {
			expanded := os.Expand(val, func(name string) string {
				if v, ok := env[name]; ok {
					return v
				}
				return os.Getenv(name)
			})
			if expanded != val {
				env[key] = expanded
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	for _, key := range order {
		os.Setenv(key, env[key])
	}
}
