package notify

import (
	"APIKiller/core/data"
)

type Notify interface {
	Notify(item *data.DataItem)
	GetQueue() chan *data.DataItem
}
