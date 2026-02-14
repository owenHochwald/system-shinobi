package dojo

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"system-shinobi/sensei/internal/process"
	"system-shinobi/sensei/internal/sysinfo"
)

// ScrollType identifies which scroll (tab) is active
type ScrollType int

const (
	ScrollShuriken ScrollType = iota // !shuriken - process killer
	ScrollShadow                     // !shadow  - process monitor
	ScrollClone                      // !clone   - system info
)

// Model is the top-level BubbleTea model for the Dojo TUI
type Model struct {
	currentScroll ScrollType
	width, height int

	// !shuriken state
	processes   []process.Process
	selectedIdx int
	confirmKill bool
	killResult  string

	// !shadow state
	shadowProcs []process.Process

	// !clone state
	sysInfo sysinfo.Info

	// shared
	cpuPercent float64
	err        string
}

// BubbleTea messages
type (
	processListMsg   []process.Process
	shadowRefreshMsg []process.Process
	cpuUpdateMsg     float64
	sysInfoMsg       sysinfo.Info
	killResultMsg    struct{ err error }
	tickMsg          time.Time
	errMsg           string
)

// NewModel creates a new Dojo model
func NewModel() Model {
	return Model{
		currentScroll: ScrollShuriken,
		cpuPercent:    -1,
	}
}

// Init returns the initial commands to run
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchProcesses,
		fetchSysInfo,
		fetchCPU,
		tickEvery(2*time.Second),
	)
}

// Commands that fetch data asynchronously
func fetchProcesses() tea.Msg {
	procs, err := process.ListTop(30)
	if err != nil {
		return errMsg(err.Error())
	}
	return processListMsg(procs)
}

func fetchShadow() tea.Msg {
	procs, err := process.ListTop(20)
	if err != nil {
		return errMsg(err.Error())
	}
	return shadowRefreshMsg(procs)
}

func fetchSysInfo() tea.Msg {
	return sysInfoMsg(sysinfo.Collect())
}

func fetchCPU() tea.Msg {
	cpu, err := process.GetCPUPercent()
	if err != nil {
		return cpuUpdateMsg(-1)
	}
	return cpuUpdateMsg(cpu)
}

func killProcess(pid int) tea.Cmd {
	return func() tea.Msg {
		err := process.Kill(pid)
		return killResultMsg{err: err}
	}
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
