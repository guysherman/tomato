package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guysherman/tomato/breakMode"
	"github.com/guysherman/tomato/focusMode"
)

type Tomato struct {
	currentView View
}

func (m Tomato) Init() tea.Cmd {
	return nil
}

func (m Tomato) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case focusMode.FocusCompleteMsg:
		return handleFocusComplete(m, msg)
	default:
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	}
}

func handleFocusComplete(m Tomato, msg focusMode.FocusCompleteMsg) (tea.Model, tea.Cmd) {
	m.currentView = breakMode.NewBreakMode("5m", time.Second)
	return m, nil
}

func (m Tomato) View() string {
	return m.currentView.View()
}

func main() {
	m := Tomato{
		currentView: focusMode.NewFocusMode("25m", time.Second, 120, 40),
	}

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
