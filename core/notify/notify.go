package notify

import (
	"APIKiller/core/data"
)

type Notify interface {
	Notify(item *data.DataItem)
}

var (
	NotificationQueue chan *data.DataItem
	notifier          Notify
)

//
// CreateNotification
//  @Description: create new notification and throw it into queue
//  @param notification
//
func CreateNotification(notification *data.DataItem) {
	if NotificationQueue != nil {
		NotificationQueue <- notification
	}
}

func BindNotifier(n Notify) {
	notifier = n

	// init notificationQueue
	NotificationQueue = make(chan *data.DataItem, 1024)

	// notification queue handle
	go func() {
		var item *data.DataItem
		for {
			item = <-NotificationQueue
			notifier.Notify(item)
		}
	}()
}
