package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type checkFunc func() (string, string)

func color(s, code string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", code, s)
}

func green(s string) string  { return color(s, "32") }
func red(s string) string    { return color(s, "31") }
func yellow(s string) string { return color(s, "33") }
func cyan(s string) string   { return color(s, "36") }

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose project configuration and health",
	Long:  `Run diagnostics on the project: verify config, environment, dependencies, database, and more.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checks := []struct {
			Name string
			Fn   checkFunc
		}{
			{"Project root", checkProjectRoot},
			{"Config file", checkConfigFile},
			{"Environment file", checkEnvFile},
			{"Critical env vars", checkCriticalEnvVars},
			{"Secret key", checkSecretKey},
			{"Variable expansion", checkVarExpansion},
			{"Python runtime", checkPythonRuntime},
			{"Python dependencies", checkDeps},
			{"Django health", checkDjangoHealth},
			{"Database connection", checkDBConnection},
			{"Migrations", checkMigrations},
			{"ASGI app", checkASGI},
			{"Docker", checkDocker},
			{"Git status", checkGit},
			{"Celery", checkCelery},
		}

		fmt.Println(cyan("dj doctor") + "\n")

		pass, fail, warn := 0, 0, 0
		for _, c := range checks {
			status, msg := c.Fn()
			label := status
			switch status {
			case "PASS":
				label = green("PASS")
				pass++
			case "FAIL":
				label = red("FAIL")
				fail++
			case "WARN":
				label = yellow("WARN")
				warn++
			case "INFO":
				label = cyan("INFO")
			}
			fmt.Printf("[%s]  %s: %s\n", label, c.Name, msg)
		}

		total := pass + fail + warn
		fmt.Printf("\n%d total — %s, %s, %s\n", total, green(fmt.Sprintf("%d passed", pass)), yellow(fmt.Sprintf("%d warnings", warn)), red(fmt.Sprintf("%d failures", fail)))
		if fail > 0 {
			fmt.Println("Fix failures and run 'dj doctor' again.")
		}
		return nil
	},
}

func runCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func runCaptureTimeout(d time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timed out after %v", d)
	}
	return strings.TrimSpace(string(out)), err
}

func runManageCapture(args ...string) (string, error) {
	manager := viper.GetString("manager")
	if manager == "" {
		return "", fmt.Errorf("manager not configured")
	}
	parts := strings.Fields(manager)
	allArgs := append(parts[1:], args...)
	return runCaptureTimeout(30*time.Second, parts[0], allArgs...)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func checkProjectRoot() (string, string) {
	if fileExists("manage.py") {
		return "PASS", "manage.py found"
	}
	return "FAIL", "manage.py not found — are you in a Django project root?"
}

func checkConfigFile() (string, string) {
	if fileExists("dj.toml") {
		return "PASS", "dj.toml found"
	}
	return "WARN", "dj.toml not found — run 'dj init' to create one"
}

func checkEnvFile() (string, string) {
	envFile := viper.GetString("env_file")
	if envFile == "" {
		envFile = ".env"
	}
	info, err := os.Stat(envFile)
	if err != nil {
		return "FAIL", fmt.Sprintf("%s not found", envFile)
	}
	if info.Size() == 0 {
		return "WARN", fmt.Sprintf("%s is empty", envFile)
	}
	return "PASS", fmt.Sprintf("%s found (%d bytes)", envFile, info.Size())
}

func checkCriticalEnvVars() (string, string) {
	critical := []string{"SECRET_KEY", "APP_NAME"}
	missing := make([]string, 0)
	for _, k := range critical {
		if os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	dbs := []string{"DATABASE_URL", "DATABASE_NAME"}
	dbMissing := true
	for _, k := range dbs {
		if os.Getenv(k) != "" {
			dbMissing = false
			break
		}
	}
	if dbMissing {
		missing = append(missing, "DATABASE_URL or DATABASE_NAME")
	}
	if len(missing) == 0 {
		return "PASS", "all critical vars set"
	}
	return "FAIL", fmt.Sprintf("missing: %s", strings.Join(missing, ", "))
}

func checkSecretKey() (string, string) {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return "FAIL", "SECRET_KEY not set"
	}
	placeholder := strings.ToLower(key)
	if placeholder == "changeme" || placeholder == "change-me" || placeholder == "change_me" {
		return "FAIL", "SECRET_KEY is a placeholder — run 'dj secret' to generate one"
	}
	if len(key) < 50 {
		return "WARN", fmt.Sprintf("SECRET_KEY is short (%d chars, recommend ≥50)", len(key))
	}
	return "PASS", fmt.Sprintf("SECRET_KEY set (%d chars)", len(key))
}

func checkVarExpansion() (string, string) {
	envFile := viper.GetString("env_file")
	if envFile == "" {
		envFile = ".env"
	}
	f, err := os.Open(envFile)
	if err != nil {
		return "WARN", "cannot read env file for expansion check"
	}
	defer f.Close()

	re := regexp.MustCompile(`\$\{([^}]+)\}|\$([a-zA-Z_][a-zA-Z0-9_]*)`)
	unresolved := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		matches := re.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			name := m[1]
			if name == "" {
				name = m[2]
			}
			if os.Getenv(name) == "" {
				unresolved[name] = true
			}
		}
	}
	if len(unresolved) > 0 {
		keys := make([]string, 0, len(unresolved))
		for k := range unresolved {
			keys = append(keys, k)
		}
		return "WARN", fmt.Sprintf("unresolved variables: %s", strings.Join(keys, ", "))
	}
	return "PASS", "all variables resolve"
}

func checkPythonRuntime() (string, string) {
	manager := viper.GetString("dep.command")
	if manager == "" {
		manager = "uv"
	}
	runtime := strings.Fields(manager)[0]
	if _, err := exec.LookPath(runtime); err != nil {
		return "FAIL", fmt.Sprintf("%s not found on PATH", runtime)
	}
	return "PASS", fmt.Sprintf("%s found", runtime)
}

func checkDeps() (string, string) {
	manager := viper.GetString("dep.command")
	if manager == "" {
		manager = "uv sync"
	}
	parts := strings.Fields(manager)
	out, err := runCaptureTimeout(30*time.Second, parts[0], append(parts[1:], "--dry-run")...)
	if err != nil {
		return "WARN", "dependencies not installed or lock file missing — run 'dj dep'"
	}
	if strings.Contains(out, "No requirements") || strings.Contains(out, "nothing to") {
		return "PASS", "all dependencies up to date"
	}
	return "WARN", "some dependencies need updating — run 'dj dep'"
}

func checkDjangoHealth() (string, string) {
	if !fileExists("manage.py") {
		return "INFO", "skipped (no manage.py)"
	}
	out, err := runManageCapture("check", "--deploy")
	if err != nil {
		return "WARN", fmt.Sprintf("Django check failed: %s", truncate(out, 80))
	}
	if strings.Contains(out, "no issues") || out == "" {
		return "PASS", "no deployment issues detected"
	}
	return "WARN", fmt.Sprintf("issues found:\n%s", indent(out, "         "))
}

func checkDBConnection() (string, string) {
	if !fileExists("manage.py") {
		return "INFO", "skipped (no manage.py)"
	}
	out, err := runManageCapture("check", "--database", "default")
	if err != nil {
		return "FAIL", fmt.Sprintf("database connection failed: %s", truncate(out, 80))
	}
	return "PASS", "database reachable"
}

func checkMigrations() (string, string) {
	if !fileExists("manage.py") {
		return "INFO", "skipped (no manage.py)"
	}
	out, err := runManageCapture("showmigrations", "--list")
	if err != nil {
		return "WARN", "could not check migrations"
	}
	unapplied := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "[ ]") {
			unapplied++
		}
	}
	if unapplied == 0 {
		return "PASS", "all migrations applied"
	}
	return "WARN", fmt.Sprintf("%d unapplied migrations — run 'dj migrate'", unapplied)
}

func checkASGI() (string, string) {
	command := viper.GetString("serve.command")
	if command == "" {
		return "WARN", "serve.command not configured"
	}
	re := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_.]*):\w+`)
	match := re.FindString(command)
	if match == "" {
		return "INFO", "could not detect ASGI module from serve.command"
	}
	parts := strings.SplitN(match, ":", 2)
	modulePath := strings.ReplaceAll(parts[0], ".", "/") + ".py"
	if fileExists(modulePath) {
		return "PASS", fmt.Sprintf("ASGI module found: %s", modulePath)
	}
	pyPath := modulePath
	if fileExists(pyPath) {
		return "PASS", fmt.Sprintf("ASGI module found: %s", pyPath)
	}
	return "WARN", fmt.Sprintf("ASGI module not found at %s", modulePath)
}

func checkDocker() (string, string) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "INFO", "docker not on PATH"
	}
	out, err := runCaptureTimeout(5*time.Second, "docker", "info")
	if err != nil {
		return "WARN", "Docker daemon not running"
	}
	image := viper.GetString("docker.image")
	if image == "" {
		return "INFO", fmt.Sprintf("Docker available (no docker.image configured)")
	}
	out, err = runCaptureTimeout(5*time.Second, "docker", "ps", "--filter", fmt.Sprintf("name=%s", image), "--format", "{{.Names}}")
	if err == nil && out != "" {
		return "PASS", fmt.Sprintf("container '%s' is running", image)
	}
	return "WARN", fmt.Sprintf("container '%s' not running — run 'dj docker build'", image)
}

func checkGit() (string, string) {
	branch, err := runCapture("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "INFO", "not a git repository"
	}
	dirty, _ := runCapture("git", "status", "--porcelain")
	dirtyCount := 0
	if dirty != "" {
		dirtyCount = len(strings.Split(strings.TrimSpace(dirty), "\n"))
	}
	behind, _ := runCapture("git", "rev-list", "--count", "HEAD..@{upstream}")
	behindStr := "up to date"
	if behind != "" && behind != "0" {
		behindStr = fmt.Sprintf("%s commits behind remote", behind)
	}
	msg := fmt.Sprintf("branch: %s, clean: %t, %s", branch, dirtyCount == 0, behindStr)
	if dirtyCount > 0 {
		return "WARN", msg
	}
	return "INFO", msg
}

func checkCelery() (string, string) {
	app := viper.GetString("celery.app")
	if app == "" {
		return "INFO", "celery.app not configured"
	}
	modulePath := strings.ReplaceAll(app, ".", "/")
	if fileExists(modulePath + ".py") {
		return "INFO", fmt.Sprintf("celery app module found: %s.py", modulePath)
	}
	if fileExists(filepath.Join(modulePath, "__init__.py")) {
		return "INFO", fmt.Sprintf("celery app module found: %s/", modulePath)
	}
	return "WARN", fmt.Sprintf("celery app module not found: %s", app)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
