package ui

import (
	"fmt"
	"regexp"
	"strings"
)

const totalSegments = 14
const labelWidth = 22

func RenderProgress(currentStep, totalSteps int) string {
	filled := (currentStep * totalSegments) / totalSteps
	bar := Display.Render(strings.Repeat("▪", filled)) +
		Secondary.Render(strings.Repeat("▫", totalSegments-filled))
	label := Secondary.Render(fmt.Sprintf("%d / %d", currentStep, totalSteps))
	return fmt.Sprintf("  %s   %s", bar, label)
}

func RenderConfirmedRow(labelText, value string) string {
	lbl := fmt.Sprintf("%-*s", labelWidth, strings.ToUpper(labelText))
	return fmt.Sprintf("  %s  %s%s",
		Success.Render(Sym.Success),
		Secondary.Render(lbl),
		Display.Render(value),
	)
}

func RenderHint(text string) string {
	return "  " + Secondary.Render(text)
}

func RenderWarningBlock(title, body string) string {
	div := Secondary.Render(strings.Repeat("─", 61))
	icon := Warning.Render(Sym.Warning)
	lines := make([]string, 0)
	for _, line := range strings.Split(body, "\n") {
		lines = append(lines, "     "+Primary.Render(line))
	}
	return strings.Join([]string{
		div,
		fmt.Sprintf("  %s  %s", icon, Warning.Render(title)),
		"",
		strings.Join(lines, "\n"),
		div,
	}, "\n")
}

func RenderConfigBox(content string) string {
	lines := strings.Split(content, "\n")
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	width += 4
	top := "  ╭" + strings.Repeat("─", width) + "╮"
	bottom := "  ╰" + strings.Repeat("─", width) + "╯"
	middle := make([]string, 0, len(lines))
	for _, line := range lines {
		middle = append(middle, fmt.Sprintf("  │  %-*s  │", width-4, Secondary.Render(line)))
	}
	return strings.Join(append(append([]string{top}, middle...), bottom), "\n")
}

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func StripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}
