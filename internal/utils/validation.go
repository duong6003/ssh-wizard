package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	aliasRe    = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	hostnameRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?)*$`)
	usernameRe = regexp.MustCompile(`^[a-z_][a-z0-9_-]*$`)
)

func ValidateAlias(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return "Host alias is required"
	}
	if len(v) > 64 {
		return "Must be 64 characters or fewer"
	}
	if !aliasRe.MatchString(v) {
		return "Only letters, numbers, - _ . allowed"
	}
	return true
}

func ValidateHostname(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return "Hostname is required"
	}
	if strings.ContainsAny(v, " \t") {
		return "Not a valid hostname or IP address"
	}
	if ip := net.ParseIP(v); ip != nil {
		return true
	}
	if looksLikeIPv4(v) {
		return "Not a valid hostname or IP address"
	}
	if hostnameRe.MatchString(v) {
		return true
	}
	return "Not a valid hostname or IP address"
}

func looksLikeIPv4(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) != 4 {
		return false
	}
	for _, ch := range parts[0] {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func ValidateUsername(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return "Username is required"
	}
	if len(v) > 0 && v[0] >= '0' && v[0] <= '9' {
		return "Username must start with a letter or underscore"
	}
	if len(v) > 32 {
		return "Must be 32 characters or fewer"
	}
	if !usernameRe.MatchString(v) {
		return "Lowercase letters, numbers, - and _ only"
	}
	return true
}

func ValidatePort(v string) interface{} {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return "Must be a number"
	}
	if n < 1 || n > 65535 {
		return fmt.Sprintf("Must be between 1 and 65535")
	}
	return true
}
