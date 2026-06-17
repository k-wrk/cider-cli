package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateDocumentsReport(key string) (model, tea.Cmd) {
	switch key {
	case "left", "p":
		if m.reportPage > 0 {
			m.reportPage--
		}
	case "right", "n":
		totalPages := (len(m.largeFiles) + m.itemsPerPage - 1) / m.itemsPerPage
		if m.reportPage < totalPages-1 {
			m.reportPage++
		}
	case "esc", "b":
		m.state = stateMainMenu
	}
	return m, nil
}
