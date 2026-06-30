package ssh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type KeyType string

const (
	KeyTypeED25519 KeyType = "ed25519"
	KeyTypeRSA     KeyType = "rsa"
)

type GeneratedKey struct {
	PrivateKeyPath string
	PublicKeyPath  string
	Fingerprint    string
	KeyType        KeyType
}

func GenerateKey(keyType KeyType, outputPath, passphrase string) (*GeneratedKey, error) {
	args := []string{"-t", string(keyType)}
	if keyType == KeyTypeRSA {
		args = append(args, "-b", "4096")
	}
	args = append(args, "-f", outputPath, "-N", "", "-C", "ssh-wizard")

	cmd := exec.Command("ssh-keygen", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ssh-keygen failed: %w\n%s", err, stderr.String())
	}

	if passphrase != "" {
		// TODO: Avoid command-line passphrase exposure with an expect/pty flow.
		pCmd := exec.Command("ssh-keygen", "-p", "-f", outputPath, "-P", "", "-N", passphrase)
		var pStderr bytes.Buffer
		pCmd.Stderr = &pStderr
		if err := pCmd.Run(); err != nil {
			return nil, fmt.Errorf("ssh-keygen -p failed: %w\n%s", err, pStderr.String())
		}
	}

	fingerprint, err := GetFingerprint(outputPath)
	if err != nil {
		return nil, err
	}

	return &GeneratedKey{
		PrivateKeyPath: outputPath,
		PublicKeyPath:  outputPath + ".pub",
		Fingerprint:    fingerprint,
		KeyType:        keyType,
	}, nil
}

func GetFingerprint(privateKeyPath string) (string, error) {
	cmd := exec.Command("ssh-keygen", "-l", "-f", privateKeyPath)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not read fingerprint: %w", err)
	}
	parts := strings.Fields(string(out))
	if len(parts) < 2 {
		return strings.TrimSpace(string(out)), nil
	}
	return parts[1], nil
}
