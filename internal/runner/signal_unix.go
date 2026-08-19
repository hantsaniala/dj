//go:build unix

package runner

import (
	"os"
	"os/exec"
	"syscall"
)

// setProcessGroup places the child in its own process group. The terminal's
// Ctrl+C broadcast then only reaches dj, not the whole child tree at once.
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// forwardSignal sends sig to the child's entire process group so the whole
// subtree (uv -> uvicorn, reloader, etc.) shuts down together.
func forwardSignal(cmd *exec.Cmd, sig os.Signal) {
	if cmd.Process == nil {
		return
	}
	if s, ok := sig.(syscall.Signal); ok {
		// A negative PID targets the process group.
		_ = syscall.Kill(-cmd.Process.Pid, s)
	}
}
