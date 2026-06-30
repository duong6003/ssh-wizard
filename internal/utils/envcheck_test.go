package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/duong6003/ssh-wizard/internal/utils"
)

func TestSupportsUnicodeTerminal(t *testing.T) {
	assert.True(t, utils.SupportsUnicodeTerminal(map[string]string{"WT_SESSION": "1"}, utils.PlatformWindows))
	assert.True(t, utils.SupportsUnicodeTerminal(map[string]string{"TERM": "xterm-256color", "LANG": "en_US.UTF-8"}, utils.PlatformLinux))
	assert.False(t, utils.SupportsUnicodeTerminal(map[string]string{"TERM": "dumb"}, utils.PlatformLinux))
	assert.False(t, utils.SupportsUnicodeTerminal(map[string]string{"SSH_WIZARD_ASCII": "1", "WT_SESSION": "1"}, utils.PlatformWindows))
	assert.False(t, utils.SupportsUnicodeTerminal(map[string]string{}, utils.PlatformWindows))
}

func TestDetectTerminalName(t *testing.T) {
	assert.Equal(t, "Windows Terminal", utils.DetectTerminalName(map[string]string{"WT_SESSION": "abc"}))
	assert.Equal(t, "ConEmu/Cmder", utils.DetectTerminalName(map[string]string{"ConEmuANSI": "ON"}))
	assert.Equal(t, "xterm-256color", utils.DetectTerminalName(map[string]string{"TERM": "xterm-256color"}))
}
