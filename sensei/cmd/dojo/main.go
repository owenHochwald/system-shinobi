package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"system-shinobi/sensei/internal/dojo"
)

func main() {
	model := dojo.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("Error running dojo: %v", err)
		os.Exit(1)
	}
}
