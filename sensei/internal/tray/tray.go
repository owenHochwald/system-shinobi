package tray

import (
	"fmt"

	"fyne.io/systray"
	"system-shinobi/sensei/internal/icon"
)

var currentIcons map[icon.IconState][]byte
var templateIcons map[icon.IconState][]byte

// Setup initializes the system tray with menu items and returns references to them
func Setup(icons map[icon.IconState][]byte, templates map[icon.IconState][]byte) (*systray.MenuItem, *systray.MenuItem, *systray.MenuItem) {
	currentIcons = icons
	templateIcons = templates

	// Set initial icon to idle state
	UpdateIcon(icon.StateIdle)

	// Set tooltip
	systray.SetTooltip("System Shinobi - CPU Monitor")

	// Create menu items
	cpuLabel := systray.AddMenuItem("🥷 CPU: --% [Idle]", "Current CPU usage and ninja state")
	cpuLabel.Disable() // Make it read-only

	systray.AddSeparator()

	dojoItem := systray.AddMenuItem("Open Dojo (Terminal UI)", "Launch the Dojo process manager")

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit Shinobi", "Exit System Shinobi")

	return cpuLabel, dojoItem, quit
}

// UpdateIcon swaps the systray icon based on the current state
func UpdateIcon(state icon.IconState) {
	if templateIcons != nil && currentIcons != nil {
		systray.SetTemplateIcon(templateIcons[state], currentIcons[state])
	} else if currentIcons != nil {
		systray.SetIcon(currentIcons[state])
	}
}

// UpdateLabel updates the CPU percentage and state display in the menu
func UpdateLabel(cpuLabel *systray.MenuItem, percent float64, state icon.IconState) {
	cpuLabel.SetTitle(FormatCpuLabel(percent, state))
}

// FormatCpuLabel formats the CPU percentage and state for display
func FormatCpuLabel(percent float64, state icon.IconState) string {
	stateNames := map[icon.IconState]string{
		icon.StateIdle:   "Idle",
		icon.StateLow:    "Low",
		icon.StateMedium: "Medium",
		icon.StateHigh:   "High",
	}
	name := stateNames[state]
	return fmt.Sprintf("🥷 CPU: %.1f%% [%s]", percent, name)
}
