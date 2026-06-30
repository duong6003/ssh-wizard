package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Display     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	Primary     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	Secondary   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	Success     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Warning     = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	Error       = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	Interactive = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

func Label(text string) string {
	return Secondary.Render(strings.ToUpper(text))
}

var Sym = struct {
	Prompt   string
	Success  string
	Error    string
	Warning  string
	Bullet   string
	Active   string
	Inactive string
}{
	Prompt:   "›",
	Success:  "✓",
	Error:    "✗",
	Warning:  "⚠",
	Bullet:   "·",
	Active:   "●",
	Inactive: "○",
}
