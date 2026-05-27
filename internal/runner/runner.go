package runner

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func Run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		time.Sleep(50 * time.Millisecond)
		fmt.Fprintln(os.Stderr)
	}
	return err
}
