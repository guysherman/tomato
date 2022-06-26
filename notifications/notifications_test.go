package notifications

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotifications(t *testing.T) {
	Convey("Notification", t, func() {
		Convey("Generates sequence with title and body", func() {
			sequences := []string{}
			n := NewNotification("Title", "Body", Focus, func(s string) { sequences = append(sequences, s) })

			n.Send()

			So(sequences[0], ShouldEqual, "\x1b]99;i=0:d=0:p=title;Title\x1b\\")
			So(sequences[1], ShouldEqual, "\x1b]99;i=0:d=1:p=body;Body\x1b\\")
		})
	})
}
