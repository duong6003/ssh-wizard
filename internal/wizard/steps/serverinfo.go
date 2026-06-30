package steps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/duong6003/ssh-wizard/internal/ui"
	"github.com/duong6003/ssh-wizard/internal/utils"
	"github.com/duong6003/ssh-wizard/internal/wizard"
)

type serverField int

const (
	fieldAlias serverField = iota
	fieldHostname
	fieldUsername
	fieldPort
	fieldCount
)

type ServerInfoModel struct {
	state   *wizard.State
	inputs  [fieldCount]textinput.Model
	focused serverField
	errors  [fieldCount]string
}

func NewServerInfo(state *wizard.State) tea.Model {
	m := ServerInfoModel{state: state}
	placeholders := [fieldCount]string{"prod", "192.168.1.1", "ubuntu", "22"}
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = placeholders[i]
		t.CharLimit = 64
		if serverField(i) == fieldPort {
			t.SetValue("22")
		}
		m.inputs[i] = t
	}
	m.inputs[fieldAlias].Focus()
	return m
}

func (m ServerInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ServerInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m = m.validate(m.focused)
			if m.errors[m.focused] == "" {
				m.focused = (m.focused + 1) % fieldCount
				m.refocus()
			}
		case tea.KeyShiftTab, tea.KeyUp:
			if m.focused > 0 {
				m.focused--
				m.refocus()
			}
		case tea.KeyEnter:
			m = m.validate(m.focused)
			if m.errors[m.focused] != "" {
				break
			}
			if m.focused < fieldCount-1 {
				m.focused++
				m.refocus()
				break
			}
			allValid := true
			for i := serverField(0); i < fieldCount; i++ {
				m = m.validate(i)
				if m.errors[i] != "" {
					allValid = false
					m.focused = i
					m.refocus()
					break
				}
			}
			if allValid {
				port, _ := strconv.Atoi(m.inputs[fieldPort].Value())
				m.state.Server = &wizard.ServerConfig{
					Alias:    m.inputs[fieldAlias].Value(),
					Hostname: m.inputs[fieldHostname].Value(),
					Username: m.inputs[fieldUsername].Value(),
					Port:     port,
				}
				return m, func() tea.Msg { return wizard.AdvanceMsg{} }
			}
		}
	}

	var cmds [fieldCount]tea.Cmd
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds[:]...)
}

func (m ServerInfoModel) validate(f serverField) ServerInfoModel {
	value := m.inputs[f].Value()
	var result interface{}
	switch f {
	case fieldAlias:
		result = utils.ValidateAlias(value)
	case fieldHostname:
		result = utils.ValidateHostname(value)
	case fieldUsername:
		result = utils.ValidateUsername(value)
	case fieldPort:
		result = utils.ValidatePort(value)
	}
	if result == true {
		m.errors[f] = ""
	} else if text, ok := result.(string); ok {
		m.errors[f] = text
	}
	return m
}

func (m *ServerInfoModel) refocus() {
	for i := range m.inputs {
		if serverField(i) == m.focused {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m ServerInfoModel) View() string {
	labels := [fieldCount]string{"HOST ALIAS", "HOSTNAME OR IP ADDRESS", "USERNAME", "PORT"}
	hints := [fieldCount]string{
		"Short nickname used in ssh commands  (e.g. prod, dev-box)",
		"Server IP or domain  (e.g. 192.168.1.1, server.example.com)",
		"The Linux user account on the remote server",
		"Default is 22",
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("SERVER INFORMATION", 1, 7))

	for i := serverField(0); i < fieldCount; i++ {
		if serverField(i) < m.focused {
			b.WriteString(ui.RenderConfirmedRow(labels[i], m.inputs[i].Value()) + "\n")
			continue
		}
		b.WriteString(fmt.Sprintf("\n  %s\n", ui.Label(labels[i])))
		b.WriteString("  " + m.inputs[i].View() + "\n")
		if m.errors[i] != "" {
			b.WriteString("  " + ui.Error.Render(ui.Sym.Error+" "+m.errors[i]) + "\n")
		} else {
			b.WriteString(ui.RenderHint(hints[i]) + "\n")
		}
		if serverField(i) == m.focused {
			break
		}
	}

	b.WriteString(ui.RenderNavHint([]string{"↑↓ / Tab navigate", "⏎ confirm", "Ctrl+C quit"}))
	return b.String()
}
