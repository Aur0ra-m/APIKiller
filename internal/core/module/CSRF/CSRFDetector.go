package CSRF

import (
	http2 "APIKiller/internal/core/ahttp"
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/module"
	"APIKiller/pkg/logger"
	"github.com/antlabs/strsim"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type CSRFDetector struct {
	csrfTokenPattern   string
	csrfInvalidPattern []string
	samesitePolicy     map[string]string
	mu                 sync.Mutex
}

func (d *CSRFDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	logger.Debugln("[Detect] CSRF detect")

	srcResp := item.SourceResponse
	srcReq := item.SourceRequest

	// same-site check with lock
	d.mu.Lock()
	if d.samesitePolicy[srcReq.Host] == "" {
		d.getSameSitePolicy(item)
	}
	d.mu.Unlock()

	policy := d.samesitePolicy[srcReq.Host]
	if policy == "Strict" {
		return
	} else if policy == "Lax" && item.Method != "GET" {
		return
	} else {
		// no same-site policy or the policy is fail
	}

	// cors--Access-Control-Allow-Origin
	value := srcResp.Header.Get("Access-Control-Allow-Origin")
	if value != "" && value != "*" {
		return
	}

	// copy newReq
	newReq := http2.RequestClone(srcReq)

	// delete referer and origin
	if newReq.Header.Get("Referer") != "" {
		newReq.Header.Del("Referer")
	}

	if newReq.Header.Get("Origin") != "" {
		newReq.Header.Del("Origin")
	}

	// find token position and detect before delete csrf token
	// 1. row query

	if newReq.URL.RawQuery != "" {
		editedRawQuery := []string{}

		for _, kv := range strings.Split(newReq.URL.RawQuery, "&") {
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
		newReq.URL.RawQuery = strings.Join(editedRawQuery, "&")
	}

	// 2. post body(application/x-www-form-urlencoded,multipart/form-data )
	for k, _ := range newReq.PostForm {
		match, _ := regexp.Match(d.csrfTokenPattern, []byte(k))
		if match {
			newReq.PostForm.Del(k)
		}
	}

	// make newReq
	newResp := http2.DoRequest(newReq)
	if newResp == nil {
		return
	}

	// judge and save result
	if d.judge(srcResp, newResp) {
		return data.BuildResult(item, "CSRF", newReq, newResp)
	}

	return nil
}

// getSameSitePolicy
//
//	@Description: get same-site policy from response received from the request deleted cookie
//	@receiver d
//	@param
//	@param item
func (d *CSRFDetector) getSameSitePolicy(item *data.DataItem) {
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

	//logger.Infoln(fmt.Sprintf("Host: %s, Same-Site policy: %s", key, policy))
}

// judge
//
//	@Description:
//	@receiver d
//	@return bool  true -- exists vulnerable
func (d *CSRFDetector) judge(srcResponse, response *http.Response) bool {
	bytes2, _ := ioutil.ReadAll(response.Body)

	// black keyword match
	for _, s := range d.csrfInvalidPattern {
		if strings.Contains(string(bytes2), s) {
			return false
		}
	}

	// body similarity compare
	bytes, _ := ioutil.ReadAll(srcResponse.Body)

	sim := strsim.Compare(string(bytes), string(bytes2))
	if sim > 0.9 {
		return true
	}

	return false
}

func NewCSRFDetector() module.Detecter {
	if viper.GetInt("app.module.CSRFDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] csrf detector module")

	// instantiate CSRFDetector
	detector := &CSRFDetector{
		csrfTokenPattern:   viper.GetString("app.module.CSRFDetector.csrfTokenPattern"),
		csrfInvalidPattern: viper.GetStringSlice("app.module.CSRFDetector.csrfInvalidPattern"),
		samesitePolicy:     make(map[string]string, 100),
	}

	return detector
}
