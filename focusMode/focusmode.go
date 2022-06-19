package focusMode

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var inactiveButtonStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")).
	Background(lipgloss.Color("7")).
	Padding(0, 3).
	Margin(1)

var activeButtonStyle = inactiveButtonStyle.Copy().
	Foreground(lipgloss.Color("255")).
	Background(lipgloss.Color("1")).
	Margin(1).
	Underline(true)
var border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("1")).Padding(2, 2, 0)

var duration = time.Second * 25

type activeButton int64

const (
	startPauseButton activeButton = iota
	stopButton
)

type FocusMode struct {
	width            int
	height           int
	timer            timer.Model
	started          bool
	originalDuration time.Duration
	originalInterval time.Duration
	progressBar      progress.Model
	percentComplete  float64
	keymaps          []key.Binding
	help             help.Model
	activeButton     activeButton
}

func NewFocusMode(duration string, interval time.Duration, width int, height int) FocusMode {
	focusDuration, err := time.ParseDuration(duration)
	if err != nil {
		focusDuration = time.Minute * 25
	}

	return FocusMode{
		timer: timer.NewWithInterval(focusDuration, interval),
		progressBar: progress.New(
			progress.WithGradient("#FF0000", "#00FF00"),
			progress.WithoutPercentage(),
			progress.WithWidth(int(float64(width)*0.64)),
		),
		originalDuration: focusDuration,
		originalInterval: interval,
		started:          false,
		percentComplete:  0,
		keymaps: []key.Binding{
			key.NewBinding(
				key.WithKeys(tea.KeySpace.String()),
				key.WithHelp("space", "Starts the timer"),
			),
			key.NewBinding(
				key.WithKeys(tea.KeySpace.String()),
				key.WithHelp("space", "Pauses the timer"),
				key.WithDisabled(),
			),
			key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "Stops, and resets the timer"),
				key.WithDisabled(),
			),
			key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quits the application"),
			),
		},
		help:         help.NewModel(),
		width:        width,
		height:       height,
		activeButton: startPauseButton,
	}
}

func (m FocusMode) getStartPauseButton() string {
	var buttonStyle lipgloss.Style
	if m.activeButton == startPauseButton {
		buttonStyle = activeButtonStyle
	} else {
		buttonStyle = inactiveButtonStyle
	}
	startPauseButtonText := m.getStartPauseButtonText()
	startPauseButton := buttonStyle.Render(startPauseButtonText)
	return startPauseButton
}

func (m FocusMode) getStopButton() string {
	var buttonStyle lipgloss.Style
	if m.activeButton == stopButton {
		buttonStyle = activeButtonStyle
	} else {
		buttonStyle = inactiveButtonStyle
	}
	cancelButton := buttonStyle.Render("Stop")
	return cancelButton
}

func (m FocusMode) View() string {
	startPauseButton := m.getStartPauseButton()
	cancelButton := m.getStopButton()
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, startPauseButton, cancelButton)

	pbar := m.progressBar.ViewAs(m.progressBar.Percent())
	timeLeft := fmt.Sprintf("\n%s\n", m.timer.View())
	help := fmt.Sprintf("\n\n%s", m.help.ShortHelpView(m.keymaps))
	ui := lipgloss.JoinVertical(lipgloss.Center, pbar, timeLeft, buttons, help)
	block := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, border.Render(ui))
	return block
}

func (m FocusMode) getStartPauseButtonText() string {
	if !m.started {
		return "Start"
	} else if !m.timer.Running() {
		return "Resume"
	} else {
		return "Pause"
	}
}

func (m FocusMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeySpace.String():
			return handleSpacebar(m)
		case tea.KeyEnter.String():
			return handleEnterPressed(m)
		case "s":
			return handleSPressed(m)
		case "h", tea.KeyLeft.String():
			return handleHPressed(m)
		case "l", tea.KeyRight.String():
			return handleLPressed(m)
		case "q":
			return m, tea.Quit
		}
	case timer.TickMsg:
		return handleTickMessage(m, msg)
	case timer.StartStopMsg:
		return handleStartStopMessage(m, msg)
	case tea.WindowSizeMsg:
		return handleResizeMessage(m, msg)
	}
	return m, nil
}

func (m FocusMode) Init() tea.Cmd {
	return m.timer.Init()
}

func (m FocusMode) PercentComplete() float64 {
	return m.percentComplete
}

func handleTickMessage(m FocusMode, msg timer.TickMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	timeUsed := m.originalDuration.Hours() - m.timer.Timeout.Hours()
	m.percentComplete = timeUsed / m.originalDuration.Hours()
	m.progressBar.SetPercent(m.percentComplete)
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

func handleSpacebar(m FocusMode) (tea.Model, tea.Cmd) {
	return startPauseTimer(m)
}

func startPauseTimer(m FocusMode) (tea.Model, tea.Cmd) {
	if !m.started {
		m.started = true
		m.keymaps[0].SetEnabled(false)
		m.keymaps[1].SetEnabled(true)
		m.keymaps[2].SetEnabled(true)
		return m, m.timer.Init()
	} else {
		return m, m.timer.Toggle()
	}
}

func handleEnterPressed(m FocusMode) (tea.Model, tea.Cmd) {
	if m.activeButton == startPauseButton {
		return startPauseTimer(m)
	} else {
		return stopTimer(m)
	}
}

func stopTimer(m FocusMode) (tea.Model, tea.Cmd) {
	newModel := NewFocusMode(m.originalDuration.String(), m.originalInterval, m.width, m.height)
	return newModel, nil
}

func handleSPressed(m FocusMode) (tea.Model, tea.Cmd) {
	return stopTimer(m)
}

func handleHPressed(m FocusMode) (tea.Model, tea.Cmd) {
	m.activeButton = startPauseButton
	return m, nil
}

func handleLPressed(m FocusMode) (tea.Model, tea.Cmd) {
	m.activeButton = stopButton
	return m, nil
}

func handleStartStopMessage(m FocusMode, msg timer.StartStopMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	m.keymaps[0].SetEnabled(!m.timer.Running())
	m.keymaps[1].SetEnabled(m.timer.Running())
	m.keymaps[2].SetEnabled(true)
	return m, cmd
}

func handleResizeMessage(m FocusMode, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.progressBar.Width = int(float64(msg.Width) * 0.64)

	return m, nil
}
