package util

import (
	"APIKiller/internal/core/data"
	"fmt"
	"net/http"
	"time"
)

//
// BuildResult
//  @Description: build result through create a new *data.DateItem
//  @param dataItem
//  @param vulnType
//  @param vulnReq
//  @param vulnResp
//  @return *data.DataItem
//
func BuildResult(dataItem *data.DataItem, vulnType string, vulnReq *http.Request, vulnResp *http.Response) *data.DataItem {
	return &data.DataItem{
		Id:             GenerateRandomId(),
		Domain:         dataItem.Domain,
		Url:            dataItem.Url,
		Method:         dataItem.Method,
		Https:          dataItem.Https,
		SourceRequest:  dataItem.SourceRequest,
		SourceResponse: dataItem.SourceResponse,
		VulnType:       vulnType,
		VulnRequest:    vulnReq,
		VulnResponse:   vulnResp,
		ReportTime:     fmt.Sprintf("%v", time.Now().Unix()),
		CheckState:     false,
	}
}
