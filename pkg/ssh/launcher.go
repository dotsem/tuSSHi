// Package ssh manages command execution for native SSH interactive client sessions.
package ssh

import (
	"os/exec"
)

// NewSSHCommand constructs an exec.Cmd configured to launch the native ssh client
// for the given host alias. Since standard input, output, and error will be bound
// directly to the terminal, the user will experience the native interactive SSH session.
func NewSSHCommand(alias string) *exec.Cmd {
	// #nosec G204 - alias is a user-selected SSH host name/alias from the config file
	cmd := exec.Command("ssh", alias)
	return cmd
}
