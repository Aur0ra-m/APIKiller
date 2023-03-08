package openRedirectDetector

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/log"
	"context"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
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

	// get raw query parameter list
	for key, _ := range srcReq.URL.Query() {
		if slices.Contains(d.rawQueryParams, key) {
			// clone a new newReq
			newReq := ahttp.RequestClone(srcReq)

			// replace value of params in new newReq
			newReq.URL.Query()[key] = []string{d.evilLink}

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
