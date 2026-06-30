package utils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ssh-wizard/internal/utils"
)

func TestValidateAlias(t *testing.T) {
	assert.Equal(t, true, utils.ValidateAlias("prod"))
	assert.Equal(t, true, utils.ValidateAlias("dev-box"))
	assert.Equal(t, true, utils.ValidateAlias("my.server"))
	assert.Contains(t, utils.ValidateAlias(""), "required")
	assert.Contains(t, utils.ValidateAlias("my server"), "letters")
	assert.Contains(t, utils.ValidateAlias(strings.Repeat("a", 65)), "64")
}

func TestValidateHostname(t *testing.T) {
	assert.Equal(t, true, utils.ValidateHostname("192.168.1.1"))
	assert.Equal(t, true, utils.ValidateHostname("server.example.com"))
	assert.Contains(t, utils.ValidateHostname(""), "required")
	assert.Contains(t, utils.ValidateHostname("192.168.x.x"), "valid")
}

func TestValidateUsername(t *testing.T) {
	assert.Equal(t, true, utils.ValidateUsername("ubuntu"))
	assert.Equal(t, true, utils.ValidateUsername("_admin"))
	assert.Contains(t, utils.ValidateUsername(""), "required")
	assert.Contains(t, utils.ValidateUsername("1user"), "letter")
}

func TestValidatePort(t *testing.T) {
	assert.Equal(t, true, utils.ValidatePort("22"))
	assert.Equal(t, true, utils.ValidatePort("2222"))
	assert.Contains(t, utils.ValidatePort("abc"), "number")
	assert.Contains(t, utils.ValidatePort("0"), "1")
	assert.Contains(t, utils.ValidatePort("65536"), "65535")
}
