package steps

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"ssh-wizard/internal/wizard"
)

func TestDoneAddAnotherServerResetsStateAndRestarts(t *testing.T) {
	state := &wizard.State{
		Server:            &wizard.ServerConfig{Alias: "prod", Hostname: "example.com", Username: "ubuntu", Port: 22},
		ConfigWritten:     true,
		ConnectionSuccess: true,
	}
	model := DoneModel{state: state}

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.NotNil(t, cmd)
	assert.Nil(t, state.Server)
	assert.False(t, state.ConfigWritten)
	assert.False(t, state.ConnectionSuccess)
	assert.IsType(t, wizard.RestartMsg{}, cmd())
}
