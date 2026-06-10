// Package ssh manages command execution for native SSH interactive client sessions.
package ssh

import (
	"fmt"
	"os/exec"
)

// NewSSHCommand constructs an exec.Cmd configured to launch the native ssh client
// for the given host alias. Since standard input, output, and error will be bound
// directly to the terminal, the user will experience the native interactive SSH session.
func NewSSHCommand(alias string) *exec.Cmd {
	script := fmt.Sprintf(
		"echo -e '\\ntuSSHi connecting to \\033[1;35m%s\\033[0m...\\n><>   ><>   ><>   ><>   ><>   ><>'; exec ssh %s",
		alias,
		alias,
	)
	// #nosec G204 - alias is a user-selected SSH host name/alias from the config file
	return exec.Command("bash", "-c", script)
}
