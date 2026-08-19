//go:build unix

package runner

import (
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// TestRunSigintForwardsAndReturns verifies that when the dj process receives
// SIGINT (as the terminal broadcasts on Ctrl+C), Run forwards it to the
// child's process group, waits for the child to shut down, and then returns
// without hanging.
func TestRunSigintForwardsAndReturns(t *testing.T) {
	cmd := exec.Command("bash", "-c",
		`trap 'exit 0' INT; while :; do sleep 0.1; done`)

	done := make(chan error, 1)
	go func() { done <- Run(cmd) }()

	// Give the child time to start, then deliver SIGINT to ourselves exactly
	// as a Ctrl+C foreground-group broadcast would (the child is now in its
	// own process group).
	time.Sleep(400 * time.Millisecond)
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		t.Fatalf("send SIGINT to self: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after SIGINT (child may be orphaned/hanging)")
	}
}
