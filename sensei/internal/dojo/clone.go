package dojo

import (
	"fmt"
	"strings"

	"system-shinobi/sensei/internal/sysinfo"
)

// renderClone renders the !clone system info scroll
func (m Model) renderClone() string {
	var b strings.Builder

	b.WriteString(scrollTitleStyle.Render("!CLONE - System Intelligence"))
	b.WriteString("\n\n")

	info := m.sysInfo

	rows := []struct{ label, value string }{
		{"Hostname", info.Hostname},
		{"OS", info.OSVersion},
		{"Uptime", sysinfo.FormatUptime(info.Uptime)},
		{"CPU Model", info.CPUModel},
		{"CPU Cores", fmt.Sprintf("%d", info.Cores)},
		{"Memory", fmt.Sprintf("%s / %s",
			sysinfo.FormatMemory(info.MemUsed),
			sysinfo.FormatMemory(info.MemTotal))},
		{"CPU Usage", formatCPUStatus(m.cpuPercent)},
	}

	for _, r := range rows {
		label := infoLabelStyle.Render(r.label)
		value := infoValueStyle.Render(r.value)
		b.WriteString(fmt.Sprintf("  %s %s\n", label, value))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  [r] Refresh system info"))

	return b.String()
}

func formatCPUStatus(cpu float64) string {
	if cpu < 0 {
		return "-- (probe not connected)"
	}
	return fmt.Sprintf("%.1f%%", cpu)
}
