package process

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// ListTop returns the top n processes sorted by CPU usage
func ListTop(n int) ([]Process, error) {
	out, err := exec.Command("ps", "-Aceo", "pid,pcpu,pmem,comm").Output()
	if err != nil {
		return nil, fmt.Errorf("ps command failed: %w", err)
	}
	procs := parsePsOutput(string(out))
	SortByCPU(procs)
	if len(procs) > n {
		procs = procs[:n]
	}
	return procs, nil
}

// Kill sends SIGTERM to the process with the given PID
func Kill(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}

// GetCPUPercent returns total CPU usage by summing all process CPU percentages
// and dividing by the number of logical cores (ps reports per-core percentages)
func GetCPUPercent() (float64, error) {
	out, err := exec.Command("ps", "-Aceo", "pcpu").Output()
	if err != nil {
		return 0, fmt.Errorf("ps command failed: %w", err)
	}
	var total float64
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n")[1:] {
		val, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		if err == nil {
			total += val
		}
	}
	return total, nil
}

// parsePsOutput parses the output of ps -Aceo pid,pcpu,pmem,comm
func parsePsOutput(output string) []Process {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil
	}

	var procs []Process
	for _, line := range lines[1:] { // skip header
		p, ok := parseLine(line)
		if ok {
			procs = append(procs, p)
		}
	}
	return procs
}

// parseLine parses a single line of ps output
func parseLine(line string) (Process, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return Process{}, false
	}

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return Process{}, false
	}
	cpu, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Process{}, false
	}
	mem, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Process{}, false
	}
	// comm can contain spaces, so join remaining fields
	name := strings.Join(fields[3:], " ")

	return Process{PID: pid, Name: name, CPU: cpu, Memory: mem}, true
}
