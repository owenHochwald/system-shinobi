package process

import "sort"

// Process represents a running system process
type Process struct {
	PID    int
	Name   string
	CPU    float64
	Memory float64
}

// SortByCPU sorts processes by CPU usage descending
func SortByCPU(procs []Process) {
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].CPU > procs[j].CPU
	})
}
