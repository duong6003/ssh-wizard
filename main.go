package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"ssh-wizard/internal/wizard"
	_ "ssh-wizard/internal/wizard/steps"
)

func main() {
	p := tea.NewProgram(wizard.NewRootModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
