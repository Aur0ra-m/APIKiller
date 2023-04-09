package DoS

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/util"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

type resourceSizeDetector struct {
	sizeParams []string
}

func (d *resourceSizeDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	srcReq := item.SourceRequest
	srcResp := item.SourceResponse

	for _, param := range d.sizeParams {

		// replace value of params in new newReq
		newReq := ahttp.ModifyParam(srcReq, param, "10000")
		if newReq == nil {
			continue
		}

		// do newReq
		newResp := ahttp.DoRequest(newReq)

		if newResp == nil {
			return
		}

		// judge
		if d.judge(srcResp, newResp) {
			return util.BuildResult(item, "DoS-ResourceSizeNotStrict", newReq, newResp)
		}

		return nil
	}
	return nil

}

func (d *resourceSizeDetector) judge(srcResp, newResp *http.Response) bool {
	srcCL, _ := strconv.Atoi(srcResp.Header.Get("Content-Length"))
	newCL, _ := strconv.Atoi(newResp.Header.Get("Content-Length"))
	if newCL/10 > srcCL { // successfully
		return true
	}
	return false
}

func newResourceSizeDetector() *resourceSizeDetector {
	return &resourceSizeDetector{
		sizeParams: viper.GetStringSlice("app.module.DoSDetector.sizeParam"),
	}
}
