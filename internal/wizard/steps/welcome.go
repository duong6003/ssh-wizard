package steps

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/wizard"
)

type WelcomeModel struct{}

func NewWelcome() tea.Model { return WelcomeModel{} }

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
	b.WriteString("  " + ui.Display.Render("ssh-wizard") + "\n\n\n")
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
	b.WriteString("  " + ui.Secondary.Render(strings.Repeat("─", 61)) + "\n")
	b.WriteString("  " + ui.Secondary.Render("Press Enter to begin   Ctrl+C to quit") + "\n")
	return b.String()
}
