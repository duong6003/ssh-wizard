package ssh_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	ssh "ssh-wizard/internal/ssh"
)

const sampleConfig = `
Host prod
  HostName 192.168.1.1
  User ubuntu
  Port 22
  IdentityFile ~/.ssh/id_ed25519_prod

Host dev-box
  HostName 10.0.0.2
  User ubuntu
`

func TestParseSSHConfig(t *testing.T) {
	entries := ssh.ParseSSHConfig(sampleConfig)
	assert.Len(t, entries, 2)
	assert.Equal(t, "prod", entries[0].Host)
	assert.Equal(t, "192.168.1.1", entries[0].Hostname)
	assert.Equal(t, "ubuntu", entries[0].User)
	assert.Equal(t, 22, entries[0].Port)
	assert.Equal(t, "~/.ssh/id_ed25519_prod", entries[0].IdentityFile)
}

func TestFindConflict(t *testing.T) {
	entries := ssh.ParseSSHConfig(sampleConfig)
	conflict := ssh.FindConflict(entries, "prod")
	assert.NotNil(t, conflict)
	assert.Nil(t, ssh.FindConflict(entries, "staging"))
}

func TestFormatEntry(t *testing.T) {
	entry := ssh.SSHConfigEntry{
		Host:                "my-server",
		Hostname:            "1.2.3.4",
		User:                "ubuntu",
		Port:                22,
		IdentityFile:        "~/.ssh/key",
		IdentitiesOnly:      true,
		ServerAliveInterval: 60,
	}
	out := ssh.FormatEntry(entry)
	assert.Contains(t, out, "Host my-server")
	assert.Contains(t, out, "HostName 1.2.3.4")
	assert.Contains(t, out, "IdentitiesOnly yes")
	assert.Contains(t, out, "ServerAliveInterval 60")
}

func TestAppendEntry(t *testing.T) {
	existing := "Host old\n  HostName old.com\n  User root\n"
	entry := ssh.SSHConfigEntry{Host: "new", Hostname: "new.com", User: "ubuntu", Port: 22}
	result := ssh.AppendEntry(existing, entry)
	assert.Contains(t, result, "Host old")
	assert.Contains(t, result, "Host new")
	oldIdx := strings.Index(result, "Host old")
	newIdx := strings.Index(result, "Host new")
	assert.Greater(t, newIdx, oldIdx)
}
