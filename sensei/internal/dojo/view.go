package dojo

import (
	"fmt"
	"strings"
)

const ninjaHeader = `
  ╔═══════════════════════════════════════╗
  ║   🥷  SYSTEM SHINOBI DOJO  🥷        ║
  ╚═══════════════════════════════════════╝`

// View renders the full TUI
func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(headerStyle.Render(ninjaHeader))
	b.WriteString("\n\n")

	// Tab bar
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Active scroll content
	switch m.currentScroll {
	case ScrollShuriken:
		b.WriteString(m.renderShuriken())
	case ScrollShadow:
		b.WriteString(m.renderShadow())
	case ScrollClone:
		b.WriteString(m.renderClone())
	}

	// Error display
	if m.err != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("  Error: " + m.err))
	}

	// Fill remaining space
	content := b.String()
	lines := strings.Count(content, "\n")
	for i := lines; i < m.height-2; i++ {
		b.WriteString("\n")
	}

	// Status bar at bottom
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m Model) renderTabs() string {
	tabs := []struct {
		name   string
		scroll ScrollType
	}{
		{"!shuriken", ScrollShuriken},
		{"!shadow", ScrollShadow},
		{"!clone", ScrollClone},
	}

	var parts []string
	for _, t := range tabs {
		if t.scroll == m.currentScroll {
			parts = append(parts, activeTabStyle.Render(t.name))
		} else {
			parts = append(parts, inactiveTabStyle.Render(t.name))
		}
	}

	return "  " + strings.Join(parts, " ")
}

func (m Model) renderStatusBar() string {
	cpuStr := "--"
	if m.cpuPercent >= 0 {
		cpuStr = fmt.Sprintf("%.1f%%", m.cpuPercent)
	}
	status := fmt.Sprintf(" CPU: %s  |  [Tab] Switch  [q] Quit  [1/2/3] Jump ", cpuStr)
	return statusBarStyle.Render(status)
}
