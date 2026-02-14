package sysinfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Info holds system information for the !clone scroll
type Info struct {
	Hostname  string
	OSVersion string
	Uptime    time.Duration
	CPUModel  string
	Cores     int
	MemTotal  uint64 // bytes
	MemUsed   uint64 // bytes
}

// Collect gathers system info from macOS APIs
func Collect() Info {
	info := Info{
		Cores: runtime.NumCPU(),
	}

	info.Hostname, _ = os.Hostname()
	info.OSVersion = getOSVersion()
	info.Uptime = getUptime()
	info.CPUModel = getCPUModel()
	info.MemTotal, info.MemUsed = getMemory()

	return info
}

// FormatMemory formats bytes as a human-readable string
func FormatMemory(bytes uint64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 1.0 {
		return fmt.Sprintf("%.1f GB", gb)
	}
	mb := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0f MB", mb)
}

// FormatUptime formats a duration as a human-readable string
func FormatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

func getOSVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "unknown"
	}
	return "macOS " + strings.TrimSpace(string(out))
}

func getUptime() time.Duration {
	tv, err := unix.SysctlTimeval("kern.boottime")
	if err != nil {
		return 0
	}
	boot := time.Unix(tv.Sec, int64(tv.Usec)*1000)
	return time.Since(boot)
}

func getCPUModel() string {
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getMemory() (total, used uint64) {
	totalMem, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0, 0
	}
	total = totalMem

	// Parse vm_stat output for memory usage
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		return total, 0
	}

	pageSize := uint64(unix.Getpagesize())
	var freePages, inactivePages uint64

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Pages free:") {
			freePages = parseVmStatValue(line)
		} else if strings.HasPrefix(line, "Pages inactive:") {
			inactivePages = parseVmStatValue(line)
		}
	}

	free := (freePages + inactivePages) * pageSize
	if free > total {
		return total, 0
	}
	used = total - free
	return
}

func parseVmStatValue(line string) uint64 {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return 0
	}
	s := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[1]), "."))
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}
