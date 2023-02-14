package core

import (
	"APIKiller/core/data"
	"APIKiller/core/database"
	"APIKiller/core/module"
	"APIKiller/core/notify"
	"APIKiller/core/origin"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"fmt"
	"time"
)

func NewHandler(ctx context.Context, httpItem *origin.HttpItem) {
	r := httpItem.Req

	// assembly DataItem
	item := &data.DataItem{
		Id:             util.GenerateRandomId(),
		Domain:         r.Host,
		Url:            r.URL.Path,
		Https:          httpItem.Https,
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
	detectors := ctx.Value("modules").([]module.Detecter)
	for _, detector := range detectors {
		if detector == nil {
			continue
		}

		detector.Detect(ctx, item)
	}

	// notify
	notifier := ctx.Value("notifier").(notify.Notify)
	notifier.Notify(ctx, item)

	// print result and save result
	logger.Infoln(fmt.Sprintf("%v %v checkout: %v", item.Domain, item.Url, item.VulnType))
	db := ctx.Value("db").(database.Database)
	db.AddInfo(item)

}
