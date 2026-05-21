package runner

import (
	"fmt"
	"os/exec"
)

func Require(binary, installHint string) error {
	_, err := exec.LookPath(binary)
	if err != nil {
		return fmt.Errorf("%s not found on PATH\n  Install: %s", binary, installHint)
	}
	return nil
}

func RequireOne(hint string, binaries ...string) error {
	for _, b := range binaries {
		if _, err := exec.LookPath(b); err == nil {
			return nil
		}
	}
	return fmt.Errorf("none of %v found on PATH\n  Install: %s", binaries, hint)
}
