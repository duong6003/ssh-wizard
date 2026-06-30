package steps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/duong6003/ssh-wizard/internal/ui"
)

func TestAuthViewRendersStatusMessage(t *testing.T) {
	ui.ConfigureTerminalSymbols(false)
	t.Cleanup(func() { ui.ConfigureTerminalSymbols(true) })

	model := AuthModel{
		phase:     authPhaseChoosing,
		statusMsg: "~/.ssh/id_ed25519 not found — generating a new key instead.",
	}

	view := ui.StripANSI(model.View())

	assert.Contains(t, view, "!")
	assert.Contains(t, view, "~/.ssh/id_ed25519 not found - generating a new key instead.")
}
