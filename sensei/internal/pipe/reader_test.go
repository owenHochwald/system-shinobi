package pipe

import (
	"strings"
	"testing"
	"time"
)

func TestParseSingleReading(t *testing.T) {
	input := `{"cpu_percent":45.3,"timestamp":1707860342}` + "\n"
	reader := NewPipeReaderFromReader(strings.NewReader(input))
	reader.Start()
	defer reader.Stop()

	select {
	case reading := <-reader.Readings():
		if reading.CpuPercent != 45.3 {
			t.Errorf("Expected cpu_percent 45.3, got %f", reading.CpuPercent)
		}
		if reading.Timestamp != 1707860342 {
			t.Errorf("Expected timestamp 1707860342, got %d", reading.Timestamp)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for reading")
	}
}

func TestParseMultipleReadings(t *testing.T) {
	input := `{"cpu_percent":10.0,"timestamp":1707860340}` + "\n" +
		`{"cpu_percent":45.3,"timestamp":1707860341}` + "\n" +
		`{"cpu_percent":89.7,"timestamp":1707860342}` + "\n"

	reader := NewPipeReaderFromReader(strings.NewReader(input))
	reader.Start()
	defer reader.Stop()

	expected := []struct {
		cpu       float64
		timestamp int64
	}{
		{10.0, 1707860340},
		{45.3, 1707860341},
		{89.7, 1707860342},
	}

	for i, exp := range expected {
		select {
		case reading := <-reader.Readings():
			if reading.CpuPercent != exp.cpu {
				t.Errorf("Reading %d: expected cpu_percent %f, got %f", i, exp.cpu, reading.CpuPercent)
			}
			if reading.Timestamp != exp.timestamp {
				t.Errorf("Reading %d: expected timestamp %d, got %d", i, exp.timestamp, reading.Timestamp)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timeout waiting for reading %d", i)
		}
	}
}

func TestMalformedLineSkipped(t *testing.T) {
	input := `{"invalid json` + "\n" +
		`{"cpu_percent":45.3,"timestamp":1707860342}` + "\n"

	reader := NewPipeReaderFromReader(strings.NewReader(input))
	reader.Start()
	defer reader.Stop()

	// Should skip the malformed line and get the valid one
	select {
	case reading := <-reader.Readings():
		if reading.CpuPercent != 45.3 {
			t.Errorf("Expected cpu_percent 45.3, got %f", reading.CpuPercent)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for reading")
	}
}

func TestZeroCpu(t *testing.T) {
	input := `{"cpu_percent":0.0,"timestamp":1707860342}` + "\n"
	reader := NewPipeReaderFromReader(strings.NewReader(input))
	reader.Start()
	defer reader.Stop()

	select {
	case reading := <-reader.Readings():
		if reading.CpuPercent != 0.0 {
			t.Errorf("Expected cpu_percent 0.0, got %f", reading.CpuPercent)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for reading")
	}
}

func TestHundredCpu(t *testing.T) {
	input := `{"cpu_percent":100.0,"timestamp":1707860342}` + "\n"
	reader := NewPipeReaderFromReader(strings.NewReader(input))
	reader.Start()
	defer reader.Stop()

	select {
	case reading := <-reader.Readings():
		if reading.CpuPercent != 100.0 {
			t.Errorf("Expected cpu_percent 100.0, got %f", reading.CpuPercent)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for reading")
	}
}
