package core

import (
	"APIKiller/core/data"
	"APIKiller/core/database"
	"APIKiller/core/module"
	"APIKiller/core/notify"
	"APIKiller/core/origin"
	logger "APIKiller/logger"
	"APIKiller/util"
	"context"
	"fmt"
	"sync"
	"time"
)

func NewHandler(ctx context.Context, httpItem *origin.TransferItem) {
	r := httpItem.Req

	// assembly DataItem
	item := &data.DataItem{
		Id:             util.GenerateRandomId(),
		Domain:         r.Host,
		Url:            r.URL.Path,
		Https:          r.URL.Scheme == "https",
		Method:         r.Method,
		SourceRequest:  r,
		SourceResponse: httpItem.Resp,
		VulnType:       []string{},
		VulnRequest:    nil,
		VulnResponse:   nil,
		ReportTime:     fmt.Sprintf("%v", time.Now().Unix()),
		CheckState:     false,
	}

	// enum all modules and detect
	modules := ctx.Value("modules").([]module.Detecter)
	var wg sync.WaitGroup
	for i, _ := range modules {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if modules[i] == nil {
				return
			}

			modules[i].Detect(ctx, item)
		}(i)
	}
	wg.Wait()

	// notify
	notifier := ctx.Value("notifier")
	if notifier != nil {
		notifier := ctx.Value("notifier").(notify.Notify)
		notifier.NotifyQueue() <- item
	}

	// print result and save result
	logger.Infoln(fmt.Sprintf("%v %v checkout: %v", item.Domain, item.Url, item.VulnType))
	db := ctx.Value("db").(database.Database)
	db.ItemAddQueue() <- item

}
