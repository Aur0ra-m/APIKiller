package notify

import (
	"APIKiller/internal/core/data"
)

type Notify interface {
	//
	// Notify
	//  @Description:
	//  @param item
	//
	Notify(item *data.DataItem)
}

var (
	notificationQueue chan *data.DataItem
	notifier          Notify
)

//
// CreateNotification
//  @Description: create new notification and throw it into queue
//  @param notification
//
func CreateNotification(notification *data.DataItem) {
	if notificationQueue != nil {
		notificationQueue <- notification
	}
}

func BindNotifier(n Notify) {
	notifier = n

	// init notificationQueue
	notificationQueue = make(chan *data.DataItem, 1024)

	// notification queue handle
	go func() {
		var item *data.DataItem
		for {
			item = <-notificationQueue
			notifier.Notify(item)
		}
	}()
}
