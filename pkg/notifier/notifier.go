package notifier

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/types"
)

type Notify interface {
	Notify(item *types.DataItem)
	NotifyQueue() chan *types.DataItem
	SetNotifyQueue(chan *types.DataItem)
}

func NewNotify(cfg *config.NotifierConfig) Notify {
	var (
		notifier Notify
	)

	if cfg.Lark["webhookUrl"] != "" {
		notifier = NewLarkNotifier(cfg)
	} else if cfg.Dingding["webhookUrl"] != "" {
		notifier = NewDingdingNotifer(cfg)
	} else {
		return nil
	}

	// init notify queue
	notifier.SetNotifyQueue(make(chan *types.DataItem, 30))
	// message queue
	go func() {
		var item *types.DataItem
		for {
			item = <-notifier.NotifyQueue()
			notifier.Notify(item)
		}
	}()

	return notifier
}
