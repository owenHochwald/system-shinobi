package tray

import (
	"fmt"

	"fyne.io/systray"
	"system-shinobi/sensei/internal/icon"
)

var currentIcons map[icon.IconState][]byte
var templateIcons map[icon.IconState][]byte

// Setup initializes the system tray with menu items and returns references to them
func Setup(icons map[icon.IconState][]byte, templates map[icon.IconState][]byte) (*systray.MenuItem, *systray.MenuItem) {
	currentIcons = icons
	templateIcons = templates

	// Set initial icon to idle state
	UpdateIcon(icon.StateIdle)

	// Set tooltip
	systray.SetTooltip("System Shinobi - CPU Monitor")

	// Create menu items
	cpuLabel := systray.AddMenuItem("CPU: --%", "Current CPU usage")
	cpuLabel.Disable() // Make it read-only

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit", "Exit System Shinobi")

	return cpuLabel, quit
}

// UpdateIcon swaps the systray icon based on the current state
func UpdateIcon(state icon.IconState) {
	if templateIcons != nil && currentIcons != nil {
		// Use template icon for macOS light/dark mode support
		systray.SetTemplateIcon(templateIcons[state], currentIcons[state])
	} else if currentIcons != nil {
		// Fallback to regular icon
		systray.SetIcon(currentIcons[state])
	}
}

// UpdateLabel updates the CPU percentage display in the menu
func UpdateLabel(cpuLabel *systray.MenuItem, percent float64) {
	cpuLabel.SetTitle(FormatCpuLabel(percent))
}

// FormatCpuLabel formats the CPU percentage for display
func FormatCpuLabel(percent float64) string {
	return fmt.Sprintf("CPU: %.1f%%", percent)
}
