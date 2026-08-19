//go:build !unix

package runner

import (
	"os"
	"os/exec"
)

// setProcessGroup is a no-op on platforms without Unix process groups.
func setProcessGroup(cmd *exec.Cmd) {}

// forwardSignal is a no-op on platforms without Unix process groups; the child
// shares the parent process group and receives terminal signals directly.
func forwardSignal(cmd *exec.Cmd, sig os.Signal) {}
