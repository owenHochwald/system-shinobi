package dojo

import "github.com/charmbracelet/lipgloss"

// Ninja theme colors (matching icon states from icon/generator.go)
var (
	colorIdle   = lipgloss.Color("#505050") // gray
	colorLow    = lipgloss.Color("#4CAF50") // green
	colorMedium = lipgloss.Color("#FFC107") // amber
	colorHigh   = lipgloss.Color("#F44336") // red
	colorDim    = lipgloss.Color("#555555")
	colorBright = lipgloss.Color("#EEEEEE")
	colorBg     = lipgloss.Color("#1A1A2E") // dark navy
	colorAccent = lipgloss.Color("#16213E") // slightly lighter navy
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBright).
			Background(colorBg).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBg).
			Background(colorLow).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(colorDim).
				Padding(0, 2)

	scrollTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorLow).
				MarginBottom(1)

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorBright).
				Underline(true)

	selectedRowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorBg).
				Background(colorLow)

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHigh).
			Background(colorAccent).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorHigh).
			Bold(true)

	infoLabelStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Width(14)

	infoValueStyle = lipgloss.NewStyle().
			Foreground(colorBright)
)

// cpuColor returns a lipgloss style colored by CPU percentage
func cpuColor(cpu float64) lipgloss.Style {
	base := lipgloss.NewStyle()
	if cpu >= 70 {
		return base.Foreground(colorHigh)
	}
	if cpu >= 40 {
		return base.Foreground(colorMedium)
	}
	if cpu >= 15 {
		return base.Foreground(colorLow)
	}
	return base.Foreground(colorDim)
}
