package timerview

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

type activeButton int64

const (
	startPauseButton activeButton = iota
	stopButton
)

type StopBehavior func(TimerView) (tea.Model, tea.Cmd)
type TimeoutBehavior func()

type TimerViewStyle struct {
	activeButtonStyle   lipgloss.Style
	inactiveButtonStyle lipgloss.Style
	borderStyle         lipgloss.Style
	progressBarColor    string
	startText           string
	pauseText           string
	resumeText          string
	stopText            string
	stopHelpText        string
	width               int
	height              int
	onStop              StopBehavior
	onTimeout           TimeoutBehavior
}

type TimerView struct {
	timer            timer.Model
	started          bool
	originalDuration time.Duration
	originalInterval time.Duration
	progressBar      progress.Model
	percentComplete  float64
	keymaps          []key.Binding
	help             help.Model
	activeButton     activeButton
	style            TimerViewStyle
}

func NewTimerView(duration string, interval time.Duration, style TimerViewStyle) TimerView {
	focusDuration, err := time.ParseDuration(duration)
	if err != nil {
		focusDuration = time.Minute * 25
	}

	return TimerView{
		timer: timer.NewWithInterval(focusDuration, interval),
		progressBar: progress.New(
			progress.WithSolidFill(style.progressBarColor),
			progress.WithoutPercentage(),
			progress.WithWidth(int(float64(style.width)*0.64)),
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
				key.WithHelp("s", style.stopHelpText),
				key.WithDisabled(),
			),
			key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quits the application"),
			),
		},
		help:         help.NewModel(),
		activeButton: startPauseButton,
		style:        style,
	}
}

func (m TimerView) getStartPauseButton() string {
	var buttonStyle lipgloss.Style
	if m.activeButton == startPauseButton {
		buttonStyle = m.style.activeButtonStyle
	} else {
		buttonStyle = m.style.inactiveButtonStyle
	}
	startPauseButtonText := m.getStartPauseButtonText()
	startPauseButton := buttonStyle.Render(startPauseButtonText)
	return startPauseButton
}

func (m TimerView) getStopButton() string {
	var buttonStyle lipgloss.Style
	if m.activeButton == stopButton {
		buttonStyle = m.style.activeButtonStyle
	} else {
		buttonStyle = m.style.inactiveButtonStyle
	}
	cancelButton := buttonStyle.Render(m.style.stopText)
	return cancelButton
}

func (m TimerView) View() string {
	startPauseButton := m.getStartPauseButton()
	cancelButton := m.getStopButton()
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, startPauseButton, cancelButton)

	pbar := m.progressBar.ViewAs(m.progressBar.Percent())
	timeLeft := fmt.Sprintf("\n%s\n", m.timer.View())
	help := fmt.Sprintf("\n\n%s", m.help.ShortHelpView(m.keymaps))
	ui := lipgloss.JoinVertical(lipgloss.Center, pbar, timeLeft, buttons, help)
	block := lipgloss.Place(m.style.width, m.style.height, lipgloss.Center, lipgloss.Center, m.style.borderStyle.Render(ui))
	return block
}

func (m TimerView) getStartPauseButtonText() string {
	if !m.started {
		return m.style.startText
	} else if !m.timer.Running() {
		return m.style.resumeText
	} else {
		return m.style.pauseText
	}
}

func (m TimerView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyMessage(m, msg)
	case timer.TickMsg:
		return handleTickMessage(m, msg)
	case timer.StartStopMsg:
		return handleStartStopMessage(m, msg)
	case tea.WindowSizeMsg:
		return handleResizeMessage(m, msg)
	case timer.TimeoutMsg:
		return handleTimeoutMessage(m, msg)
	}
	return m, nil
}

func (m TimerView) Init() tea.Cmd {
	return m.timer.Init()
}

func (m TimerView) PercentComplete() float64 {
	return m.percentComplete
}
func handleKeyMessage(m TimerView, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	return m, nil
}

func handleTickMessage(m TimerView, msg timer.TickMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	timeUsed := m.originalDuration.Hours() - m.timer.Timeout.Hours()
	m.percentComplete = timeUsed / m.originalDuration.Hours()
	m.progressBar.SetPercent(m.percentComplete)
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

func handleSpacebar(m TimerView) (tea.Model, tea.Cmd) {
	return startPauseTimer(m)
}

func handleEnterPressed(m TimerView) (tea.Model, tea.Cmd) {
	if m.activeButton == startPauseButton {
		return startPauseTimer(m)
	} else {
		return stopTimer(m)
	}
}

func startPauseTimer(m TimerView) (tea.Model, tea.Cmd) {
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

func stopTimer(m TimerView) (tea.Model, tea.Cmd) {
	newModel := NewTimerView(m.originalDuration.String(), m.originalInterval, m.style)
	return newModel, nil
}

func handleSPressed(m TimerView) (tea.Model, tea.Cmd) {
	return m.style.onStop(m)
}

func handleHPressed(m TimerView) (tea.Model, tea.Cmd) {
	m.activeButton = startPauseButton
	return m, nil
}

func handleLPressed(m TimerView) (tea.Model, tea.Cmd) {
	m.activeButton = stopButton
	return m, nil
}

func handleStartStopMessage(m TimerView, msg timer.StartStopMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	m.keymaps[0].SetEnabled(!m.timer.Running())
	m.keymaps[1].SetEnabled(m.timer.Running())
	m.keymaps[2].SetEnabled(true)
	return m, cmd
}

func handleResizeMessage(m TimerView, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.style.width = msg.Width
	m.style.height = msg.Height
	m.progressBar.Width = int(float64(msg.Width) * 0.64)

	return m, nil
}

func handleTimeoutMessage(m TimerView, msg timer.TimeoutMsg) (tea.Model, tea.Cmd) {
	if m.style.onTimeout != nil {
		m.style.onTimeout()
	}
	return m, focusComplete
}

func focusComplete() tea.Msg {
	return TimerCompleteMsg{}
}
