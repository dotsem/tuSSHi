package ssh

import (
	"errors"
	"os/exec"
	"strings"
)

// AuthResult holds the outcome of an SSH service authentication check.
type AuthResult struct {
	OK    bool
	Error string // first line of stderr, only populated on exit code 255
}

// CheckAuth tests whether the SSH key for the given host alias authenticates
// successfully. It uses BatchMode to suppress interactive prompts and interprets
// exit code 255 as an SSH-level failure. Exit codes 0 and 1 both indicate the
// handshake succeeded — the server simply closed the session without opening a shell,
// which is the expected behaviour for Git hosting services.
func CheckAuth(alias string) AuthResult {
	// #nosec G204 — alias is a Host block name from the user's own ssh config
	cmd := exec.Command("ssh",
		"-T",
		"-o", "BatchMode=yes",
		"-o", "ConnectTimeout=5",
		"-o", "StrictHostKeyChecking=accept-new",
		alias,
	)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return AuthResult{OK: true}
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() != 255 {
		return AuthResult{OK: true}
	}

	errMsg := firstLine(strings.TrimSpace(stderr.String()))
	if errMsg == "" {
		errMsg = err.Error()
	}
	return AuthResult{OK: false, Error: errMsg}
}

func firstLine(s string) string {
	if before, _, ok := strings.Cut(s, "\n"); ok {
		return before
	}
	return s
}
