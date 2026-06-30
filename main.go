package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/ui"
	"ssh-wizard/internal/utils"
	"ssh-wizard/internal/wizard"
	_ "ssh-wizard/internal/wizard/steps"
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
