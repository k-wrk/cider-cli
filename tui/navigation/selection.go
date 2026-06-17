package navigation

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateSelection(key string) (model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		start := m.selectionPage * m.itemsPerPage
		end := start + m.itemsPerPage
		if end > len(m.items) {
			end = len(m.items)
		}
		pageSize := end - start
		if m.cursor < pageSize-1 {
			m.cursor++
		}
	case "left", "p":
		if m.selectionPage > 0 {
			m.selectionPage--
			m.cursor = 0
		}
	case "right", "n":
		totalPages := (len(m.items) + m.itemsPerPage - 1) / m.itemsPerPage
		if m.selectionPage < totalPages-1 {
			m.selectionPage++
			m.cursor = 0
		}
	case " ":
		start := m.selectionPage * m.itemsPerPage
		itemIdx := start + m.cursor
		if itemIdx < len(m.items) {
			m.items[itemIdx].Selected = !m.items[itemIdx].Selected
		}
	case "esc", "b":
		m.state = stateMainMenu
	case "enter":
		m.itemsToClean = []int{}
		var itemNames []string
		for i, item := range m.items {
			if item.Selected && item.Size > 0 {
				m.itemsToClean = append(m.itemsToClean, i)
				itemNames = append(itemNames, item.Name)
			}
		}
		m.freedSize = 0
		m.currentCleanIdx = 0

		if len(m.itemsToClean) == 0 {
			m.state = stateFinished
			return m, nil
		}

		m.previousState = stateSelection
		m.confirmPrompt = fmt.Sprintf("Do you want to move the following selected items to the trash?\n\n- %s", strings.Join(itemNames, "\n- "))
		m.onConfirm = func() (model, tea.Cmd) {
			m.state = stateCleaning
			cmd := m.progress.SetPercent(0)
			return m, tea.Batch(cmd, m.cleanNext())
		}
		m.state = stateConfirmTrash
	}
	return m, nil
}
