package steps

import (
	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/wizard"
)

type stub struct{ label string }

func (s stub) Init() tea.Cmd                       { return nil }
func (s stub) Update(tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s stub) View() string                        { return "  " + s.label + "\n" }

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

func NewWelcome() tea.Model                    { return stub{"Welcome"} }
func NewServerInfo(*wizard.State) tea.Model    { return stub{"Server Info"} }
func NewAuth(*wizard.State) tea.Model          { return stub{"Auth"} }
func NewKeyInstall(*wizard.State) tea.Model    { return stub{"Key Install"} }
func NewSSHConfigStep(*wizard.State) tea.Model { return stub{"SSH Config"} }
func NewConnTest(*wizard.State) tea.Model      { return stub{"Connection Test"} }
func NewDone(*wizard.State) tea.Model          { return stub{"Done"} }
