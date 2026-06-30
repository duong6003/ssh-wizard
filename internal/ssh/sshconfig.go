package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ssh-wizard/internal/utils"
)

type SSHConfigEntry struct {
	Host                string
	Hostname            string
	User                string
	Port                int
	IdentityFile        string
	IdentitiesOnly      bool
	ServerAliveInterval int
	Compression         bool
	ForwardAgent        bool
}

func ParseSSHConfig(content string) []SSHConfigEntry {
	var entries []SSHConfigEntry
	var current *SSHConfigEntry

	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var key, value string
		if idx := strings.IndexAny(line, " ="); idx != -1 {
			key = strings.ToLower(strings.TrimSpace(line[:idx]))
			value = strings.TrimSpace(strings.TrimLeft(line[idx:], " ="))
		} else {
			continue
		}

		if key == "host" {
			if current != nil {
				entries = append(entries, *current)
			}
			current = &SSHConfigEntry{Host: value, Port: 22}
			continue
		}
		if current == nil {
			continue
		}

		switch key {
		case "hostname":
			current.Hostname = value
		case "user":
			current.User = value
		case "port":
			if n, err := strconv.Atoi(value); err == nil {
				current.Port = n
			}
		case "identityfile":
			current.IdentityFile = value
		case "identitiesonly":
			current.IdentitiesOnly = strings.EqualFold(value, "yes")
		case "serveraliveinterval":
			if n, err := strconv.Atoi(value); err == nil {
				current.ServerAliveInterval = n
			}
		case "compression":
			current.Compression = strings.EqualFold(value, "yes")
		case "forwardagent":
			current.ForwardAgent = strings.EqualFold(value, "yes")
		}
	}
	if current != nil {
		entries = append(entries, *current)
	}
	return entries
}

func FindConflict(entries []SSHConfigEntry, alias string) *SSHConfigEntry {
	for i := range entries {
		if strings.EqualFold(entries[i].Host, alias) {
			return &entries[i]
		}
	}
	return nil
}

func FormatEntry(entry SSHConfigEntry) string {
	lines := []string{fmt.Sprintf("Host %s", entry.Host)}
	lines = append(lines, fmt.Sprintf("    HostName %s", entry.Hostname))
	lines = append(lines, fmt.Sprintf("    User %s", entry.User))
	lines = append(lines, fmt.Sprintf("    Port %d", entry.Port))
	if entry.IdentityFile != "" {
		lines = append(lines, fmt.Sprintf("    IdentityFile %s", entry.IdentityFile))
	}
	if entry.IdentitiesOnly {
		lines = append(lines, "    IdentitiesOnly yes")
	}
	if entry.ServerAliveInterval > 0 {
		lines = append(lines, fmt.Sprintf("    ServerAliveInterval %d", entry.ServerAliveInterval))
	}
	if entry.Compression {
		lines = append(lines, "    Compression yes")
	}
	if entry.ForwardAgent {
		lines = append(lines, "    ForwardAgent yes")
	}
	return strings.Join(lines, "\n")
}

func AppendEntry(existing string, entry SSHConfigEntry) string {
	trimmed := strings.TrimRight(existing, "\n ")
	sep := ""
	if trimmed != "" {
		sep = "\n\n"
	}
	return trimmed + sep + FormatEntry(entry) + "\n"
}

func ReadConfig() (string, []SSHConfigEntry, error) {
	path := utils.GetSSHConfigPath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	content := string(data)
	return content, ParseSSHConfig(content), nil
}

func WriteConfig(content string) error {
	path := utils.GetSSHConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0600)
}
