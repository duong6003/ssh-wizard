package ui_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ssh-wizard/internal/ui"
)

func TestRenderProgress(t *testing.T) {
	result := ui.StripANSI(ui.RenderProgress(1, 7))
	assert.Contains(t, result, "▪▪▫▫▫▫▫▫▫▫▫▫▫▫")
	assert.Contains(t, result, "1 / 7")

	result = ui.StripANSI(ui.RenderProgress(7, 7))
	assert.Contains(t, result, "▪▪▪▪▪▪▪▪▪▪▪▪▪▪")
}

func TestRenderConfirmedRow(t *testing.T) {
	result := ui.StripANSI(ui.RenderConfirmedRow("HOST ALIAS", "my-server"))
	assert.Contains(t, result, "✓")
	assert.Contains(t, result, "HOST ALIAS")
	assert.Contains(t, result, "my-server")
}

func TestRenderWarningBlock(t *testing.T) {
	result := ui.StripANSI(ui.RenderWarningBlock("file exists", "please check"))
	assert.Contains(t, result, "⚠")
	assert.Contains(t, result, "file exists")
	assert.Contains(t, result, "please check")
}
