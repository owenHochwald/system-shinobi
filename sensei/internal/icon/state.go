package icon

// IconState represents different CPU load states
type IconState int

const (
	StateIdle   IconState = iota // 0-15%
	StateLow                     // 15-40%
	StateMedium                  // 40-70%
	StateHigh                    // 70-100%
)

// Classify determines the IconState based on CPU percentage
func Classify(cpuPercent float64) IconState {
	if cpuPercent < 15.0 {
		return StateIdle
	}
	if cpuPercent < 40.0 {
		return StateLow
	}
	if cpuPercent < 70.0 {
		return StateMedium
	}
	return StateHigh
}

// String returns the string representation of an IconState
func (s IconState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateLow:
		return "low"
	case StateMedium:
		return "medium"
	case StateHigh:
		return "high"
	default:
		return "unknown"
	}
}
