package steps

import "ssh-wizard/internal/wizard"

func init() {
	wizard.RegisterStepConstructors(wizard.StepConstructors{
		Welcome:       NewWelcome,
		ServerInfo:    NewServerInfo,
		Auth:          NewAuth,
		KeyInstall:    NewKeyInstall,
		SSHConfigStep: NewSSHConfigStep,
		ConnTest:      NewConnTest,
		Done:          NewDone,
	})
}
