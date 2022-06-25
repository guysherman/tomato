package main

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guysherman/tomato/breakMode"
	"github.com/guysherman/tomato/focusMode"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("Main", t, func() {
		Convey("FocusCompleteMsg transitions to BreakMode", func() {
			var t tea.Model
			t = Tomato{}
			msg := focusMode.FocusCompleteMsg{}
			t, cmd := t.Update(msg)

			Convey("currentView becomes BreakMode", func() {
				So(fmt.Sprintf("%T", t.(Tomato).currentView), ShouldResemble, fmt.Sprintf("%T", breakMode.BreakMode{}))
			})

			Convey("cmd is nil", func() {
				So(cmd, ShouldBeNil)
			})
		})
	})
}
