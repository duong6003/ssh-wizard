package steps

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/utils"
	"ssh-wizard/internal/wizard"
)

type DoneModel struct{ state *wizard.State }

func NewDone(state *wizard.State) tea.Model { return DoneModel{state: state} }

func (m DoneModel) Init() tea.Cmd { return nil }

func (m DoneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyEnter || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DoneModel) View() string {
	var b strings.Builder
	b.WriteString(ui.RenderHeader("READY", 7, 7))
	b.WriteString("\n  " + ui.Display.Render(m.state.Server.Alias) + " " + ui.Secondary.Render("is configured.") + "\n")
	b.WriteString(ui.RenderDivider())

	b.WriteString(ui.RenderKeyValue("HOST ALIAS", m.state.Server.Alias) + "\n")
	b.WriteString(ui.RenderKeyValue("SERVER", m.state.Server.Username+"@"+m.state.Server.Hostname) + "\n")
	b.WriteString(ui.RenderKeyValue("PORT", fmt.Sprintf("%d", m.state.Server.Port)) + "\n")
	if m.state.Key != nil && m.state.Key.PrivateKeyPath != "" {
		b.WriteString(ui.RenderKeyValue("KEY", m.state.Key.PrivateKeyPath) + "\n")
	}
	b.WriteString(ui.RenderKeyValue("CONFIG", utils.GetSSHConfigPath()) + "\n")

	b.WriteString(ui.RenderDivider())
	b.WriteString("  " + ui.Secondary.Render("Connect in terminal:") + "\n\n")
	b.WriteString("    " + ui.Interactive.Render("ssh "+m.state.Server.Alias) + "\n\n")
	b.WriteString("  " + ui.Secondary.Render("Connect in VS Code:") + "\n\n")
	b.WriteString("    " + ui.Secondary.Render("Cmd+Shift+P") + " " + ui.Sym.Arrow + " " + ui.Interactive.Render("Remote-SSH: Connect to Host") + " " + ui.Sym.Arrow + " " + ui.Interactive.Render(m.state.Server.Alias) + "\n")
	b.WriteString(ui.RenderNavHint([]string{"⏎ exit", "A add another server"}))
	return b.String()
}
