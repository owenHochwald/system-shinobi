package dojo

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"system-shinobi/sensei/internal/process"
	"system-shinobi/sensei/internal/sysinfo"
)

// Update handles all BubbleTea messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case processListMsg:
		m.processes = []process.Process(msg)
		return m, nil

	case shadowRefreshMsg:
		m.shadowProcs = []process.Process(msg)
		return m, nil

	case cpuUpdateMsg:
		m.cpuPercent = float64(msg)
		return m, nil

	case sysInfoMsg:
		m.sysInfo = sysinfo.Info(msg)
		return m, nil

	case killResultMsg:
		if msg.err != nil {
			m.killResult = errorStyle.Render(fmt.Sprintf("  Kill failed: %v", msg.err))
		} else {
			m.killResult = scrollTitleStyle.Render("  Target eliminated.")
		}
		m.confirmKill = false
		// Refresh process list after kill
		return m, fetchProcesses

	case tickMsg:
		// Periodic refresh: update shadow processes and CPU
		cmds := []tea.Cmd{
			fetchCPU,
			tickEvery(2 * time.Second),
		}
		if m.currentScroll == ScrollShadow {
			cmds = append(cmds, fetchShadow)
		}
		return m, tea.Batch(cmds...)

	case errMsg:
		m.err = string(msg)
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		if !m.confirmKill {
			return m, tea.Quit
		}
	case "tab":
		m.confirmKill = false
		m.killResult = ""
		m.currentScroll = (m.currentScroll + 1) % 3
		return m, m.scrollEnterCmd()
	case "shift+tab":
		m.confirmKill = false
		m.killResult = ""
		m.currentScroll = (m.currentScroll + 2) % 3 // wraps backward
		return m, m.scrollEnterCmd()
	case "1":
		m.currentScroll = ScrollShuriken
		m.confirmKill = false
		m.killResult = ""
		return m, fetchProcesses
	case "2":
		m.currentScroll = ScrollShadow
		return m, fetchShadow
	case "3":
		m.currentScroll = ScrollClone
		return m, tea.Batch(fetchSysInfo, fetchCPU)
	case "r":
		return m, m.scrollEnterCmd()
	}

	// Scroll-specific keys
	switch m.currentScroll {
	case ScrollShuriken:
		return m.handleShurikenKey(msg)
	}

	return m, nil
}

func (m Model) handleShurikenKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIdx > 0 {
			m.selectedIdx--
			m.confirmKill = false
			m.killResult = ""
		}
	case "down", "j":
		max := m.visibleProcessCount() - 1
		if m.selectedIdx < max {
			m.selectedIdx++
			m.confirmKill = false
			m.killResult = ""
		}
	case "enter":
		if m.confirmKill && m.selectedIdx < len(m.processes) {
			return m, killProcess(m.processes[m.selectedIdx].PID)
		}
		m.confirmKill = true
	case "esc":
		m.confirmKill = false
	}
	return m, nil
}

// scrollEnterCmd returns the command to run when entering a scroll
func (m Model) scrollEnterCmd() tea.Cmd {
	switch m.currentScroll {
	case ScrollShuriken:
		return fetchProcesses
	case ScrollShadow:
		return fetchShadow
	case ScrollClone:
		return tea.Batch(fetchSysInfo, fetchCPU)
	}
	return nil
}
