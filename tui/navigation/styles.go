package navigation

import (
	"github.com/charmbracelet/lipgloss"
)

// HSL-Tailored / Dark Sleek Color Palette (Adapted Dracula Premium)
var (
	purpleTheme = lipgloss.Color("#7D56F4")
	cyanTheme   = lipgloss.Color("#00F0FF")
	greenTheme  = lipgloss.Color("#04B575")
	redTheme    = lipgloss.Color("#FF4A7A")
	yellowTheme = lipgloss.Color("#FFB86C")
	grayTheme   = lipgloss.Color("#6272A4")
	bgTheme     = lipgloss.Color("#282A36")
	fgTheme     = lipgloss.Color("#F8F8F2")
)

// Box, Border, and Text Styles for the Console GUI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(fgTheme).
			Background(purpleTheme).
			Padding(0, 2).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(cyanTheme).
			Italic(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(grayTheme).
			Italic(true)

	itemStyle = lipgloss.NewStyle().
			Foreground(fgTheme)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(cyanTheme).
				Bold(true)

	sizeStyle = lipgloss.NewStyle().
			Foreground(greenTheme).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(grayTheme).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleTheme).
			Padding(1, 2).
			Width(65)

	doneBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(greenTheme).
			Padding(1, 2).
			Width(65)

	reportBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cyanTheme).
			Padding(1, 2).
			Width(75)

	menuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleTheme).
			Padding(1, 2).
			Width(60)
)

// Visual helper to highlight text in cyan color
func cyanStyle(str string) string {
	return lipgloss.NewStyle().Foreground(cyanTheme).Render(str)
}
