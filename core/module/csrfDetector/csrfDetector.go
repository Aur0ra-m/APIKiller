package csrfDetector

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	"APIKiller/util"
	"context"
	"github.com/antlabs/strsim"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type CsrfDetector struct {
	csrfTokenPattern   string
	csrfInvalidPattern string
}

func (d *CsrfDetector) Detect(ctx context.Context, item *data.DataItem) {
	resp := item.SourceResponse

	// same-site

	// cors--Access-Control-Allow-Origin
	value := resp.Header.Get("Access-Control-Allow-Origin")
	if value != "" && value != "*" {
		return
	}

	// copy request
	request := http2.RequestClone(item.SourceRequest)

	// delete referer and origin
	if request.Header.Get("Referer") != "" {
		request.Header.Del("Referer")
	}

	if request.Header.Get("Origin") != "" {
		request.Header.Del("Origin")
	}

	// find token position and detect before delete csrf token
	// 1. row query

	if request.URL.RawQuery != "" {
		editedRawQuery := []string{}

		for _, kv := range strings.Split(request.URL.RawQuery, "&") {
			splits := strings.Split(kv, "=")
			key := splits[0]
			//value := splits[1]

			match, _ := regexp.Match(d.csrfTokenPattern, []byte(key))
			if match {
				continue
			}

			// add to editedRawQuery
			editedRawQuery = append(editedRawQuery, kv)
		}
		request.URL.RawQuery = strings.Join(editedRawQuery, "&")
	}

	// 2. post body(application/x-www-form-urlencoded,multipart/form-data )
	for k, _ := range request.PostForm {
		match, _ := regexp.Match(d.csrfTokenPattern, []byte(k))
		if match {
			request.PostForm.Del(k)
		}
	}

	// make request
	response := http2.DoRequest(request)
	if response == nil {
		return
	}

	// judge and save result
	if d.judge(resp, response) {

		item.VulnType = append(item.VulnType, "csrf")
		item.VulnRequest = append(item.VulnRequest, request)
		item.VulnResponse = append(item.VulnResponse, response)
	}
}

//
// judge
//  @Description:
//  @receiver d
//  @return bool  true -- exists vulnerable
//
func (d *CsrfDetector) judge(srcResponse, response *http.Response) bool {
	// get response body
	bytes, _ := ioutil.ReadAll(srcResponse.Body)
	bytes2, _ := ioutil.ReadAll(response.Body)

	sim := strsim.Compare(string(bytes), string(bytes2))
	if sim > 0.9 {
		return true
	}
	return false
}

func NewCsrfDetector(ctx context.Context) module.Detecter {
	if util.GetConfig(ctx, "app.detectors.csrfDetector.option") != "1" {
		return nil
	}

	// get config
	csrfTokenPattern := util.GetConfig(ctx, "app.detectors.csrfDetector.csrfTokenPattern")
	csrfInvalidPattern := util.GetConfig(ctx, "app.detectors.csrfDetector.csrfInvalidPattern")

	// instantiate csrfDetector
	detector := &CsrfDetector{
		csrfTokenPattern:   csrfTokenPattern,
		csrfInvalidPattern: csrfInvalidPattern,
	}

	return detector
}
