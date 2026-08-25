package ssh_test

import (
	"os"
	"path/filepath"
	"testing"

	"tusshi/internal/ssh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeygen(t *testing.T) {
	t.Run("generate and read ed25519 keypair", func(t *testing.T) {
		tempDir := t.TempDir()
		keyPath := filepath.Join(tempDir, "id_ed25519_test")

		err := ssh.GenerateKey(keyPath, ssh.KeyTypeED25519, "test-key")
		require.NoError(t, err)

		_, err = os.Stat(keyPath)
		assert.NoError(t, err, "private key file should exist")

		pubKey, err := ssh.ReadPublicKey(keyPath)
		assert.NoError(t, err, "reading public key should succeed")
		assert.Contains(t, pubKey, "ssh-ed25519")
		assert.Contains(t, pubKey, "test-key")
	})

	t.Run("returns error when reading non-existent public key", func(t *testing.T) {
		_, err := ssh.ReadPublicKey("/non/existent/path/id_ed25519")
		assert.Error(t, err)
	})
}
