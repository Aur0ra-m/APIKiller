package core

import (
	"APIKiller/core/data"
	"APIKiller/core/database"
	"APIKiller/core/module"
	"APIKiller/core/notify"
	"APIKiller/core/origin"
	"APIKiller/logger"
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
	modules := module.Modules
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
