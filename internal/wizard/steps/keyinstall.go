package steps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gossh "ssh-wizard/internal/ssh"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/wizard"
)

type installPhase int

const (
	installPhasePassword installPhase = iota
	installPhaseRunning
	installPhaseDone
	installPhaseError
)

type installStepState struct {
	label  string
	status string
}

type installProgressMsg struct{ step gossh.InstallStep }
type installDoneMsg struct{}
type installErrMsg struct{ err *gossh.InstallKeyError }

type KeyInstallModel struct {
	state     *wizard.State
	phase     installPhase
	passInput textinput.Model
	spinner   spinner.Model
	steps     []installStepState
	errorMsg  string
}

func NewKeyInstall(state *wizard.State) tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := textinput.New()
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '·'

	m := KeyInstallModel{
		state:     state,
		spinner:   s,
		passInput: p,
		steps: []installStepState{
			{"Connecting to server", "pending"},
			{"Creating ~/.ssh/", "pending"},
			{"Setting permissions (700)", "pending"},
			{"Appending public key", "pending"},
			{"Verifying installation", "pending"},
		},
	}

	if state.Key != nil && state.Key.Method == wizard.AuthExisting {
		state.KeyInstalled = true
		m.phase = installPhaseDone
	} else if state.Key != nil && state.Key.Method == wizard.AuthPassword {
		m.phase = installPhasePassword
		m.passInput.Focus()
	} else {
		m.phase = installPhaseRunning
	}

	return m
}

func (m KeyInstallModel) Init() tea.Cmd {
	if m.phase == installPhaseDone {
		return func() tea.Msg { return wizard.AdvanceMsg{} }
	}
	if m.phase == installPhaseRunning {
		return tea.Batch(m.spinner.Tick, m.startInstall(""))
	}
	return m.spinner.Tick
}

func (m KeyInstallModel) startInstall(password string) tea.Cmd {
	state := m.state
	return func() tea.Msg {
		opts := gossh.InstallKeyOptions{
			Hostname:       state.Server.Hostname,
			Port:           state.Server.Port,
			Username:       state.Server.Username,
			Password:       password,
			PrivateKeyPath: state.Key.PrivateKeyPath,
			PublicKeyPath:  state.Key.PublicKeyPath,
		}
		err := gossh.InstallPublicKey(opts, func(step gossh.InstallStep) {})
		if err != nil {
			if installErr, ok := err.(*gossh.InstallKeyError); ok {
				return installErrMsg{installErr}
			}
			return installErrMsg{&gossh.InstallKeyError{Message: err.Error()}}
		}
		return installDoneMsg{}
	}
}

func (m KeyInstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		for i, step := range m.steps {
			if step.status == "running" {
				m.steps[i].status = "done"
			}
		}
		for i, step := range m.steps {
			if step.status == "pending" {
				m.steps[i].status = "running"
				break
			}
		}
		return m, nil

	case installDoneMsg:
		for i := range m.steps {
			m.steps[i].status = "done"
		}
		m.phase = installPhaseDone
		m.state.KeyInstalled = true
		return m, func() tea.Msg { return wizard.AdvanceMsg{} }

	case installErrMsg:
		m.phase = installPhaseError
		m.errorMsg = msg.err.Message
		for i, step := range m.steps {
			if step.status == "running" {
				m.steps[i].status = "error"
				break
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.phase == installPhasePassword && msg.Type == tea.KeyEnter {
			password := m.passInput.Value()
			m.phase = installPhaseRunning
			m.steps[0].status = "running"
			return m, tea.Batch(m.spinner.Tick, m.startInstall(password))
		}
		if m.phase == installPhaseError && msg.String() == "r" {
			for i := range m.steps {
				m.steps[i].status = "pending"
			}
			m.phase = installPhaseRunning
			m.errorMsg = ""
			return m, tea.Batch(m.spinner.Tick, m.startInstall(""))
		}
	}

	var spinCmd tea.Cmd
	m.spinner, spinCmd = m.spinner.Update(msg)
	var passCmd tea.Cmd
	m.passInput, passCmd = m.passInput.Update(msg)
	return m, tea.Batch(spinCmd, passCmd)
}

func (m KeyInstallModel) View() string {
	var b strings.Builder
	b.WriteString(ui.RenderHeader("KEY INSTALLATION", 3, 7))

	if m.phase == installPhasePassword {
		b.WriteString(fmt.Sprintf("\n  %s %s\n",
			ui.Label("password for"),
			ui.Interactive.Render(m.state.Server.Username+"@"+m.state.Server.Hostname),
		))
		b.WriteString("  " + ui.Secondary.Render("Used once to install your key. Not stored.") + "\n\n")
		b.WriteString("  " + m.passInput.View() + "\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("\n  %s\n\n",
		ui.Secondary.Render("Installing your public key on "+m.state.Server.Username+"@"+m.state.Server.Hostname+" ···"),
	))

	for _, step := range m.steps {
		switch step.status {
		case "done":
			b.WriteString(fmt.Sprintf("  %s  %s\n", ui.Success.Render(ui.Sym.Success), ui.Primary.Render(step.label)))
		case "running":
			b.WriteString(fmt.Sprintf("  %s  %s\n", m.spinner.View(), ui.Primary.Render(step.label)))
		case "error":
			b.WriteString(fmt.Sprintf("  %s  %s\n", ui.Error.Render(ui.Sym.Error), ui.Error.Render(step.label)))
		default:
			b.WriteString(fmt.Sprintf("  %s  %s\n", ui.Secondary.Render(ui.Sym.Inactive), ui.Secondary.Render(step.label)))
		}
	}

	if m.phase == installPhaseError {
		b.WriteString(ui.RenderWarningBlock(m.errorMsg, "Press R to retry or Ctrl+C to quit"))
	}

	return b.String()
}
