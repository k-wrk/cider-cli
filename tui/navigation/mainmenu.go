package navigation

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateMainMenu(key string) (model, tea.Cmd) {
	m.statusMsg = ""
	switch key {
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		}
	case "down", "j":
		if m.menuCursor < len(menuOptions)-1 {
			m.menuCursor++
		}
	case "enter":
		switch m.menuCursor {
		case 0:
			m.state = stateSelection
			m.cursor = 0
		case 1:
			m.state = stateScanningReport
			return m, tea.Batch(m.scanDocumentsReport(), m.spinner.Tick)
		case 2:
			var snapshotsSize int64
			for _, item := range m.items {
				if item.ID == "snapshots" {
					snapshotsSize = item.Size
					break
				}
			}
			m.previousState = stateMainMenu
			m.confirmPrompt = "Do you want to delete Time Machine Local Snapshots (requires sudo privileges)?"
			m.onConfirm = func() (model, tea.Cmd) {
				c := exec.Command("sudo", "tmutil", "deletelocalsnapshots", "/")
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					if err == nil {
						return cleanFinishedMsg{freedBytes: snapshotsSize}
					}
					return scanFinishedMsg{}
				})
			}
			m.state = stateConfirmTrash
			return m, nil
		case 3:
			m.state = stateScanningApps
			return m, tea.Batch(m.scanAppsSuggestions(), m.spinner.Tick)
		case 4:
			m.state = stateScanningHF
			return m, tea.Batch(m.scanHuggingFace(), m.spinner.Tick)
		case 5:
			m.state = stateScanningOllama
			return m, tea.Batch(m.scanOllama(), m.spinner.Tick)
		case 6:
			m.state = stateScanningBrowsers
			return m, tea.Batch(m.scanBrowsers(), m.spinner.Tick)
		case 7:
			m.state = stateScanningDev
			return m, tea.Batch(m.scanDevTools(), m.spinner.Tick)
		case 8:
			m.state = stateScanningAppSupport
			return m, tea.Batch(m.scanAppSupport(), m.spinner.Tick)
		case 9:
			m.state = stateScanningDocker
			return m, tea.Batch(m.scanDocker(), m.spinner.Tick)
		case 10:
			_ = exec.Command("open", "x-apple.systempreferences:com.apple.settings.Storage").Run()
			m.statusMsg = "Settings opened successfully!"
		case 11:
			return m, tea.Quit
		}
	}
	return m, nil
}
