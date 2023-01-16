package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTea() {
	p := tea.NewProgram(InitWTimer(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %s", err)
		os.Exit(1)
	}
}
