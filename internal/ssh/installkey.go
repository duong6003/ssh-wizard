package ssh

import (
	"fmt"
	"os"
	"strings"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type InstallStep string

const (
	InstallStepConnecting    InstallStep = "connecting"
	InstallStepCreatingDir   InstallStep = "creating-dir"
	InstallStepSettingPerms  InstallStep = "setting-perms"
	InstallStepInstallingKey InstallStep = "installing-key"
	InstallStepVerifying     InstallStep = "verifying"
	InstallStepDone          InstallStep = "done"
)

type InstallKeyOptions struct {
	Hostname       string
	Port           int
	Username       string
	Password       string
	PrivateKeyPath string
	Passphrase     string
	PublicKeyPath  string
}

type InstallKeyError struct {
	Step    InstallStep
	Message string
	Cause   error
}

func (e *InstallKeyError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Step, e.Message, e.Cause)
}

func InstallPublicKey(opts InstallKeyOptions, onStep func(InstallStep)) error {
	onStep(InstallStepConnecting)

	pubKeyData, err := os.ReadFile(opts.PublicKeyPath)
	if err != nil {
		return &InstallKeyError{Step: InstallStepConnecting, Message: "cannot read public key", Cause: err}
	}
	pubKey := strings.TrimSpace(string(pubKeyData))

	var authMethods []gossh.AuthMethod
	if opts.Password != "" {
		authMethods = append(authMethods, gossh.Password(opts.Password))
	}
	if opts.PrivateKeyPath != "" {
		keyData, err := os.ReadFile(opts.PrivateKeyPath)
		if err == nil {
			var signer gossh.Signer
			if opts.Passphrase != "" {
				signer, err = gossh.ParsePrivateKeyWithPassphrase(keyData, []byte(opts.Passphrase))
			} else {
				signer, err = gossh.ParsePrivateKey(keyData)
			}
			if err == nil {
				authMethods = append(authMethods, gossh.PublicKeys(signer))
			}
		}
	}

	cfg := &gossh.ClientConfig{
		User:            opts.Username,
		Auth:            authMethods,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", opts.Hostname, opts.Port)
	client, err := gossh.Dial("tcp", addr, cfg)
	if err != nil {
		return &InstallKeyError{Step: InstallStepConnecting, Message: "connection failed", Cause: err}
	}
	defer client.Close()

	run := func(step InstallStep, command string) error {
		onStep(step)
		sess, err := client.NewSession()
		if err != nil {
			return &InstallKeyError{Step: step, Message: "session failed", Cause: err}
		}
		defer sess.Close()
		if err := sess.Run(command); err != nil {
			return &InstallKeyError{Step: step, Message: "command failed", Cause: err}
		}
		return nil
	}

	if err := run(InstallStepCreatingDir, "mkdir -p ~/.ssh"); err != nil {
		return err
	}
	if err := run(InstallStepSettingPerms, "chmod 700 ~/.ssh"); err != nil {
		return err
	}

	onStep(InstallStepInstallingKey)
	checkSess, _ := client.NewSession()
	var existingOut strings.Builder
	checkSess.Stdout = &existingOut
	checkSess.Run("cat ~/.ssh/authorized_keys 2>/dev/null || true")
	checkSess.Close()

	if !strings.Contains(existingOut.String(), pubKey) {
		escaped := strings.ReplaceAll(pubKey, "'", "'\\''")
		appendCmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", escaped)
		if err := run(InstallStepInstallingKey, appendCmd); err != nil {
			return err
		}
	}

	onStep(InstallStepVerifying)
	verifySess, _ := client.NewSession()
	var verifyOut strings.Builder
	verifySess.Stdout = &verifyOut
	verifySess.Run("cat ~/.ssh/authorized_keys")
	verifySess.Close()
	if !strings.Contains(verifyOut.String(), pubKey) {
		return &InstallKeyError{Step: InstallStepVerifying, Message: "key not found after installation"}
	}

	onStep(InstallStepDone)
	return nil
}
