package openRedirectDetector

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"context"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

type OpenRedirectDetector struct {
	rawQueryParams []string
	failFlag       []string
	evilLink       string
}

func (d *OpenRedirectDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] Open-Redirect detect")

	srcResp := item.SourceResponse
	srcReq := item.SourceRequest

	for _, param := range d.rawQueryParams {
		newReq := ahttp.ModifyQueryParam(srcReq, param, d.evilLink)
		if newReq == nil {
			logger.Debugln("parameter not found")
			continue
		}

		// do newReq
		newResp := ahttp.DoRequest(newReq)

		// judge
		if d.judge(srcResp, newResp) {
			item.VulnType = append(item.VulnType, "open-redirect")
			item.VulnRequest = append(item.VulnRequest, newReq)
			item.VulnResponse = append(item.VulnResponse, newResp)
		}

		return
	}
}

func (d *OpenRedirectDetector) judge(srcResp, newResp *http.Response) bool {
	if newResp.StatusCode == srcResp.StatusCode && strings.Index(newResp.Header.Get("Location"), d.evilLink) != -1 {
		return true
	}

	// black list
	if newResp.Body != nil {
		bytes, _ := ioutil.ReadAll(newResp.Body)
		for _, flag := range d.failFlag {
			if strings.Contains(string(bytes), flag) {
				return false
			}
		}
	}

	return true
}

func NewOpenRedirectDetector(ctx context.Context) module.Detecter {
	if viper.GetInt("app.module.openRedirectDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] Open-Redirect detect module")

	d := &OpenRedirectDetector{
		rawQueryParams: viper.GetStringSlice("app.module.openRedirectDetector.rawQueryParams"),
		evilLink:       "https://www.baidu.com",
		failFlag:       viper.GetStringSlice("app.module.openRedirectDetector.failFlag"),
	}

	return d
}
