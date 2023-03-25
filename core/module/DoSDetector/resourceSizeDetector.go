package DoSDetector

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"context"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

type resourceSizeDetector struct {
	sizeParams []string
}

func (d *resourceSizeDetector) Detect(ctx context.Context, item *data.DataItem) {
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

		// judge
		if d.judge(srcResp, newResp) {
			item.VulnType = append(item.VulnType, "DoS-ResourceSizeNotStrict")
			item.VulnRequest = append(item.VulnRequest, newReq)
			item.VulnResponse = append(item.VulnResponse, newResp)
		}

		return
	}
	return

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
