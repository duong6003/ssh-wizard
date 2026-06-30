package steps

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	gossh "github.com/duong6003/ssh-wizard/internal/ssh"
	"github.com/duong6003/ssh-wizard/internal/ui"
	"github.com/duong6003/ssh-wizard/internal/utils"
	"github.com/duong6003/ssh-wizard/internal/wizard"
)

type configPhase int

const (
	configPhasePreview configPhase = iota
	configPhaseConflict
	configPhaseWriting
	configPhaseDone
)

type configWrittenMsg struct{ err error }

type SSHConfigStepModel struct {
	state    *wizard.State
	phase    configPhase
	entry    gossh.SSHConfigEntry
	existing string
	conflict *gossh.SSHConfigEntry
	preview  string
	cursor   int
	writeErr string
}

func NewSSHConfigStep(state *wizard.State) tea.Model {
	m := SSHConfigStepModel{state: state}

	existing, _, _ := gossh.ReadConfig()
	m.existing = existing

	m.entry = gossh.SSHConfigEntry{
		Host:                state.Server.Alias,
		Hostname:            state.Server.Hostname,
		User:                state.Server.Username,
		Port:                state.Server.Port,
		ServerAliveInterval: 60,
	}
	if state.Key != nil && state.Key.Method != wizard.AuthPassword {
		m.entry.IdentityFile = state.Key.PrivateKeyPath
		m.entry.IdentitiesOnly = true
	}

	entries := gossh.ParseSSHConfig(existing)
	if conflict := gossh.FindConflict(entries, state.Server.Alias); conflict != nil {
		m.conflict = conflict
		m.phase = configPhaseConflict
	}

	m.preview = gossh.FormatEntry(m.entry)
	return m
}

func (m SSHConfigStepModel) Init() tea.Cmd { return nil }

func (m SSHConfigStepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case configWrittenMsg:
		if msg.err != nil {
			m.writeErr = msg.err.Error()
			return m, nil
		}
		m.phase = configPhaseDone
		m.state.ConfigWritten = true
		return m, func() tea.Msg { return wizard.AdvanceMsg{} }

	case tea.KeyMsg:
		switch m.phase {
		case configPhaseConflict:
			switch msg.Type {
			case tea.KeyUp:
				if m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyDown:
				if m.cursor < 1 {
					m.cursor++
				}
			case tea.KeyEnter:
				if m.cursor == 0 {
					m.phase = configPhasePreview
				} else {
					return m, func() tea.Msg { return wizard.AdvanceMsg{} }
				}
			}
		case configPhasePreview:
			if msg.Type == tea.KeyEnter {
				m.phase = configPhaseWriting
				content := gossh.AppendEntry(m.existing, m.entry)
				return m, func() tea.Msg {
					err := gossh.WriteConfig(content)
					return configWrittenMsg{err: err}
				}
			}
		}
	}
	return m, nil
}

func (m SSHConfigStepModel) View() string {
	var b strings.Builder
	b.WriteString(ui.RenderHeader("SSH CONFIGURATION", 4, 7))

	switch m.phase {
	case configPhaseConflict:
		b.WriteString(ui.RenderWarningBlock(
			`Host "`+m.state.Server.Alias+`" already exists in ~/.ssh/config`,
			"Your new entry will replace the existing one.\n\nExisting:\n"+gossh.FormatEntry(*m.conflict),
		))
		b.WriteString("\n")
		choices := []string{"Overwrite existing entry", "Skip (keep existing)"}
		for i, choice := range choices {
			if i == m.cursor {
				b.WriteString("  " + ui.Display.Render(ui.Sym.Active+"  "+choice) + "\n")
			} else {
				b.WriteString("  " + ui.Secondary.Render(ui.Sym.Inactive+"  "+choice) + "\n")
			}
		}
		b.WriteString(ui.RenderNavHint([]string{"↑↓ select", "⏎ confirm"}))

	case configPhasePreview:
		b.WriteString("\n  " + ui.Secondary.Render(ui.NormalizeGlyphs("Preview — entry to be written to ~/.ssh/config:")) + "\n\n")
		b.WriteString(ui.RenderConfigBox(m.preview) + "\n")
		if m.writeErr != "" {
			b.WriteString("\n  " + ui.Error.Render(ui.Sym.Error+"  "+m.writeErr) + "\n")
		}
		b.WriteString(ui.RenderNavHint([]string{"⏎ write to file", "Ctrl+C abort"}))

	case configPhaseWriting:
		b.WriteString("\n  " + ui.Secondary.Render("Writing to ~/.ssh/config "+ui.Sym.Ellipsis) + "\n")

	case configPhaseDone:
		b.WriteString(ui.RenderConfirmedRow("WRITTEN TO", utils.GetSSHConfigPath()) + "\n")
	}

	return b.String()
}
