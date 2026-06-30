package steps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gossh "ssh-wizard/internal/ssh"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/wizard"
)

type connPhase int

const (
	connPhaseRunning connPhase = iota
	connPhaseDone
	connPhaseError
)

type connTestDoneMsg struct{ result gossh.ConnTestResult }

type ConnTestModel struct {
	state   *wizard.State
	phase   connPhase
	spinner spinner.Model
	steps   []gossh.ConnStep
	result  *gossh.ConnTestResult
	cursor  int
}

func NewConnTest(state *wizard.State) tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return ConnTestModel{state: state, spinner: s}
}

func (m ConnTestModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.startTest())
}

func (m ConnTestModel) startTest() tea.Cmd {
	alias := m.state.Server.Alias
	return func() tea.Msg {
		result := gossh.TestConnection(alias, func(steps []gossh.ConnStep) {})
		return connTestDoneMsg{result: result}
	}
}

func (m ConnTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connTestDoneMsg:
		m.result = &msg.result
		m.steps = msg.result.Steps
		if msg.result.Success {
			m.phase = connPhaseDone
			m.state.ConnectionSuccess = true
			return m, nil
		}
		m.phase = connPhaseError
		return m, nil

	case tea.KeyMsg:
		switch m.phase {
		case connPhaseDone:
			if msg.Type == tea.KeyEnter {
				return m, func() tea.Msg { return wizard.AdvanceMsg{} }
			}
		case connPhaseError:
			switch msg.String() {
			case "r", "t":
				m.phase = connPhaseRunning
				m.steps = nil
				m.result = nil
				return m, tea.Batch(m.spinner.Tick, m.startTest())
			case "q":
				return m, tea.Quit
			case "enter":
				return m, func() tea.Msg { return wizard.AdvanceMsg{} }
			}
		}
	}

	var spinCmd tea.Cmd
	m.spinner, spinCmd = m.spinner.Update(msg)
	return m, spinCmd
}

func (m ConnTestModel) View() string {
	var b strings.Builder
	b.WriteString(ui.RenderHeader("CONNECTION TEST", 5, 7))
	b.WriteString(fmt.Sprintf("\n  %s\n\n", ui.Secondary.Render("Testing ssh "+m.state.Server.Alias+" "+ui.Sym.Ellipsis)))

	for _, step := range m.steps {
		duration := ""
		if step.DurationMs > 0 {
			duration = fmt.Sprintf("  %s", ui.Secondary.Render(fmt.Sprintf("%.2fs", float64(step.DurationMs)/1000)))
		}
		switch step.Status {
		case gossh.StatusDone:
			b.WriteString(fmt.Sprintf("  %s  %-30s%s\n", ui.Success.Render(ui.Sym.Success), ui.Primary.Render(step.Name), duration))
		case gossh.StatusRunning:
			b.WriteString(fmt.Sprintf("  %s  %s\n", m.spinner.View(), ui.Primary.Render(step.Name+" "+ui.Sym.Ellipsis)))
		case gossh.StatusError:
			b.WriteString(fmt.Sprintf("  %s  %s\n", ui.Error.Render(ui.Sym.Error), ui.Error.Render(step.Name)))
		default:
			b.WriteString(fmt.Sprintf("  %s  %s\n", ui.Secondary.Render(ui.Sym.Inactive), ui.Secondary.Render(step.Name)))
		}
	}

	if m.phase == connPhaseError && m.result != nil {
		b.WriteString("\n")
		b.WriteString(ui.RenderWarningBlock(
			failureTitle(m.result.ErrorCode),
			failureBody(m.result.ErrorCode, m.state.Server.Hostname),
		))
		b.WriteString("\n")
		b.WriteString(ui.RenderNavHint([]string{"R retry", "T test again", "⏎ continue anyway", "Q quit"}))
	} else if m.phase == connPhaseDone {
		b.WriteString(ui.RenderNavHint([]string{"⏎ continue"}))
	}

	return b.String()
}

func failureTitle(code string) string {
	titles := map[string]string{
		"DNS_FAILED":         "Could not resolve hostname",
		"CONNECTION_REFUSED": "Connection refused on port 22",
		"CONNECTION_TIMEOUT": "Connection timed out",
		"AUTH_FAILED":        "Authentication failed",
		"HOST_KEY_CHANGED":   "WARNING: Server fingerprint has changed",
	}
	if title, ok := titles[code]; ok {
		return title
	}
	return "Connection failed (" + code + ")"
}

func failureBody(code, hostname string) string {
	bodies := map[string]string{
		"DNS_FAILED":         "Could not resolve \"" + hostname + "\".\n\n     Check:\n       · Is the hostname correct?\n       · Are you on the right network/VPN?\n       · ping " + hostname,
		"CONNECTION_REFUSED": "Server is reachable but rejected port 22.\n\n     Check:\n       · Is SSH daemon running? (sudo systemctl status ssh)\n       · Is port 22 open in firewall?",
		"CONNECTION_TIMEOUT": "Server did not respond within 10s.\n\n     Check:\n       · Is the server online?\n       · Try: nc -zv " + hostname + " 22",
		"AUTH_FAILED":        "Server rejected your SSH key.\n\n     Try:\n       · Run key installation again: press R\n       · Check auth log: sudo journalctl -u ssh",
		"HOST_KEY_CHANGED":   "Server fingerprint changed — may be rebuilt or MITM.\n\n     To clear old key:\n       ssh-keygen -R " + hostname,
	}
	if body, ok := bodies[code]; ok {
		return body
	}
	return "An unexpected error occurred. Check ssh output manually."
}
