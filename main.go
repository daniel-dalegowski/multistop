package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type stopwatch struct {
	name      string
	elapsed   time.Duration
	running   bool
	startedAt time.Time
}

func (s *stopwatch) current() time.Duration {
	if s.running {
		return s.elapsed + time.Since(s.startedAt)
	}
	return s.elapsed
}

func (s *stopwatch) toggle() {
	if s.running {
		s.elapsed += time.Since(s.startedAt)
		s.running = false
	} else {
		s.startedAt = time.Now()
		s.running = true
	}
}

func (s *stopwatch) reset() {
	s.running = false
	s.elapsed = 0
}

type mode int

const (
	modeNormal mode = iota
	modeAdding
	modeRenaming
)

type model struct {
	watches  []*stopwatch
	cursor   int
	mode     mode
	input    textinput.Model
	quitting bool
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	ti := textinput.New()
	ti.CharLimit = 40
	ti.Width = 30

	return model{
		watches: []*stopwatch{
			{name: "Stopwatch 1"},
		},
		cursor: 0,
		mode:   modeNormal,
		input:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tickCmd()

	case tea.KeyMsg:
		if m.mode == modeAdding || m.mode == modeRenaming {
			switch msg.String() {
			case "enter":
				name := m.input.Value()
				if name == "" {
					name = "Stopwatch"
				}
				if m.mode == modeAdding {
					m.watches = append(m.watches, &stopwatch{name: name})
					m.cursor = len(m.watches) - 1
				} else {
					m.watches[m.cursor].name = name
				}
				m.mode = modeNormal
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			case "esc":
				m.mode = modeNormal
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.watches)-1 {
				m.cursor++
			}

		case " ":
			if len(m.watches) > 0 {
				m.watches[m.cursor].toggle()
			}

		case "r":
			if len(m.watches) > 0 {
				m.watches[m.cursor].reset()
			}

		case "n":
			m.mode = modeAdding
			m.input.SetValue("")
			m.input.Placeholder = fmt.Sprintf("Stopwatch %d", len(m.watches)+1)
			m.input.Focus()
			return m, textinput.Blink

		case "e":
			if len(m.watches) > 0 {
				m.mode = modeRenaming
				m.input.SetValue(m.watches[m.cursor].name)
				m.input.Focus()
				return m, textinput.Blink
			}

		case "d":
			if len(m.watches) > 0 {
				m.watches = append(m.watches[:m.cursor], m.watches[m.cursor+1:]...)
				if m.cursor >= len(m.watches) && m.cursor > 0 {
					m.cursor--
				}
			}
		}
	}
	return m, nil
}

var (
	greenDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("●")
	redDot     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("●")
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	nameStyle   = lipgloss.NewStyle().Bold(true).Width(20)
	timeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginTop(1)
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).MarginBottom(1)
)

func formatDuration(d time.Duration) string {
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	s := int((d % time.Minute) / time.Second)
	cs := int((d % time.Second) / (10 * time.Millisecond))
	return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, cs)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := titleStyle.Render("Multi-Stopwatch") + "\n"

	for i, w := range m.watches {
		dot := redDot
		if w.running {
			dot = greenDot
		}

		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		name := w.name
		if m.mode == modeRenaming && i == m.cursor {
			name = m.input.View()
		} else {
			name = nameStyle.Render(name)
		}

		line := fmt.Sprintf("%s%s %s  %s", cursor, dot, name, timeStyle.Render(formatDuration(w.current())))
		s += line + "\n"
	}

	if m.mode == modeAdding {
		s += "\n" + dimStyle.Render("new stopwatch name: ") + m.input.View()
	}

	s += helpStyle.Render("\nspace start/pause  r reset  n new  e rename  d delete  ↑/↓ select  q quit")

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
