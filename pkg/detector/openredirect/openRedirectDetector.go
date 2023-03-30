package openredirect

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/detector"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
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

func (d *OpenRedirectDetector) Detect(item *types.DataItem) {
	logger.Debug("[Detect] Open-Redirect detect")

	srcResp := item.SourceResponse
	srcReq := item.SourceRequest

	// get raw query parameter list
	for key, _ := range srcReq.URL.Query() {
		if slices.Contains(d.rawQueryParams, key) {
			// clone a new newReq
			newReq := gohttp.RequestClone(srcReq)

			// replace value of params in new newReq
			newReq.URL.Query()[key] = []string{d.evilLink}

			// do newReq
			newResp := gohttp.DoRequest(newReq)

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

func (d *OpenRedirectDetector) judge(oldResp, newResp *http.Response) bool {
	if newResp.StatusCode == oldResp.StatusCode && strings.Index(newResp.Header.Get("Location"), d.evilLink) != -1 {
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

func NewOpenRedirectDetector(cfg *config.DetectorConfig) detector.Detector {
	openCfg := cfg.OpenRedirect
	if openCfg.Enable {
		return nil
	}

	logger.Info("[Load Module] Open-Redirect detect module\n")

	return &OpenRedirectDetector{
		rawQueryParams: openCfg.RawQueryParams,
		evilLink:       "https://www.baidu.com",
		failFlag:       openCfg.FailFlag,
	}
}
