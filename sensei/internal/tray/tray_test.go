package tray

import (
	"testing"

	"system-shinobi/sensei/internal/icon"
)

func TestFormatCpuLabel(t *testing.T) {
	tests := []struct {
		percent  float64
		state    icon.IconState
		expected string
	}{
		{0.0, icon.StateIdle, "🥷 CPU: 0.0% [Idle]"},
		{45.3, icon.StateMedium, "🥷 CPU: 45.3% [Medium]"},
		{100.0, icon.StateHigh, "🥷 CPU: 100.0% [High]"},
		{12.5, icon.StateIdle, "🥷 CPU: 12.5% [Idle]"},
		{25.0, icon.StateLow, "🥷 CPU: 25.0% [Low]"},
	}

	for _, tt := range tests {
		result := FormatCpuLabel(tt.percent, tt.state)
		if result != tt.expected {
			t.Errorf("FormatCpuLabel(%f, %v) = %q, expected %q", tt.percent, tt.state, result, tt.expected)
		}
	}
}
