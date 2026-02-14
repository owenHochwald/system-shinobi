package tray

import "testing"

func TestFormatCpuLabel(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{0.0, "CPU: 0.0%"},
		{45.3, "CPU: 45.3%"},
		{100.0, "CPU: 100.0%"},
		{12.5, "CPU: 12.5%"},
	}

	for _, tt := range tests {
		result := FormatCpuLabel(tt.percent)
		if result != tt.expected {
			t.Errorf("FormatCpuLabel(%f) = %q, expected %q", tt.percent, result, tt.expected)
		}
	}
}
