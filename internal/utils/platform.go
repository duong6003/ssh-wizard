package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformMacOS   Platform = "macos"
	PlatformWindows Platform = "windows"
	PlatformWSL     Platform = "wsl"
)

func DetectPlatform() Platform {
	if isWSL() {
		return PlatformWSL
	}
	switch runtime.GOOS {
	case "darwin":
		return PlatformMacOS
	case "windows":
		return PlatformWindows
	default:
		return PlatformLinux
	}
}

func isWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(data))
	return strings.Contains(lower, "microsoft") || strings.Contains(lower, "wsl")
}

func ExpandTilde(path string) string {
	home, _ := os.UserHomeDir()
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func GetSSHDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh")
}

func GetSSHConfigPath() string {
	return filepath.Join(GetSSHDir(), "config")
}

func GetDefaultKeyPath(alias, keyType string) string {
	return filepath.Join(GetSSHDir(), "id_"+keyType+"_"+alias)
}
