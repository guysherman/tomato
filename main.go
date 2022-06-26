package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guysherman/tomato/timerview"
)

type timerMode int

const (
	focus timerMode = iota
	shortBreak
	longBreak
)

type Tomato struct {
	currentView      View
	mode             timerMode
	tomatoCount      int
	currentWidth     int
	currentHeight    int
	focusTime        string
	shortBreakTime   string
	longBreakTime    string
	longBreakTomatos int
}

func (m Tomato) Init() tea.Cmd {
	return nil
}

func (m Tomato) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timerview.TimerCompleteMsg:
		return handleTimerComplete(m, msg)
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.currentWidth = msg.Width
		m.currentHeight = msg.Height
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	default:
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	}
}

func handleTimerComplete(m Tomato, msg timerview.TimerCompleteMsg) (tea.Model, tea.Cmd) {
	if m.mode == focus {
		m.tomatoCount++
		if m.tomatoCount%m.longBreakTomatos == 0 {
			m.mode = longBreak
		} else {
			m.mode = shortBreak
		}
	} else {
		m.mode = focus
	}

	m.currentView = m.viewForMode()

	return m, nil
}

func (m Tomato) viewForMode() View {
	if m.mode == focus {
		return timerview.NewFocusMode(m.focusTime, time.Second, m.currentWidth, m.currentHeight)
	} else if m.mode == shortBreak {
		return timerview.NewBreakMode(m.shortBreakTime, time.Second, m.currentWidth, m.currentHeight)
	} else {
		return timerview.NewBreakMode(m.longBreakTime, time.Second, m.currentWidth, m.currentHeight)
	}
}

func (m Tomato) View() string {
	return m.currentView.View()
}

func main() {
	var focusTimeFlag = flag.String("f", "25m", "Sets the length of the focus period, expressed in <number><unit> eg 25m")
	var shortBreakTimeFlag = flag.String("s", "5m", "Sets the length of the short break, expressed in <number><unit> eg 5m")
	var longBreakTimeFlag = flag.String("l", "15m", "Sets the length of the long break, expressed in <number><unit> eg 15m")
	var longBreakTomatosFlag = flag.Int("L", 4, "Sets the number of tomatos per long break, expressed in <number> eg 4")

	flag.Parse()
	m := Tomato{
		currentView:      timerview.NewFocusMode(*focusTimeFlag, time.Second, 120, 40),
		mode:             focus,
		tomatoCount:      0,
		currentWidth:     120,
		currentHeight:    40,
		focusTime:        *focusTimeFlag,
		shortBreakTime:   *shortBreakTimeFlag,
		longBreakTime:    *longBreakTimeFlag,
		longBreakTomatos: *longBreakTomatosFlag,
	}

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
