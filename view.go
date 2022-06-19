package main

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	Init() tea.Cmd
	View() string
	Update(m tea.Msg) (tea.Model, tea.Cmd)
}
