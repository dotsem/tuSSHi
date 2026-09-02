// Package ssh provides helpers for SSH key generation and service authentication checks.
package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"tusshi/internal/utils"
)

// Key type identifiers supported by ssh-keygen.
const (
	KeyTypeED25519 = "ed25519"
	KeyTypeRSA     = "rsa"
	KeyTypeECDSA   = "ecdsa"
)

// KeyTypeOptions contains the list of supported SSH key types.
var KeyTypeOptions = []string{KeyTypeED25519, KeyTypeRSA, KeyTypeECDSA}

// GenerateKey runs ssh-keygen to create a new keypair at the given path.
// keyType must be one of: ed25519, rsa, ecdsa. An empty passphrase is always used.
func GenerateKey(path, keyType, comment string) error {
	if keyType == "" {
		keyType = KeyTypeED25519
	}

	path = utils.ExpandTilde(path)

	args := []string{"-t", keyType, "-f", path, "-N", "", "-C", comment}
	if keyType == KeyTypeRSA {
		args = append(args, "-b", "4096")
	}

	// #nosec G204 — path and args are controlled by the user via the TUI wizard
	cmd := exec.Command("ssh-keygen", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh-keygen failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// ReadPublicKey reads the public key file corresponding to the given private key path.
func ReadPublicKey(privateKeyPath string) (string, error) {
	pubPath := utils.ExpandTilde(privateKeyPath) + ".pub"
	data, err := os.ReadFile(pubPath) // #nosec G304 — user-provided path from TUI
	if err != nil {
		return "", fmt.Errorf("reading public key %q: %w", pubPath, err)
	}
	return strings.TrimSpace(string(data)), nil
}
