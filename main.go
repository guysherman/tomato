package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guysherman/tomato/focusMode"
)

type model struct {
	currentView View
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m.currentView.Update(msg)
}

func (m model) View() string {
	return m.currentView.View()
}

func main() {
	m := model{
		currentView: focusMode.NewFocusMode("25m", time.Second, 120, 40),
	}

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
