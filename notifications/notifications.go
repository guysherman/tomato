package notifications

import "fmt"

type NotificationAction int

const (
	Focus NotificationAction = iota
	Report
	FocusAndReport
	None
)

type Notification struct {
	Title  string
	Body   string
	Action NotificationAction
	sender SendEscapeSequence
}

var notificationId int = 0

type SendEscapeSequence func(string)

func NewNotification(title string, body string, action NotificationAction, sender SendEscapeSequence) Notification {
	return Notification{
		Title:  title,
		Body:   body,
		Action: action,
		sender: sender,
	}
}

func (n Notification) Send() {
	for _, s := range n.kittyEscapeSequence() {
		n.sender(s)
	}

	notificationId++
}

func (n Notification) kittyEscapeSequence() []string {
	titleSequence := fmt.Sprintf("\x1b]99;i=%d:d=0:p=title;%s\x1b\\", notificationId, n.Title)
	bodySequence := fmt.Sprintf("\x1b]99;i=%d:d=1:p=body;%s\x1b\\", notificationId, n.Body)
	return []string{titleSequence, bodySequence}
}
