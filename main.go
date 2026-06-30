package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duong6003/ssh-wizard/internal/ui"
	"github.com/duong6003/ssh-wizard/internal/utils"
	"github.com/duong6003/ssh-wizard/internal/wizard"
	_ "github.com/duong6003/ssh-wizard/internal/wizard/steps"
)

func main() {
	env := utils.CheckEnvironment()
	ui.ConfigureTerminalSymbols(env.SupportsUnicode)

	p := tea.NewProgram(wizard.NewRootModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
