package timerview

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewBreakMode(duration string, interval time.Duration, width int, height int) TimerView {
	inactiveButtonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("7")).
		Padding(0, 3).
		Margin(1)

	activeButtonStyle := inactiveButtonStyle.Copy().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("2")).
		Margin(1).
		Underline(true)

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("2")).
		Padding(2, 2, 0)

	timerViewStyle := TimerViewStyle{
		inactiveButtonStyle: inactiveButtonStyle,
		activeButtonStyle:   activeButtonStyle,
		borderStyle:         border,
		progressBarColor:    "#00FF00",
		startText:           "Start",
		pauseText:           "Pause",
		resumeText:          "Resume",
		stopText:            "Skip",
		stopHelpText:        "Skips this break",
		width:               width,
		height:              height,
		onStop: func(m TimerView) (tea.Model, tea.Cmd) {
			return m, focusComplete
		},
	}

	return NewTimerView(duration, interval, timerViewStyle)
}
