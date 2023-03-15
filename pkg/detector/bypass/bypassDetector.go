package bypass

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/detector"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/types"
	"fmt"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type BypassDetector struct {
	records          []string
	midPaddings      []string
	endPaddings      []string
	codeFlags        []int
	bodyFlags        []string
	ipHeaders        []string
	ip               string
	apiVersionFormat string
	apiVersionPrefix string

	mu sync.Mutex
}

func (b *BypassDetector) Detect(item *types.DataItem) {
	srcReq := item.SourceRequest
	srcResp := item.SourceResponse

	// t represent the auth failed type(e.g. code:1, body:2), and the v represent value
	t := 0
	var v interface{}
	if !slices.Contains(b.codeFlags, srcResp.StatusCode) {

		return
	}
	if t == 0 && srcResp.Body != nil {
		allbytes, _ := ioutil.ReadAll(srcResp.Body)
		body := string(allbytes)

		for _, keyword := range b.bodyFlags {
			if strings.Contains(body, keyword) {
				t = 2
				v = keyword
			}
		}
	}

	// filter out requests that fail authentication
	if t == 0 {
		//logger.Debugln("not request that fail authentication")
		return
	}

	request := gohttp.RequestClone(srcReq)
	var (
		vulnRequest  *http.Request
		vulnResponse *http.Response
	)

	// add ip headers and assign a value of 127.0.0.1 for each header
	b.appendIpHeaders(request)

	// architecture layer detect
	b.mu.Lock()

	if slices.Contains(b.records, srcReq.Host) == false {
		// record
		b.records = append(b.records, srcReq.Host)

		// path bypass
		if vulnRequest == nil {
			vulnRequest, vulnResponse = b.pathBypass(request, srcResp, t, v)
		}

		// protocol version
	}

	b.mu.Unlock()

	// api layer
	if vulnRequest == nil {
		// api version bypass
		vulnRequest, vulnResponse = b.apiVersionBypass(request, srcResp, t, v)
	}

	// save result
	if vulnRequest != nil {
		item.VulnType = append(item.VulnType, "40x bypass")
		item.VulnRequest = append(item.VulnRequest, vulnRequest)
		item.VulnResponse = append(item.VulnResponse, vulnResponse)
	}
}

// pathBypass
//
//	@Description:
//	@receiver a
//	@param ctx
//	@param srcRequest
//	@param srcResponse
//	@param t
//	@param v
//	@return *http.Request
//	@return *http.Response
func (b *BypassDetector) pathBypass(srcRequest *http.Request, srcResponse *http.Response, t int, v interface{}) (*http.Request, *http.Response) {
	var requestClone *http.Request

	// /admin/get --> /admin/Get
	requestClone = gohttp.RequestClone(srcRequest)
	gohttp.ModifyURLPathCase(requestClone.URL)
	response := gohttp.DoRequest(requestClone)
	if b.judge(t, v, srcResponse, response) {
		return requestClone, response
	}

	// /admin/get --> /admin/./get
	for _, midPadding := range b.midPaddings {
		requestClone = gohttp.RequestClone(srcRequest)
		gohttp.ModifyURLPathMidPad(requestClone.URL, midPadding)
		response := gohttp.DoRequest(requestClone)
		if b.judge(t, v, srcResponse, response) {
			return requestClone, response
		}
	}

	// /admin/get --> /admin/get;.css
	for _, endPadding := range b.endPaddings {
		requestClone = gohttp.RequestClone(srcRequest)
		gohttp.ModifyURLPathEndPad(requestClone.URL, endPadding)
		response := gohttp.DoRequest(requestClone)
		if b.judge(t, v, srcResponse, response) {
			return requestClone, response
		}
	}

	return nil, nil
}

// appendIpHeaders
//
//	@Description:  append all possible headers to request
//	@receiver a
//	@param req
func (b *BypassDetector) appendIpHeaders(req *http.Request) {
	// assemble request
	for _, header := range b.ipHeaders {
		req.Header.Add(header, b.ip)
	}
}

// apiVersionBypass
//
//	@Description:
//	@receiver a
//	@param ctx
//	@param srcRequest
//	@param srcResponse
//	@param t
//	@param v
//	@return *http.Request
//	@return *http.Response
func (b *BypassDetector) apiVersionBypass(srcRequest *http.Request, srcResponse *http.Response, t int, v interface{}) (*http.Request, *http.Response) {

	compiler, _ := regexp.Compile("/" + b.apiVersionFormat + "/")
	foundString := compiler.FindString(srcRequest.URL.Path)

	if foundString != "" {
		// get api version
		trimedString := strings.Trim(foundString, "/")
		version, _ := strconv.Atoi(strings.Trim(trimedString, b.apiVersionPrefix))

		for i := 1; i < version; i++ {
			requestClone := gohttp.RequestClone(srcRequest)
			gohttp.ModifyURLPathAPIVerion(requestClone.URL, foundString, fmt.Sprintf("/%s%d/", b.apiVersionPrefix, i))
			response := gohttp.DoRequest(requestClone)
			if b.judge(t, v, srcResponse, response) {
				return requestClone, response
			}
		}

	}

	return nil, nil
}

// judge
//
//	@Description: Judging whether bypass successfully
//	@receiver a
//	@param ctx
//	@param t type of auth failed flag(e.g. code:1, body:2)
//	@param value
//	@param srcResp
//	@param newResp
//	@return bool
func (b *BypassDetector) judge(t int, value interface{}, srcResp, newResp *http.Response) bool {
	if t == 1 {
		//
	} else if t == 2 {
		if newResp.Body != nil {
			bytes, _ := ioutil.ReadAll(newResp.Body)
			newBody := string(bytes)

			if strings.Contains(newBody, value.(string)) {
				return false
			}
		}
	}

	if newResp.StatusCode >= 300 {
		return false
	}

	return true
}

func New40xBypassDetector(cfg *config.Config) detector.Detector {
	bypassCfg := cfg.Detector.A40xBypass
	if !bypassCfg.Enable {
		return nil
	}

	return &BypassDetector{
		records:          nil,
		midPaddings:      bypassCfg.PathFuzz["midPaddings"].([]string),
		endPaddings:      bypassCfg.PathFuzz["endPaddings"].([]string),
		codeFlags:        bypassCfg.AuthFailFlag["statusCode"].([]int),
		bodyFlags:        bypassCfg.AuthFailFlag["body"].([]string),
		ipHeaders:        bypassCfg.IpHeader,
		ip:               bypassCfg.Ip,
		apiVersionFormat: bypassCfg.ApiVersion["format"],
		apiVersionPrefix: bypassCfg.ApiVersion["prefix"],
		mu:               sync.Mutex{},
	}
}
