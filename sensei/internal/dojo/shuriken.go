package dojo

import (
	"fmt"
	"strings"
)

// renderShuriken renders the !shuriken process killer scroll
func (m Model) renderShuriken() string {
	var b strings.Builder

	b.WriteString(scrollTitleStyle.Render("!SHURIKEN - Process Assassin"))
	b.WriteString("\n\n")

	if m.killResult != "" {
		b.WriteString(m.killResult)
		b.WriteString("\n\n")
	}

	if len(m.processes) == 0 {
		b.WriteString("  Loading processes...")
		return b.String()
	}

	// Table header
	header := fmt.Sprintf("  %-7s %-7s %-7s %s", "PID", "CPU%", "MEM%", "Name")
	b.WriteString(tableHeaderStyle.Render(header))
	b.WriteString("\n")

	// Process rows
	visible := m.visibleProcessCount()
	for i, p := range m.processes {
		if i >= visible {
			break
		}

		row := fmt.Sprintf("  %-7d %-7.1f %-7.1f %s", p.PID, p.CPU, p.Memory, truncate(p.Name, 30))

		if i == m.selectedIdx {
			if m.confirmKill {
				row = confirmStyle.Render(fmt.Sprintf(" KILL PID %d (%s)? [Enter] Yes  [Esc] No ", p.PID, truncate(p.Name, 15)))
			} else {
				row = selectedRowStyle.Render(row)
			}
		} else {
			row = cpuColor(p.CPU).Render(row)
		}

		b.WriteString(row)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  [up/down] Navigate  [Enter] Kill  [r] Refresh"))

	return b.String()
}

func (m Model) visibleProcessCount() int {
	// Reserve lines for header, title, help, status bar
	available := m.height - 10
	if available < 5 {
		available = 5
	}
	if available > len(m.processes) {
		available = len(m.processes)
	}
	return available
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "~"
}
