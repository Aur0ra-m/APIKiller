package core

import (
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/database"
	"APIKiller/internal/core/module"
	"APIKiller/internal/core/notify"
	"APIKiller/internal/core/origin"
	"APIKiller/pkg/logger"
	"fmt"
	"strings"
)

var notifier notify.Notify

func NewHandler(httpItem *origin.TransferItem) {
	r := httpItem.Req

	// assembly DataItem
	item := &data.DataItem{
		Id:             "",
		Domain:         r.Host,
		Url:            r.URL.Path,
		Https:          r.URL.Scheme == "https",
		Method:         r.Method,
		SourceRequest:  r,
		SourceResponse: httpItem.Resp,
		VulnType:       "",
		VulnRequest:    nil,
		VulnResponse:   nil,
		ReportTime:     "",
		CheckState:     false,
	}

	// enum all modules and detect
	modules := module.GetModules()
	for i, _ := range modules {
		go func(x int) {
			resultDataItem := modules[x].Detect(item)

			// exist vulnerable
			if resultDataItem != nil {
				if strings.Index(resultDataItem.VulnType, module.AsyncDetectVulnTypeSeperator) <= 0 {
					logger.Infoln(fmt.Sprintf("[Found Vulnerability] %s%s-->%s", resultDataItem.Domain, resultDataItem.Url, resultDataItem.VulnType))
					// create notification
					notify.CreateNotification(resultDataItem)
				}
				// save result
				database.CreateSaveTask(resultDataItem)
			}
		}(i)
	}

}
