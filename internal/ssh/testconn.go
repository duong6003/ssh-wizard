package ssh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ConnStepStatus string

const (
	StatusPending ConnStepStatus = "pending"
	StatusRunning ConnStepStatus = "running"
	StatusDone    ConnStepStatus = "done"
	StatusError   ConnStepStatus = "error"
)

type ConnStep struct {
	Name       string
	Status     ConnStepStatus
	DurationMs int64
	Detail     string
}

type ConnTestResult struct {
	Success    bool
	Steps      []ConnStep
	ErrorCode  string
	FailedStep string
	Raw        string
}

func TestConnection(alias string, onUpdate func([]ConnStep)) ConnTestResult {
	steps := []ConnStep{
		{Name: "DNS lookup", Status: StatusPending},
		{Name: "TCP connection", Status: StatusPending},
		{Name: "Host verification", Status: StatusPending},
		{Name: "Authentication", Status: StatusPending},
		{Name: "Session established", Status: StatusPending},
	}

	update := func(i int, status ConnStepStatus, detail string) {
		steps[i].Status = status
		steps[i].Detail = detail
		steps[i].DurationMs = time.Now().UnixMilli()
		cp := make([]ConnStep, len(steps))
		copy(cp, steps)
		onUpdate(cp)
	}

	start := time.Now()
	update(0, StatusRunning, "")

	args := []string{
		"-v",
		"-o", "BatchMode=yes",
		"-o", "ConnectTimeout=10",
		"-o", "StrictHostKeyChecking=accept-new",
		alias, "exit 0",
	}

	cmd := exec.Command("ssh", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	raw := stderr.String()

	if strings.Contains(raw, "Connecting to") || strings.Contains(raw, "Trying") {
		update(0, StatusDone, fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
		update(1, StatusRunning, "")
	}
	if strings.Contains(raw, "Connection established") || strings.Contains(raw, "SSH2_MSG_KEXINIT") {
		update(1, StatusDone, "")
		update(2, StatusRunning, "")
	}
	if strings.Contains(raw, "Server host key") || strings.Contains(raw, "Authentications that") {
		update(2, StatusDone, "")
		update(3, StatusRunning, "")
	}
	if strings.Contains(raw, "Authenticated to") {
		update(3, StatusDone, "")
		update(4, StatusRunning, "")
	}

	if err == nil {
		for i := range steps {
			if steps[i].Status != StatusDone {
				update(i, StatusDone, "")
			}
		}
		return ConnTestResult{Success: true, Steps: steps}
	}

	errorCode := classifySSHError(raw, err)
	failedIdx := 0
	for i, step := range steps {
		if step.Status == StatusRunning || step.Status == StatusPending {
			failedIdx = i
			break
		}
	}
	update(failedIdx, StatusError, errorCode)

	return ConnTestResult{
		Success:    false,
		Steps:      steps,
		ErrorCode:  errorCode,
		FailedStep: steps[failedIdx].Name,
		Raw:        raw,
	}
}

func classifySSHError(stderr string, err error) string {
	switch {
	case strings.Contains(stderr, "Name or service not known"),
		strings.Contains(stderr, "Could not resolve hostname"):
		return "DNS_FAILED"
	case strings.Contains(stderr, "Connection refused"):
		return "CONNECTION_REFUSED"
	case strings.Contains(stderr, "Connection timed out"),
		strings.Contains(stderr, "No route to host"):
		return "CONNECTION_TIMEOUT"
	case strings.Contains(stderr, "Permission denied"):
		return "AUTH_FAILED"
	case strings.Contains(stderr, "REMOTE HOST IDENTIFICATION HAS CHANGED"),
		strings.Contains(stderr, "Host key verification failed"):
		return "HOST_KEY_CHANGED"
	default:
		return fmt.Sprintf("SSH_ERROR: %v", err)
	}
}
