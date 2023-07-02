package SSRF

import (
	ahttp2 "APIKiller/internal/core/ahttp"
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/module"
	"APIKiller/pkg/logger"
	util2 "APIKiller/pkg/util"
	"fmt"
	"github.com/spf13/viper"
)

type SSRFDetector struct {
	ReverseConnectionPlatform string // end with “/”
}

func NewSSRFDetector() module.Detecter {
	if viper.GetInt("app.module.SSRFDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] SSRF detect module")

	return &SSRFDetector{
		ReverseConnectionPlatform: "http://zpysri.ceye.io/",
	}
}

func (d *SSRFDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	logger.Debugln("[Detect] SSRF detect")

	//srcResp := item.SourceResponse
	srcReq := item.SourceRequest

	token := util2.GenerateRandomId()

	newReq := ahttp2.ModifyQueryParamByRegExp(srcReq, `https?://[^\s&]+`, d.ReverseConnectionPlatform+fmt.Sprintf("%s%s%s", "SSRF", module.AsyncDetectVulnTypeSeperator, token))
	if newReq == nil {
		logger.Debugln("parameter not found")
		return
	}

	// do newReq
	newResp := ahttp2.DoRequest(newReq)

	// asynchronous result
	return util2.BuildResult(item, "SSRF"+module.AsyncDetectVulnTypeSeperator+token, newReq, newResp)
}
