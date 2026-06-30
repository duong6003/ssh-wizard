package steps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/ssh"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/utils"
	"ssh-wizard/internal/wizard"
)

type authPhase int

const (
	authPhaseChoosing authPhase = iota
	authPhaseExisting
	authPhaseGenerate
	authPhasePassword
	authPhaseDone
)

type authKeyGenPhase int

const (
	genPhaseType authKeyGenPhase = iota
	genPhaseOutputPath
	genPhasePassphrase
	genPhaseGenerating
	genPhaseDone
)

type AuthModel struct {
	state         *wizard.State
	phase         authPhase
	cursor        int
	keyPathInput  textinput.Model
	keyPathErr    string
	validating    bool
	spinner       spinner.Model
	genPhase      authKeyGenPhase
	genTypeCursor int
	genPathInput  textinput.Model
	genPassInput  textinput.Model
	genKeyType    ssh.KeyType
	genOutputPath string
	validatedKey  *ssh.ValidatedKey
	generated     *ssh.GeneratedKey
	genErr        string
}

type keyValidatedMsg struct {
	key *ssh.ValidatedKey
	err error
}

type keyGeneratedMsg struct {
	key *ssh.GeneratedKey
	err error
}

func NewAuth(state *wizard.State) tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	keyPath := textinput.New()
	keyPath.Placeholder = "~/.ssh/id_ed25519"
	keyPath.CharLimit = 256

	genPath := textinput.New()
	genPath.CharLimit = 256

	passphrase := textinput.New()
	passphrase.EchoMode = textinput.EchoPassword
	passphrase.EchoCharacter = '·'
	passphrase.Placeholder = "leave blank to skip"

	return AuthModel{
		state:        state,
		spinner:      s,
		keyPathInput: keyPath,
		genPathInput: genPath,
		genPassInput: passphrase,
		genKeyType:   ssh.KeyTypeED25519,
	}
}

func (m AuthModel) Init() tea.Cmd { return m.spinner.Tick }

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case keyValidatedMsg:
		m.validating = false
		if msg.err != nil {
			m.keyPathErr = msg.err.Error()
			m.keyPathInput.Focus()
			return m, nil
		}
		m.validatedKey = msg.key
		m.phase = authPhaseDone
		path := utils.ExpandTilde(m.keyPathInput.Value())
		m.state.Key = &wizard.KeyConfig{
			Method:         wizard.AuthExisting,
			KeyType:        msg.key.KeyType,
			PrivateKeyPath: path,
			PublicKeyPath:  path + ".pub",
			Fingerprint:    msg.key.Fingerprint,
			HasPassphrase:  msg.key.HasPassphrase,
		}
		return m, func() tea.Msg { return wizard.AdvanceMsg{} }

	case keyGeneratedMsg:
		if msg.err != nil {
			m.genErr = msg.err.Error()
			m.genPhase = genPhaseType
			return m, nil
		}
		m.generated = msg.key
		m.genPhase = genPhaseDone
		m.state.Key = &wizard.KeyConfig{
			Method:         wizard.AuthGenerate,
			KeyType:        msg.key.KeyType,
			PrivateKeyPath: msg.key.PrivateKeyPath,
			PublicKeyPath:  msg.key.PublicKeyPath,
			Fingerprint:    msg.key.Fingerprint,
			HasPassphrase:  m.genPassInput.Value() != "",
		}
		return m, func() tea.Msg { return wizard.AdvanceMsg{} }

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var spinCmd tea.Cmd
	m.spinner, spinCmd = m.spinner.Update(msg)
	var cmd1, cmd2, cmd3 tea.Cmd
	m.keyPathInput, cmd1 = m.keyPathInput.Update(msg)
	m.genPathInput, cmd2 = m.genPathInput.Update(msg)
	m.genPassInput, cmd3 = m.genPassInput.Update(msg)
	return m, tea.Batch(spinCmd, cmd1, cmd2, cmd3)
}

func (m AuthModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.phase {
	case authPhaseChoosing:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < 2 {
				m.cursor++
			}
		case tea.KeyEnter:
			switch m.cursor {
			case 0:
				m.phase = authPhaseExisting
				m.keyPathInput.Focus()
			case 1:
				m.phase = authPhaseGenerate
				m.genPhase = genPhaseType
			case 2:
				m.phase = authPhasePassword
			}
		}

	case authPhaseExisting:
		if msg.Type == tea.KeyEnter && !m.validating {
			path := utils.ExpandTilde(m.keyPathInput.Value())
			if path == "" {
				path = utils.ExpandTilde("~/.ssh/id_ed25519")
				m.keyPathInput.SetValue("~/.ssh/id_ed25519")
			}
			m.validating = true
			m.keyPathErr = ""
			return m, func() tea.Msg {
				key, err := ssh.ValidateKeyFile(path)
				return keyValidatedMsg{key: key, err: err}
			}
		}

	case authPhaseGenerate:
		switch m.genPhase {
		case genPhaseType:
			switch msg.Type {
			case tea.KeyUp:
				m.genTypeCursor = 0
			case tea.KeyDown:
				m.genTypeCursor = 1
			case tea.KeyEnter:
				if m.genTypeCursor == 0 {
					m.genKeyType = ssh.KeyTypeED25519
				} else {
					m.genKeyType = ssh.KeyTypeRSA
				}
				defaultPath := utils.GetDefaultKeyPath(m.state.Server.Alias, string(m.genKeyType))
				m.genPathInput.SetValue(defaultPath)
				m.genPathInput.Focus()
				m.genPhase = genPhaseOutputPath
			}
		case genPhaseOutputPath:
			if msg.Type == tea.KeyEnter {
				m.genOutputPath = utils.ExpandTilde(m.genPathInput.Value())
				m.genPassInput.Focus()
				m.genPhase = genPhasePassphrase
			}
		case genPhasePassphrase:
			if msg.Type == tea.KeyEnter {
				m.genPhase = genPhaseGenerating
				passphrase := m.genPassInput.Value()
				outputPath := m.genOutputPath
				keyType := m.genKeyType
				return m, func() tea.Msg {
					key, err := ssh.GenerateKey(keyType, outputPath, passphrase)
					return keyGeneratedMsg{key: key, err: err}
				}
			}
		}

	case authPhasePassword:
		if msg.Type == tea.KeyEnter {
			defaultPath := utils.ExpandTilde("~/.ssh/id_ed25519")
			key, err := ssh.ValidateKeyFile(defaultPath)
			if err != nil {
				m.phase = authPhaseGenerate
				m.genPhase = genPhaseType
				return m, nil
			}
			m.state.Key = &wizard.KeyConfig{
				Method:         wizard.AuthPassword,
				KeyType:        key.KeyType,
				PrivateKeyPath: defaultPath,
				PublicKeyPath:  defaultPath + ".pub",
				Fingerprint:    key.Fingerprint,
				HasPassphrase:  key.HasPassphrase,
			}
			return m, func() tea.Msg { return wizard.AdvanceMsg{} }
		}
	}
	return m, nil
}

func (m AuthModel) View() string {
	var b strings.Builder
	b.WriteString(ui.RenderHeader("AUTHENTICATION", 2, 7))

	switch m.phase {
	case authPhaseChoosing:
		b.WriteString("\n  " + ui.Primary.Render("How should this server authenticate you?") + "\n\n")
		choices := []string{
			"Use an existing SSH key",
			"Generate a new SSH key pair",
			"Use password to install a key",
		}
		for i, choice := range choices {
			if i == m.cursor {
				b.WriteString("  " + ui.Display.Render(ui.Sym.Active) + "  " + ui.Display.Render(choice) + "\n")
			} else {
				b.WriteString("  " + ui.Secondary.Render(ui.Sym.Inactive) + "  " + ui.Secondary.Render(choice) + "\n")
			}
		}
		b.WriteString(ui.RenderNavHint([]string{"↑↓ select", "⏎ confirm"}))

	case authPhaseExisting:
		b.WriteString(ui.RenderConfirmedRow("METHOD", "Existing SSH key") + "\n")
		b.WriteString("\n  " + ui.Label("private key path") + "\n")
		if m.validating {
			b.WriteString("  " + m.spinner.View() + " " + ui.Secondary.Render("Validating key…") + "\n")
		} else {
			b.WriteString("  " + m.keyPathInput.View() + "\n")
			if m.keyPathErr != "" {
				b.WriteString("  " + ui.Error.Render(ui.Sym.Error+"  "+m.keyPathErr) + "\n")
			} else {
				b.WriteString(ui.RenderHint("Leave blank to use default: ~/.ssh/id_ed25519") + "\n")
			}
		}

	case authPhaseGenerate:
		b.WriteString(ui.RenderConfirmedRow("METHOD", "Generate new key") + "\n")
		m.viewGenerate(&b)

	case authPhasePassword:
		b.WriteString(ui.RenderConfirmedRow("METHOD", "Password + install key") + "\n")
		b.WriteString("\n  " + ui.Primary.Render("We'll use your password once to install your SSH key,") + "\n")
		b.WriteString("  " + ui.Primary.Render("then switch to key-only authentication.") + "\n\n")
		b.WriteString("  " + ui.Secondary.Render("We'll use ~/.ssh/id_ed25519 if it exists, or generate a new key.") + "\n\n")
		b.WriteString(ui.RenderNavHint([]string{"⏎ continue"}))
	}

	return b.String()
}

func (m AuthModel) viewGenerate(b *strings.Builder) {
	switch m.genPhase {
	case genPhaseType:
		if m.genErr != "" {
			b.WriteString("\n  " + ui.Error.Render(ui.Sym.Error+"  "+m.genErr) + "\n")
		}
		b.WriteString("\n  " + ui.Label("key type") + "\n\n")
		types := []struct{ label, desc string }{
			{"ED25519", "Recommended. Fast, secure, small."},
			{"RSA 4096", "Maximum compatibility with older servers."},
		}
		for i, keyType := range types {
			sym := ui.Sym.Inactive
			style := ui.Secondary
			if i == m.genTypeCursor {
				sym = ui.Sym.Active
				style = ui.Display
			}
			b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
				style.Render(sym),
				style.Render(keyType.label),
				ui.Secondary.Render(keyType.desc),
			))
		}
		b.WriteString(ui.RenderNavHint([]string{"↑↓ select", "⏎ confirm"}))

	case genPhaseOutputPath:
		b.WriteString(ui.RenderConfirmedRow("KEY TYPE", string(m.genKeyType)) + "\n")
		b.WriteString("\n  " + ui.Label("save as") + "\n")
		b.WriteString("  " + m.genPathInput.View() + "\n")
		b.WriteString(ui.RenderHint("Press Enter to accept default") + "\n")

	case genPhasePassphrase:
		b.WriteString(ui.RenderConfirmedRow("KEY TYPE", string(m.genKeyType)) + "\n")
		b.WriteString(ui.RenderConfirmedRow("SAVE AS", m.genPathInput.Value()) + "\n")
		b.WriteString("\n  " + ui.Label("passphrase") + "  " + ui.Secondary.Render("(optional)") + "\n")
		b.WriteString("  " + m.genPassInput.View() + "\n")
		b.WriteString(ui.RenderHint("Press Enter to skip. Recommended for shared machines.") + "\n")

	case genPhaseGenerating:
		b.WriteString(ui.RenderConfirmedRow("KEY TYPE", string(m.genKeyType)) + "\n")
		b.WriteString("\n  " + m.spinner.View() + " " + ui.Secondary.Render("Generating key…") + "\n")
	}
}
