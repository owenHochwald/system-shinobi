package main

import (
	"log"

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
		cpuLabel, quitItem := tray.Setup(icons, templates)

		// Open the pipe reader
		var err error
		reader, err = pipe.NewPipeReader(pipePath)
		if err != nil {
			log.Printf("Warning: Failed to open pipe at %s: %v", pipePath, err)
			log.Printf("Make sure the probe is running. Menu will show disconnected state.")
			// Continue anyway - we'll show disconnected state
			return
		}

		// Start reading from the pipe
		reader.Start()

		// Launch goroutine to process CPU readings
		go func() {
			for reading := range reader.Readings() {
				// Classify the CPU state
				state := icon.Classify(reading.CpuPercent)

				// Update the icon and label
				tray.UpdateIcon(state)
				tray.UpdateLabel(cpuLabel, reading.CpuPercent)
			}

			// If we get here, the pipe was closed (probe disconnected)
			log.Println("Pipe closed - probe disconnected")
			tray.UpdateLabel(cpuLabel, -1) // Will show "CPU: -1.0%" as a disconnected indicator
			tray.UpdateIcon(icon.StateIdle)
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
