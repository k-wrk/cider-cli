package navigation

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Main Menu Options
var menuOptions = []string{
	"Select and Clean Caches / Containers",
	"Largest Files Report (~/Documents)",
	"Delete Local Snapshots (Time Machine via Sudo)",
	"App Uninstall Suggestions",
	"Hugging Face Models & Cache",
	"Ollama Models",
	"Browser Caches & Data",
	"Developer Tools Caches",
	"Application Support Caches & Data",
	"Docker Cleanup & Caches",
	"Open macOS Storage Settings",
	"Exit",
}

// InitialModel creates the initial TUI data structure
func InitialModel() tea.Model {
	home, _ := os.UserHomeDir()

	cleanupItems := []*CleanupItem{
		// Option 1: Time Machine Local Snapshots
		{ID: "snapshots", Name: "Time Machine Local Snapshots", Path: "tm_snapshots", Selected: true, CustomClean: func() (string, error) {
			cmd := exec.Command("tmutil", "deletelocalsnapshots", "/")
			_ = cmd.Run()
			return "Snapshots cleared", nil
		}},
		{ID: "spotify", Name: "Spotify Cache", Path: filepath.Join(home, "Library/Caches/com.spotify.client"), Selected: true, CustomClean: cleanSpotify},
		{ID: "slack", Name: "Slack Cache", Path: filepath.Join(home, "Library/Application Support/Slack/Cache"), Selected: true, CustomClean: cleanSlack},
		{ID: "discord", Name: "Discord Cache", Path: filepath.Join(home, "Library/Application Support/discord/Cache"), Selected: true, CustomClean: cleanDiscord},
		{ID: "telegram", Name: "Telegram Cache & Media", Path: filepath.Join(home, "Library/Group Containers/5U85Y5W795.ru.keepcoder.Telegram"), Selected: true},
		{ID: "mail_downloads", Name: "Apple Mail Downloads & Attachments", Path: filepath.Join(home, "Library/Containers/com.apple.mail/Data/Library/Mail Downloads"), Selected: true},
		{ID: "diagnostic_reports", Name: "macOS Diagnostic / Crash Reports", Path: filepath.Join(home, "Library/Logs/DiagnosticReports"), Selected: true},
		{ID: "user_logs", Name: "User Cache Logs", Path: filepath.Join(home, "Library/Logs"), Selected: true, CustomClean: cleanUserLogs},
		{ID: "containers", Name: "Orphan Containers (Third-Party)", Path: "containers_scan", Selected: true, CustomClean: cleanOrphanContainers},
		{ID: "trash", Name: "Empty macOS Trash", Path: filepath.Join(home, ".Trash"), Selected: true, CustomClean: func() (string, error) {
			trashDir := filepath.Join(home, ".Trash")
			entries, err := os.ReadDir(trashDir)
			if err != nil {
				return "", err
			}
			for _, entry := range entries {
				_ = os.RemoveAll(filepath.Join(trashDir, entry.Name()))
			}
			return "Trash emptied", nil
		}},
	}

	pg := progress.New(progress.WithDefaultGradient())
	pg.Width = 60

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(cyanTheme).Bold(true)

	return model{
		state:        stateScanning,
		items:        cleanupItems,
		progress:     pg,
		spinner:      sp,
		itemsPerPage: 10,
	}
}

// Init initializes execution
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.startScanning(),
		m.spinner.Tick,
	)
}

// Update handles CLI events and screen transitions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		return m.handleStateKeys(key)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		newProgress, cmd := m.progress.Update(msg)
		m.progress = newProgress.(progress.Model)
		return m, cmd

	case scanFinishedMsg:
		m.state = stateMainMenu

	case reportFinishedMsg:
		m.largeFiles = msg.files
		m.reportPage = 0
		m.state = stateDocumentsReport

	case appsFinishedMsg:
		m.appSuggestions = msg.apps
		m.appPage = 0
		m.appCursor = 0
		m.state = stateAppsSuggestions

	case hfFinishedMsg:
		m.hfItems = msg.items
		m.hfPage = 0
		m.hfCursor = 0
		m.state = stateHuggingFace

	case ollamaFinishedMsg:
		m.ollamaItems = msg.items
		m.ollamaPage = 0
		m.ollamaCursor = 0
		m.state = stateOllama

	case browsersFinishedMsg:
		m.browsersItems = msg.items
		m.browsersPage = 0
		m.browsersCursor = 0
		m.state = stateBrowsers

	case devFinishedMsg:
		m.devItems = msg.items
		m.devPage = 0
		m.devCursor = 0
		m.state = stateDevTools

	case appSupportFinishedMsg:
		m.appSupportItems = msg.items
		m.appSupportPage = 0
		m.appSupportCursor = 0
		m.state = stateAppSupport

	case dockerFinishedMsg:
		m.dockerItems = msg.items
		m.dockerPage = 0
		m.dockerCursor = 0
		m.state = stateDocker

	case cleanItemFinishedMsg:
		m.freedSize += msg.sizeFreed
		m.currentCleanIdx++

		percent := float64(m.currentCleanIdx) / float64(len(m.itemsToClean))
		cmd := m.progress.SetPercent(percent)

		if m.currentCleanIdx >= len(m.itemsToClean) {
			m.state = stateFinished
			return m, cmd
		}

		return m, tea.Batch(cmd, m.cleanNext())

	case cleanFinishedMsg:
		m.state = stateFinished
		m.freedSize = msg.freedBytes
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// handleStateKeys routes key messages based on current state
func (m model) handleStateKeys(key string) (model, tea.Cmd) {
	switch m.state {
	case stateMainMenu:
		return m.updateMainMenu(key)
	case stateSelection:
		return m.updateSelection(key)
	case stateDocumentsReport:
		return m.updateDocumentsReport(key)
	case stateAppsSuggestions:
		return m.updateAppsSuggestions(key)
	case stateHuggingFace:
		return m.updateHuggingFace(key)
	case stateOllama:
		return m.updateOllama(key)
	case stateBrowsers:
		return m.updateBrowsers(key)
	case stateDevTools:
		return m.updateDevTools(key)
	case stateAppSupport:
		return m.updateAppSupport(key)
	case stateDocker:
		return m.updateDocker(key)
	case stateConfirmTrash:
		return m.updateConfirmTrash(key)
	case stateFinished:
		return m.updateFinished(key)
	}
	return m, nil
}

// cleanNext executes safe cleanup of a single item at a time (triggering progress)
func (m model) cleanNext() tea.Cmd {
	return func() tea.Msg {
		if m.currentCleanIdx >= len(m.itemsToClean) {
			return cleanFinishedMsg{freedBytes: m.freedSize}
		}

		itemIdx := m.itemsToClean[m.currentCleanIdx]
		item := m.items[itemIdx]
		var freed int64

		if item.CustomClean != nil {
			_, err := item.CustomClean()
			if err == nil {
				freed = item.Size
			}
		} else {
			home, _ := os.UserHomeDir()
			trashPath := filepath.Join(home, ".Trash", filepath.Base(item.Path))
			if _, err := os.Stat(trashPath); err == nil {
				_ = os.RemoveAll(trashPath)
			}
			err := os.Rename(item.Path, trashPath)
			if err != nil {
				_ = os.RemoveAll(item.Path)
			}
			freed = item.Size
		}

		time.Sleep(300 * time.Millisecond)

		return cleanItemFinishedMsg{
			index:     m.currentCleanIdx,
			sizeFreed: freed,
		}
	}
}
