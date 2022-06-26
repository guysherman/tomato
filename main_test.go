package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guysherman/tomato/timerview"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("Main", t, func() {
		Convey("FocusCompleteMsg transitions to BreakMode", func() {
			var t tea.Model
			t = Tomato{
				longBreakTomatos: 4,
			}
			msg := timerview.TimerCompleteMsg{}
			t, cmd := t.Update(msg)

			Convey("cmd is nil", func() {
				So(cmd, ShouldBeNil)
			})
		})
	})
}
