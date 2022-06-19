package focusMode

import (
	"fmt"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFoucusMode(t *testing.T) {
	Convey("FocusMode", t, func() {
		Convey("timer is not running", func() {
			fm := NewFocusMode("1s", time.Millisecond, 120, 40)
			Convey("Pressing spacebar starts the timer", func() {
				msg := tea.KeyMsg{
					Type: tea.KeySpace,
					Alt:  false,
				}

				fm, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldResemble, fmt.Sprintf("%T", timer.TickMsg{}))

				fm, cmd = fm.Update(msg2)
				So(fm.(FocusMode).timer.Running(), ShouldBeTrue)
				So(fm.(FocusMode).keymaps[0].Enabled(), ShouldBeFalse)
				So(fm.(FocusMode).keymaps[1].Enabled(), ShouldBeTrue)
				So(fm.(FocusMode).keymaps[2].Enabled(), ShouldBeTrue)
			})

			Convey("Pressing q exits the application", func() {
				msg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'q'},
					Alt:   false,
				}

				_, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldEqual, fmt.Sprintf("%T", tea.Quit()))
			})

			Reset(func() {
				fm = NewFocusMode("1s", time.Millisecond, 120, 40)
			})
		})

		Convey("Buttons", func() {
			var fm tea.Model
			fm = NewFocusMode("1s", time.Millisecond, 120, 40)
			Convey("l switches active button to stop", func() {
				msg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'l'},
					Alt:   false,
				}

				fm, _ := fm.Update(msg)
				So(fm.(FocusMode).activeButton, ShouldEqual, stopButton)
			})

			Convey("right switches active button to stop", func() {
				msg := tea.KeyMsg{
					Type: tea.KeyRight,
					Alt:  false,
				}

				fm, _ := fm.Update(msg)
				So(fm.(FocusMode).activeButton, ShouldEqual, stopButton)
			})

			Convey("h switches active button to start", func() {
				fmm := fm.(FocusMode)
				fmm.activeButton = stopButton
				fm = fmm
				msg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'h'},
					Alt:   false,
				}

				fm, _ := fm.Update(msg)
				So(fm.(FocusMode).activeButton, ShouldEqual, startPauseButton)
			})

			Convey("left switches active button to start", func() {
				fmm := fm.(FocusMode)
				fmm.activeButton = stopButton
				fm = fmm
				msg := tea.KeyMsg{
					Type: tea.KeyLeft,
					Alt:  false,
				}

				fm, _ := fm.Update(msg)
				So(fm.(FocusMode).activeButton, ShouldEqual, startPauseButton)
			})

			Convey("When Start is active, Enter starts timer", func() {
				msg := tea.KeyMsg{
					Type: tea.KeyEnter,
					Alt:  false,
				}

				fm, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldResemble, fmt.Sprintf("%T", timer.TickMsg{}))

				fm, cmd = fm.Update(msg2)
				So(fm.(FocusMode).started, ShouldBeTrue)
				So(fm.(FocusMode).timer.Running(), ShouldBeTrue)
			})

			Convey("When Pause is active, Enter pauses timer", func() {
				fmm := fm.(FocusMode)
				fmm.started = true
				fm = fmm
				cmd := fm.(FocusMode).timer.Init()
				tickMsg := cmd()
				fm, cmd = fm.Update(tickMsg)

				msg := tea.KeyMsg{
					Type: tea.KeyEnter,
					Alt:  false,
				}

				fm, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldResemble, fmt.Sprintf("%T", timer.StartStopMsg{}))

				fm, cmd = fm.Update(msg2)
				So(fm.(FocusMode).timer.Running(), ShouldBeFalse)
			})

			Convey("When Stop is active, Enter stops and resets the timer", func() {
				fmm := fm.(FocusMode)
				fmm.started = true
				fmm.activeButton = stopButton
				fm = fmm
				cmd := fm.(FocusMode).timer.Init()
				tickMsg := cmd()
				fm, cmd = fm.Update(tickMsg)

				msg := tea.KeyMsg{
					Type: tea.KeyEnter,
					Alt:  false,
				}

				fm, cmd := fm.Update(msg)
				So(cmd, ShouldBeNil)
				So(fm.(FocusMode).started, ShouldBeFalse)
				So(fm.(FocusMode).activeButton, ShouldEqual, startPauseButton)
			})

			Reset(func() {
				fm = NewFocusMode("1s", time.Millisecond, 120, 40)
			})
		})

		Convey("the timer is running", func() {
			var fm tea.Model = NewFocusMode("1s", time.Millisecond, 120, 40)
			fmm := fm.(FocusMode)
			fmm.started = true
			fm = fmm
			cmd := fm.(FocusMode).timer.Init()
			tickMsg := cmd()
			fm, cmd = fm.Update(tickMsg)

			Convey("Pressing spacebar pauses the timer", func() {
				msg := tea.KeyMsg{
					Type: tea.KeySpace,
					Alt:  false,
				}

				fm, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldResemble, fmt.Sprintf("%T", timer.StartStopMsg{}))

				fm, cmd = fm.Update(msg2)
				So(fm.(FocusMode).timer.Running(), ShouldBeFalse)
				So(fm.(FocusMode).keymaps[0].Enabled(), ShouldBeTrue)
				So(fm.(FocusMode).keymaps[1].Enabled(), ShouldBeFalse)
				So(fm.(FocusMode).keymaps[2].Enabled(), ShouldBeTrue)
			})

			Convey("Pressing s stops, and resets, the timer", func() {
				msg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'s'},
					Alt:   false,
				}

				fm, cmd := fm.Update(msg)
				So(cmd, ShouldBeNil)
				So(fm.(FocusMode).started, ShouldBeFalse)
				So(fm.(FocusMode).keymaps[0].Enabled(), ShouldBeTrue)
				So(fm.(FocusMode).keymaps[1].Enabled(), ShouldBeFalse)
				So(fm.(FocusMode).keymaps[2].Enabled(), ShouldBeFalse)
			})

			Convey("Pressing q exits the application", func() {
				msg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'q'},
					Alt:   false,
				}

				_, cmd := fm.Update(msg)
				msg2 := cmd()
				So(fmt.Sprintf("%T", msg2), ShouldEqual, fmt.Sprintf("%T", tea.Quit()))
			})

			Convey("Tick message increases percent complete", func() {
				msg := timer.TickMsg{
					ID:      fmm.timer.ID(),
					Timeout: false,
				}

				fm, _ := fm.Update(msg)
				So(fm.(FocusMode).PercentComplete(), ShouldAlmostEqual, 0.001)
				So(fm.(FocusMode).progressBar.Percent(), ShouldAlmostEqual, 0.001)
			})

			Reset(func() {
				fm = NewFocusMode("1s", time.Millisecond, 120, 40)
				fmm = fm.(FocusMode)
				fmm.started = true
				fm = fmm
				cmd = fm.(FocusMode).timer.Init()
				tickMsg = cmd()
				fm, cmd = fm.Update(tickMsg)
			})
		})
	})
}
