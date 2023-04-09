package SSRF

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	"APIKiller/logger"
	"APIKiller/util"
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

	token := util.GenerateRandomId()

	newReq := ahttp.ModifyQueryParamByRegExp(srcReq, `https?://[^\s&]+`, d.ReverseConnectionPlatform+fmt.Sprintf("%s%s%s", "SSRF", module.AsyncDetectVulnTypeSeperator, token))
	if newReq == nil {
		logger.Debugln("parameter not found")
		return
	}

	// do newReq
	newResp := ahttp.DoRequest(newReq)

	// asynchronous result
	return util.BuildResult(item, "SSRF"+module.AsyncDetectVulnTypeSeperator+token, newReq, newResp)
}
