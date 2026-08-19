package runner

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func Run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the child in its own process group. This way a Ctrl+C broadcast by
	// the terminal is not also delivered to the whole tree all at once, and we
	// can safely forward it to the child's group on shutdown.
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	waitCh := make(chan error, 1)
	go func() { waitCh <- cmd.Wait() }()

	sigCh := make(chan os.Signal, 1)
	defer signal.Stop(sigCh)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-sigCh:
		// Forward the signal to the child's whole process group so the full
		// uvicorn tree shuts down, then block until it exits. Waiting prevents
		// orphaned child processes from continuing to write shutdown lines to
		// the terminal after dj has already returned, which is what left the
		// prompt on a dirty line.
		forwardSignal(cmd, sig)
		<-waitCh
		// Leave the cursor on a fresh line before the shell draws its prompt.
		fmt.Fprintln(os.Stderr)
		return nil
	case err := <-waitCh:
		if err != nil {
			// Ensure the prompt does not overlap the last line of output.
			fmt.Fprintln(os.Stderr)
		}
		return err
	}
}
