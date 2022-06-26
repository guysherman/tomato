package timerview

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewFocusMode(duration string, interval time.Duration, width int, height int) TimerView {
	inactiveButtonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("7")).
		Padding(0, 3).
		Margin(1)

	activeButtonStyle := inactiveButtonStyle.Copy().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("1")).
		Margin(1).
		Underline(true)

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(2, 2, 0)

	timerViewStyle := TimerViewStyle{
		inactiveButtonStyle: inactiveButtonStyle,
		activeButtonStyle:   activeButtonStyle,
		borderStyle:         border,
		progressBarColor:    "#FF0000",
		startText:           "Start",
		pauseText:           "Pause",
		resumeText:          "Resume",
		stopText:            "Stop",
		stopHelpText:        "Stops, and resets, the timer",
		width:               width,
		height:              height,
		onStop: func(m TimerView) (tea.Model, tea.Cmd) {
			return stopTimer(m)
		},
	}

	return NewTimerView(duration, interval, timerViewStyle)
}
