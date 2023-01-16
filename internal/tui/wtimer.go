package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding   = 2
	maxWidth  = 65
	workTime  = time.Second * 20
	breakTime = time.Second * 5
)

type tickMsg time.Time

type state int

const (
	workState = iota
	breakState
)

type model struct {
	state     state
	progress  progress.Model
	stopwatch stopwatch.Model
	keymap    keymap
	help      help.Model
	quitting  bool
}

type keymap struct {
	start,
	stop,
	reset,
	quit key.Binding
}

func InitWTimer() tea.Model {
	m := model{
		state:     workState,
		progress:  progress.New(progress.WithDefaultGradient()),
		stopwatch: stopwatch.NewWithInterval(time.Second),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	m.keymap.start.SetEnabled(true)

	return &m
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), m.stopwatch.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			return m, m.stopwatch.Reset()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case tickMsg:
		cmds := []tea.Cmd{}
		if m.progress.Percent() == 1.0 {
			m.progress.SetPercent(0)
			if m.state == workState {
				m.state = breakState
			} else {
				m.state = workState
			}
			cmds = append(cmds, m.stopwatch.Reset())
		}

		cmds = append(cmds, m.progress.IncrPercent(calcPercent(m.stopwatch.Elapsed(), workTime)))
		cmds = append(cmds, tickCmd())
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)

	var b strings.Builder
	b.WriteString(pad)
	b.WriteString("\n")
	b.WriteString(pad)
	b.WriteString("We are currently")
	if m.state == workState {
		b.WriteString(" working...")
	}
	if m.state == breakState {
		b.WriteString(" on a break...")
	}
	b.WriteString("\n\n")
	b.WriteString(pad)
	b.WriteString(m.progress.View())

	s := m.stopwatch.View()
	if !m.quitting {
		b.WriteString("\n\n")
		b.WriteString(pad)
		b.WriteString("Elapsed: ")
		b.WriteString(s)
		b.WriteString("\n")
		b.WriteString(m.helpView())
	}

	return b.String()
}

func (m model) helpView() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	}))

	return b.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func calcPercent(cd, td time.Duration) float64 {
	return 0.25
}
