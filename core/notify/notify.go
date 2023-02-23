package notify

import (
	"APIKiller/core/data"
)

type Notify interface {
	Notify(item *data.DataItem)
	NotifyQueue() chan *data.DataItem
	SetNotifyQueue(chan *data.DataItem)
}
