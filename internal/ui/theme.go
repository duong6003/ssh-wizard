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

type Symbols struct {
	Prompt        string
	Success       string
	Error         string
	Warning       string
	Bullet        string
	Active        string
	Inactive      string
	ProgressFull  string
	ProgressEmpty string
	Horizontal    string
	Vertical      string
	TopLeft       string
	TopRight      string
	BottomLeft    string
	BottomRight   string
	Arrow         string
	Enter         string
	UpDown        string
	Ellipsis      string
	PasswordEcho  rune
}

var unicodeSymbols = Symbols{
	Prompt:        "\u203a",
	Success:       "\u2713",
	Error:         "\u2717",
	Warning:       "\u26a0",
	Bullet:        "\u00b7",
	Active:        "\u25cf",
	Inactive:      "\u25cb",
	ProgressFull:  "\u25aa",
	ProgressEmpty: "\u25ab",
	Horizontal:    "\u2500",
	Vertical:      "\u2502",
	TopLeft:       "\u256d",
	TopRight:      "\u256e",
	BottomLeft:    "\u2570",
	BottomRight:   "\u256f",
	Arrow:         "\u2192",
	Enter:         "\u23ce",
	UpDown:        "\u2191\u2193",
	Ellipsis:      "\u00b7\u00b7\u00b7",
	PasswordEcho:  '\u00b7',
}

var asciiSymbols = Symbols{
	Prompt:        ">",
	Success:       "OK",
	Error:         "X",
	Warning:       "!",
	Bullet:        "-",
	Active:        "*",
	Inactive:      "o",
	ProgressFull:  "#",
	ProgressEmpty: "-",
	Horizontal:    "-",
	Vertical:      "|",
	TopLeft:       "+",
	TopRight:      "+",
	BottomLeft:    "+",
	BottomRight:   "+",
	Arrow:         "->",
	Enter:         "Enter",
	UpDown:        "Up/Down",
	Ellipsis:      "...",
	PasswordEcho:  '*',
}

var Sym = unicodeSymbols
var unicodeEnabled = true

func ConfigureTerminalSymbols(supportsUnicode bool) {
	unicodeEnabled = supportsUnicode
	if supportsUnicode {
		Sym = unicodeSymbols
		return
	}
	Sym = asciiSymbols
}

func NormalizeGlyphs(text string) string {
	if unicodeEnabled {
		return text
	}
	replacements := []struct {
		from string
		to   string
	}{
		{"\u2191\u2193", Sym.UpDown},
		{"\u23ce", Sym.Enter},
		{"\u2192", Sym.Arrow},
		{"\u00b7\u00b7\u00b7", Sym.Ellipsis},
		{"\u2014", "-"},
		{"\u00b7", Sym.Bullet},
	}
	for _, r := range replacements {
		text = strings.ReplaceAll(text, r.from, r.to)
	}
	return text
}
