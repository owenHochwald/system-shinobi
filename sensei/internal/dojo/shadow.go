package dojo

import (
	"fmt"
	"strings"
)

// renderShadow renders the !shadow process monitor scroll
func (m Model) renderShadow() string {
	var b strings.Builder

	b.WriteString(scrollTitleStyle.Render("!SHADOW - Process Monitor"))
	b.WriteString("\n\n")

	if len(m.shadowProcs) == 0 {
		b.WriteString("  Loading processes...")
		return b.String()
	}

	// Table header
	header := fmt.Sprintf("  %-7s %-7s %-7s %s", "PID", "CPU%", "MEM%", "Name")
	b.WriteString(tableHeaderStyle.Render(header))
	b.WriteString("\n")

	// Process rows (read-only, no selection)
	visible := m.height - 10
	if visible < 5 {
		visible = 5
	}
	for i, p := range m.shadowProcs {
		if i >= visible {
			break
		}
		row := fmt.Sprintf("  %-7d %-7.1f %-7.1f %s", p.PID, p.CPU, p.Memory, truncate(p.Name, 30))
		b.WriteString(cpuColor(p.CPU).Render(row))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  Auto-refreshes every 2s  [r] Force refresh"))

	return b.String()
}
