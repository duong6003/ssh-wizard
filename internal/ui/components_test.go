package ui_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/duong6003/ssh-wizard/internal/ui"
)

func TestRenderProgress(t *testing.T) {
	ui.ConfigureTerminalSymbols(true)
	result := ui.StripANSI(ui.RenderProgress(1, 7))
	assert.Contains(t, result, "▪▪▫▫▫▫▫▫▫▫▫▫▫▫")
	assert.Contains(t, result, "1 / 7")

	result = ui.StripANSI(ui.RenderProgress(7, 7))
	assert.Contains(t, result, "▪▪▪▪▪▪▪▪▪▪▪▪▪▪")
}

func TestRenderConfirmedRow(t *testing.T) {
	ui.ConfigureTerminalSymbols(true)
	result := ui.StripANSI(ui.RenderConfirmedRow("HOST ALIAS", "my-server"))
	assert.Contains(t, result, "✓")
	assert.Contains(t, result, "HOST ALIAS")
	assert.Contains(t, result, "my-server")
}

func TestRenderWarningBlock(t *testing.T) {
	ui.ConfigureTerminalSymbols(true)
	result := ui.StripANSI(ui.RenderWarningBlock("file exists", "please check"))
	assert.Contains(t, result, "⚠")
	assert.Contains(t, result, "file exists")
	assert.Contains(t, result, "please check")
}

func TestASCIIFallbackSymbols(t *testing.T) {
	ui.ConfigureTerminalSymbols(false)
	t.Cleanup(func() { ui.ConfigureTerminalSymbols(true) })

	progress := ui.StripANSI(ui.RenderProgress(1, 7))
	assert.Contains(t, progress, "##------------")
	row := ui.StripANSI(ui.RenderConfirmedRow("status", "ready"))
	assert.Contains(t, row, "OK")
	box := ui.StripANSI(ui.RenderConfigBox("Host prod"))
	assert.Contains(t, box, "+")
	assert.NotContains(t, box, "╭")
}

func TestRenderConfigBoxPadsRawTextBeforeANSI(t *testing.T) {
	ui.ConfigureTerminalSymbols(true)
	box := ui.StripANSI(ui.RenderConfigBox("Host prod\nUser ubuntu"))
	lines := strings.Split(box, "\n")

	assert.Len(t, lines, 4)
	width := len([]rune(lines[0]))
	for _, line := range lines[1:] {
		assert.Equal(t, width, len([]rune(line)))
	}
}
