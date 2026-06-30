package ssh

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectKeyTypeUsesHeaderOnly(t *testing.T) {
	assert.Equal(t, KeyTypeRSA, detectKeyType("-----BEGIN RSA PRIVATE KEY-----\nabc"))
	assert.Equal(t, KeyTypeED25519, detectKeyType("-----BEGIN OPENSSH PRIVATE KEY-----\n"+strings.Repeat("rsa", 10)))
	assert.Equal(t, KeyTypeED25519, detectKeyType("-----BEGIN EC PRIVATE KEY-----\nabc"))
}
