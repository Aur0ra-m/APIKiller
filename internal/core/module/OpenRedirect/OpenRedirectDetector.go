package OpenRedirect

import (
	ahttp2 "APIKiller/internal/core/ahttp"
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/module"
	"APIKiller/pkg/logger"
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

func (d *OpenRedirectDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	logger.Debugln("[Detect] Open-Redirect detect")

	srcResp := item.SourceResponse
	srcReq := item.SourceRequest

	// filter by features of redirect
	if srcResp.Header.Get("Location") == "" {
		return
	}

	newReq := ahttp2.ModifyQueryParamByRegExp(srcReq, `https?://[^\s&]+`, d.evilLink)
	if newReq == nil {
		logger.Debugln("parameter not found")
		return
	}

	// do newReq
	newResp := ahttp2.DoRequest(newReq)

	// judge
	if d.judge(srcResp, newResp) {
		return data.BuildResult(item, "Open-Redirect", newReq, newResp)
	}

	return nil
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

	return false
}

func NewOpenRedirectDetector() module.Detecter {
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
