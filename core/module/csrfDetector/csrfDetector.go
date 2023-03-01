package csrfDetector

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/log"
	"context"
	"fmt"
	"github.com/antlabs/strsim"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type CsrfDetector struct {
	csrfTokenPattern   string
	csrfInvalidPattern string
	samesitePolicy     map[string]string
	mu                 sync.Mutex
}

func (d *CsrfDetector) Detect(ctx context.Context, item *data.DataItem) {
	resp := item.SourceResponse
	req := item.SourceRequest

	// same-site check with lock
	d.mu.Lock()
	if d.samesitePolicy[req.Host] == "" {
		d.getSameSitePolicy(ctx, item)
	}
	d.mu.Unlock()

	policy := d.samesitePolicy[req.Host]
	if policy == "Strict" {
		return
	} else if policy == "Lax" && item.Method != "GET" {
		return
	} else {
		// no same-site policy or the policy is fail
	}

	// cors--Access-Control-Allow-Origin
	value := resp.Header.Get("Access-Control-Allow-Origin")
	if value != "" && value != "*" {
		return
	}

	// copy request
	request := http2.RequestClone(req)

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

// getSameSitePolicy
//
//	@Description: get same-site policy from response received from the request deleted cookie
//	@receiver d
//	@param ctx
//	@param item
func (d *CsrfDetector) getSameSitePolicy(ctx context.Context, item *data.DataItem) {
	// copy request
	request := http2.RequestClone(item.SourceRequest)
	// delete cookie and get set-cookie header from response
	request.Header.Del("Cookie")
	response := http2.DoRequest(request)
	setCookie := response.Header.Get("Set-Cookie")

	var policy string
	// parse policy from Set-Cookie header
	if strings.Contains(setCookie, "SameSite=Lax") {
		policy = "Lax"
	} else if strings.Contains(setCookie, "SameSite=Strict") {
		policy = "Strict"
	} else { // if there is not same-site policy or the policy is None
		policy = "None"
	}

	// save policy to samesitePolicy
	key := request.Host
	d.samesitePolicy[key] = policy

	logger.Infoln(fmt.Sprintf("Host: %s, Same-Site policy: %s", key, policy))
}

// judge
//
//	@Description:
//	@receiver d
//	@return bool  true -- exists vulnerable
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
	if viper.GetInt("app.module.csrfDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] csrf detector module")

	// instantiate csrfDetector
	detector := &CsrfDetector{
		csrfTokenPattern:   viper.GetString("app.module.csrfDetector.csrfTokenPattern"),
		csrfInvalidPattern: viper.GetString("app.module.csrfDetector.csrfInvalidPattern"),
		samesitePolicy:     make(map[string]string, 100),
	}

	return detector
}
