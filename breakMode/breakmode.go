package breakMode

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BreakMode struct {
	originalDuration time.Duration
	originalInterval time.Duration
}

func NewBreakMode(duration string, interval time.Duration) BreakMode {
	breakDuration, err := time.ParseDuration(duration)
	if err != nil {
		breakDuration = time.Minute * 5
	}

	return BreakMode{
		originalDuration: breakDuration,
		originalInterval: interval,
	}
}

func (m BreakMode) Init() tea.Cmd {
	return nil
}

func (m BreakMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m BreakMode) View() string {
	return ""
}
