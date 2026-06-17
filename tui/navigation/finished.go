package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateFinished(key string) (model, tea.Cmd) {
	switch key {
	case "q":
		return m, tea.Quit
	case "esc", "b", "enter":
		m.state = stateMainMenu
		return m, tea.Batch(m.startScanning(), m.spinner.Tick)
	}
	return m, nil
}
