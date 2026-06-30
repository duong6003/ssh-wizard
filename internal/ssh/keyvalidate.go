package ssh

import (
	"fmt"
	"os"
	"strings"
)

type ValidatedKey struct {
	Path          string
	KeyType       KeyType
	Fingerprint   string
	HasPassphrase bool
}

func ValidateKeyFile(path string) (*ValidatedKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key file not found: %s", path)
		}
		return nil, fmt.Errorf("cannot read key file: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, "PRIVATE KEY") {
		return nil, fmt.Errorf("file does not appear to be an SSH private key")
	}

	keyType := detectKeyType(content)
	hasPassphrase := strings.Contains(content, "Proc-Type: 4,ENCRYPTED") ||
		(strings.Contains(content, "BEGIN OPENSSH PRIVATE KEY") && strings.Contains(content, "bcrypt"))

	fingerprint, err := GetFingerprint(path)
	if err != nil {
		return nil, fmt.Errorf("could not read fingerprint — key may be corrupted")
	}

	return &ValidatedKey{
		Path:          path,
		KeyType:       keyType,
		Fingerprint:   fingerprint,
		HasPassphrase: hasPassphrase,
	}, nil
}

func detectKeyType(content string) KeyType {
	firstLine := content
	if idx := strings.Index(content, "\n"); idx != -1 {
		firstLine = content[:idx]
	}
	upper := strings.ToUpper(firstLine)
	if strings.Contains(upper, "RSA") {
		return KeyTypeRSA
	}
	return KeyTypeED25519
}
