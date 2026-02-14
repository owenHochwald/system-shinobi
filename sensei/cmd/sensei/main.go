package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/systray"
	"system-shinobi/sensei/internal/icon"
	"system-shinobi/sensei/internal/pipe"
	"system-shinobi/sensei/internal/tray"
)

const pipePath = "/tmp/shinobi.pipe"

func main() {
	// Pre-generate all icons (both colored and template versions)
	icons := icon.GenerateAll()
	templates := make(map[icon.IconState][]byte)
	templates[icon.StateIdle] = icon.GenerateTemplate(icon.StateIdle)
	templates[icon.StateLow] = icon.GenerateTemplate(icon.StateLow)
	templates[icon.StateMedium] = icon.GenerateTemplate(icon.StateMedium)
	templates[icon.StateHigh] = icon.GenerateTemplate(icon.StateHigh)

	var reader *pipe.PipeReader

	onReady := func() {
		// Setup the system tray
		cpuLabel, dojoItem, quitItem := tray.Setup(icons, templates)

		// Open the pipe reader
		var err error
		reader, err = pipe.NewPipeReader(pipePath)
		if err != nil {
			log.Printf("Warning: Failed to open pipe at %s: %v", pipePath, err)
			log.Printf("Make sure the probe is running. Menu will show disconnected state.")
			return
		}

		// Start reading from the pipe
		reader.Start()

		// Launch goroutine to process CPU readings
		go func() {
			for reading := range reader.Readings() {
				state := icon.Classify(reading.CpuPercent)
				tray.UpdateIcon(state)
				tray.UpdateLabel(cpuLabel, reading.CpuPercent, state)
			}

			// If we get here, the pipe was closed (probe disconnected)
			log.Println("Pipe closed - probe disconnected")
			tray.UpdateLabel(cpuLabel, -1, icon.StateIdle)
			tray.UpdateIcon(icon.StateIdle)
		}()

		// Handle dojo button clicks
		go func() {
			for range dojoItem.ClickedCh {
				if err := launchDojo(); err != nil {
					log.Printf("Failed to launch dojo: %v", err)
				}
			}
		}()

		// Handle quit button clicks
		go func() {
			<-quitItem.ClickedCh
			log.Println("Quit requested")
			systray.Quit()
		}()
	}

	onExit := func() {
		if reader != nil {
			reader.Stop()
		}
		log.Println("Sensei exiting...")
	}

	systray.Run(onReady, onExit)
}

// launchDojo opens a new Terminal window running the dojo binary
func launchDojo() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}
	dojoPath := filepath.Join(filepath.Dir(exePath), "dojo")

	// Verify dojo binary exists
	if _, err := os.Stat(dojoPath); err != nil {
		return fmt.Errorf("dojo binary not found at %s: %w", dojoPath, err)
	}

	// Use AppleScript to open a new Terminal window with dojo
	cmd := exec.Command("osascript",
		"-e", `tell application "Terminal"`,
		"-e", fmt.Sprintf(`do script "%s"`, dojoPath),
		"-e", `activate`,
		"-e", `end tell`,
	)
	return cmd.Start()
}
