package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateConfirmTrash(key string) (model, tea.Cmd) {
	switch key {
	case "y", "s", "enter":
		if m.onConfirm != nil {
			return m.onConfirm()
		}
	case "n", "esc", "b":
		m.state = m.previousState
	}
	return m, nil
}
