package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTea(workDuration, breakDuration time.Duration, workEmoticon, breakEmoticon string, messageStatusPrefix string) {
	p := tea.NewProgram(InitWTimer(workDuration, breakDuration, workEmoticon, breakEmoticon, messageStatusPrefix), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %s", err)
		os.Exit(1)
	}
}
