package pipe

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

// CpuReading represents a single CPU measurement from the probe
type CpuReading struct {
	CpuPercent float64 `json:"cpu_percent"`
	Timestamp  int64   `json:"timestamp"`
}

// PipeReader reads CPU readings from a named pipe (FIFO)
type PipeReader struct {
	reader   io.ReadCloser
	readings chan CpuReading
	done     chan struct{}
}

// NewPipeReader creates a PipeReader that opens the FIFO at the given path
func NewPipeReader(path string) (*PipeReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &PipeReader{
		reader:   file,
		readings: make(chan CpuReading, 10),
		done:     make(chan struct{}),
	}, nil
}

// NewPipeReaderFromReader creates a PipeReader from an io.Reader (for testing)
func NewPipeReaderFromReader(r io.Reader) *PipeReader {
	return &PipeReader{
		reader:   io.NopCloser(r),
		readings: make(chan CpuReading, 10),
		done:     make(chan struct{}),
	}
}

// Readings returns the channel of CPU readings
func (pr *PipeReader) Readings() <-chan CpuReading {
	return pr.readings
}

// Start begins reading from the pipe in a goroutine
func (pr *PipeReader) Start() {
	go pr.readLoop()
}

// Stop closes the reader and drains the channel
func (pr *PipeReader) Stop() {
	close(pr.done)
	pr.reader.Close()
	// Drain any remaining readings
	for range pr.readings {
	}
}

func (pr *PipeReader) readLoop() {
	defer close(pr.readings)

	scanner := bufio.NewScanner(pr.reader)
	for scanner.Scan() {
		select {
		case <-pr.done:
			return
		default:
		}

		line := scanner.Text()
		var reading CpuReading
		if err := json.Unmarshal([]byte(line), &reading); err != nil {
			// Skip malformed lines with a warning
			log.Printf("Skipping malformed JSON line: %s (error: %v)", line, err)
			continue
		}

		select {
		case pr.readings <- reading:
		case <-pr.done:
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from pipe: %v", err)
	}
}
