package ui

import (
	"fmt"
	"strings"
)

func dividerLine() string {
	return Secondary.Render(strings.Repeat(Sym.Horizontal, 61))
}

func RenderHeader(title string, step, total int) string {
	progress := RenderProgress(step, total)
	titleLine := fmt.Sprintf("  %s", Display.Render(strings.ToUpper(title)))
	return strings.Join([]string{
		"",
		progress,
		"",
		"",
		titleLine,
		"  " + dividerLine(),
		"",
	}, "\n")
}

func RenderDivider() string {
	return "\n  " + dividerLine() + "\n"
}

func RenderKeyValue(labelText, value string) string {
	lbl := fmt.Sprintf("%-*s", labelWidth, strings.ToUpper(labelText))
	return fmt.Sprintf("  %s%s", Secondary.Render(lbl), Display.Render(value))
}

func RenderNavHint(hints []string) string {
	normalized := make([]string, len(hints))
	for i, hint := range hints {
		normalized[i] = NormalizeGlyphs(hint)
	}
	return "\n  " + dividerLine() + "\n  " + Secondary.Render(strings.Join(normalized, "   ")) + "\n"
}
