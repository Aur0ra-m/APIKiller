package DoSDetector

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"context"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type resourceSizeDetector struct {
	sizeParams []string
}

func (d *resourceSizeDetector) Detect(ctx context.Context, item *data.DataItem) {
	srcReq := item.SourceRequest
	srcResp := item.SourceResponse

	// get params
	for key, _ := range srcReq.URL.Query() {
		if slices.Contains(d.sizeParams, key) {
			// clone a new newReq
			newReq := ahttp.RequestClone(srcReq)

			// replace value of params in new newReq
			ahttp.ModifyURLQueryParameter(newReq, key, []string{"10000"}) // 435-->43500

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
	}

	// find parameter in body if it is a post request
	if srcReq.Method != "POST" {
		return
	}

	bytes, _ := ioutil.ReadAll(srcReq.Body)
	var param string
	for _, param = range d.sizeParams {
		if strings.Contains(string(bytes), param) {
			break
		}
	}

	// clone a new newReq
	newReq := ahttp.RequestClone(srcReq)

	// replace value of params in new newReq
	if srcReq.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		// post form
		ahttp.ModifyPostFormParameter(newReq, param, "10000")
	} else if srcReq.Header.Get("Content-Type") == "application/json" {
		// json data
		ahttp.ModifyPostJsonParameter(newReq, param, "10000")
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
