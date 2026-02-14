package icon

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		cpuPercent float64
		expected   IconState
	}{
		{0.0, StateIdle},
		{14.9, StateIdle},
		{15.0, StateLow},
		{39.9, StateLow},
		{40.0, StateMedium},
		{69.9, StateMedium},
		{70.0, StateHigh},
		{100.0, StateHigh},
	}

	for _, tt := range tests {
		result := Classify(tt.cpuPercent)
		if result != tt.expected {
			t.Errorf("Classify(%f) = %v, expected %v", tt.cpuPercent, result, tt.expected)
		}
	}
}

func TestIconStateString(t *testing.T) {
	tests := []struct {
		state    IconState
		expected string
	}{
		{StateIdle, "idle"},
		{StateLow, "low"},
		{StateMedium, "medium"},
		{StateHigh, "high"},
	}

	for _, tt := range tests {
		result := tt.state.String()
		if result != tt.expected {
			t.Errorf("%v.String() = %q, expected %q", tt.state, result, tt.expected)
		}
	}
}
