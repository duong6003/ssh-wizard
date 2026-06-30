package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type EnvironmentReport struct {
	Platform           Platform
	Terminal           string
	SupportsUnicode    bool
	SSHAvailable       bool
	SSHKeygenAvailable bool
	Warnings           []string
	RemoteTargetNote   string
}

func CheckEnvironment() EnvironmentReport {
	env := map[string]string{
		"TERM":             os.Getenv("TERM"),
		"WT_SESSION":       os.Getenv("WT_SESSION"),
		"ConEmuANSI":       os.Getenv("ConEmuANSI"),
		"LANG":             os.Getenv("LANG"),
		"LC_ALL":           os.Getenv("LC_ALL"),
		"SSH_WIZARD_ASCII": os.Getenv("SSH_WIZARD_ASCII"),
	}
	platform := DetectPlatform()
	report := EnvironmentReport{
		Platform:           platform,
		Terminal:           DetectTerminalName(env),
		SupportsUnicode:    SupportsUnicodeTerminal(env, platform),
		SSHAvailable:       commandAvailable("ssh"),
		SSHKeygenAvailable: commandAvailable("ssh-keygen"),
		RemoteTargetNote:   "Remote key installation expects a Unix-like SSH server with mkdir, chmod, cat, and echo.",
	}
	if !report.SupportsUnicode {
		report.Warnings = append(report.Warnings, "Unicode not detected; using ASCII-safe symbols.")
	}
	if !report.SSHAvailable {
		report.Warnings = append(report.Warnings, "ssh executable not found in PATH.")
	}
	if !report.SSHKeygenAvailable {
		report.Warnings = append(report.Warnings, "ssh-keygen executable not found in PATH.")
	}
	if platform == PlatformWindows {
		report.Warnings = append(report.Warnings, "Windows requires OpenSSH Client for ssh and ssh-keygen.")
	}
	return report
}

func DetectTerminalName(env map[string]string) string {
	switch {
	case env["WT_SESSION"] != "":
		return "Windows Terminal"
	case strings.EqualFold(env["ConEmuANSI"], "ON"):
		return "ConEmu/Cmder"
	case env["TERM"] != "":
		return env["TERM"]
	default:
		return runtime.GOOS + " console"
	}
}

func SupportsUnicodeTerminal(env map[string]string, platform Platform) bool {
	if env["SSH_WIZARD_ASCII"] == "1" || strings.EqualFold(env["SSH_WIZARD_ASCII"], "true") {
		return false
	}
	if env["WT_SESSION"] != "" || strings.EqualFold(env["ConEmuANSI"], "ON") {
		return true
	}
	term := strings.ToLower(env["TERM"])
	if term == "dumb" {
		return false
	}
	locale := strings.ToLower(env["LC_ALL"] + " " + env["LANG"])
	if strings.Contains(locale, "utf-8") || strings.Contains(locale, "utf8") {
		return true
	}
	if platform == PlatformWindows && term == "" {
		return false
	}
	return term != ""
}

func commandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
