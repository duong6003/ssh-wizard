package ui

import (
	"fmt"
	"strings"
)

var dividerLine = Secondary.Render(strings.Repeat("─", 61))

func RenderHeader(title string, step, total int) string {
	progress := RenderProgress(step, total)
	titleLine := fmt.Sprintf("  %s", Display.Render(strings.ToUpper(title)))
	return strings.Join([]string{
		"",
		progress,
		"",
		"",
		titleLine,
		"  " + dividerLine,
		"",
	}, "\n")
}

func RenderDivider() string {
	return "\n  " + dividerLine + "\n"
}

func RenderKeyValue(labelText, value string) string {
	lbl := fmt.Sprintf("%-*s", labelWidth, strings.ToUpper(labelText))
	return fmt.Sprintf("  %s%s", Secondary.Render(lbl), Display.Render(value))
}

func RenderNavHint(hints []string) string {
	return "\n  " + dividerLine + "\n  " + Secondary.Render(strings.Join(hints, "   ")) + "\n"
}
