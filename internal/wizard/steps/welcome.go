package steps

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duong6003/ssh-wizard/internal/ui"
	"github.com/duong6003/ssh-wizard/internal/utils"
	"github.com/duong6003/ssh-wizard/internal/wizard"
)

type WelcomeModel struct {
	env utils.EnvironmentReport
}

func NewWelcome() tea.Model { return WelcomeModel{env: utils.CheckEnvironment()} }

func (m WelcomeModel) Init() tea.Cmd { return nil }

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return m, func() tea.Msg { return wizard.AdvanceMsg{} }
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString("  " + ui.Display.Render("ssh-wizard") + "  " + ui.Secondary.Render(wizard.Version) + "\n\n\n")
	b.WriteString("  " + ui.Primary.Render("Set up VS Code Remote SSH in under 5 minutes.") + "\n\n")
	b.WriteString("  " + ui.Primary.Render("This wizard will:") + "\n\n")
	for _, item := range []string{
		"Manage your SSH keys",
		"Install your public key on the remote server",
		"Generate a clean ~/.ssh/config entry",
		"Test the connection",
		"Confirm VS Code is ready to connect",
	} {
		b.WriteString("  " + ui.Secondary.Render(ui.Sym.Bullet) + "  " + ui.Primary.Render(item) + "\n")
	}
	b.WriteString("\n  " + ui.Secondary.Render("No manual file editing required.") + "\n\n")
	b.WriteString("  " + ui.Secondary.Render(strings.Repeat(ui.Sym.Horizontal, 61)) + "\n")
	b.WriteString("  " + ui.Label("environment") + "\n")
	b.WriteString("  " + ui.Secondary.Render("Client: "+string(m.env.Platform)+" / "+m.env.Terminal) + "\n")
	if m.env.SSHAvailable && m.env.SSHKeygenAvailable && len(m.env.Warnings) == 0 {
		b.WriteString("  " + ui.Success.Render(ui.Sym.Success+"  ssh and ssh-keygen ready") + "\n")
	} else {
		for _, warning := range m.env.Warnings {
			b.WriteString("  " + ui.Warning.Render(ui.Sym.Warning+"  "+ui.NormalizeGlyphs(warning)) + "\n")
		}
	}
	b.WriteString("  " + ui.Secondary.Render(ui.NormalizeGlyphs(m.env.RemoteTargetNote)) + "\n\n")
	b.WriteString("  " + ui.Secondary.Render(strings.Repeat(ui.Sym.Horizontal, 61)) + "\n")
	b.WriteString("  " + ui.Secondary.Render("Press Enter to begin   Ctrl+C to quit") + "\n")
	return b.String()
}
