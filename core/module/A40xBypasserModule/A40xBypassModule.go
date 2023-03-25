package A40xBypasserModule

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type A40xBypassModule struct {
	typeFlag         string
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

func (d *A40xBypassModule) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] 40x bypass")
	srcReq := item.SourceRequest
	srcResp := item.SourceResponse

	// t represent the auth failed type(e.g. code:1, body:2), and the v represent value
	t := 0
	var v interface{}

	if !slices.Contains(d.codeFlags, srcResp.StatusCode) {

		return
	}
	if t == 0 && srcResp.Body != nil {
		allbytes, _ := ioutil.ReadAll(srcResp.Body)
		body := string(allbytes)

		for _, keyword := range d.bodyFlags {
			if strings.Contains(body, keyword) {
				t = 2
				v = keyword
			}
		}
	}

	// filter out requests that fail authentication
	if t == 0 {
		logger.Debugln("not request that fail authentication")
		return
	}

	request := ahttp.RequestClone(srcReq)
	var (
		vulnRequest  *http.Request
		vulnResponse *http.Response
	)

	// add ip headers and assign d value of 127.0.0.1 for each header
	d.appendIpHeaders(request)

	// architecture layer detect
	d.mu.Lock()

	if slices.Contains(d.records, srcReq.Host) == false {
		// record
		d.records = append(d.records, srcReq.Host)

		// path bypass
		if vulnRequest == nil {
			vulnRequest, vulnResponse = d.pathBypass(ctx, request, srcResp, t, v)
		}

		// protocol version
	}

	d.mu.Unlock()

	// api layer
	if vulnRequest == nil {
		// api version bypass
		vulnRequest, vulnResponse = d.apiVersionBypass(ctx, request, srcResp, t, v)
	}

	// save result
	if vulnRequest != nil {
		item.VulnType = append(item.VulnType, "40x bypass")
		item.VulnRequest = append(item.VulnRequest, vulnRequest)
		item.VulnResponse = append(item.VulnResponse, vulnResponse)
	}
}

//
// pathBypass
//  @Description:
//  @receiver a
//  @param ctx
//  @param srcRequest
//  @param srcResponse
//  @param t
//  @param v
//  @return *http.Request
//  @return *http.Response
//
func (d *A40xBypassModule) pathBypass(ctx context.Context, srcRequest *http.Request, srcResponse *http.Response, t int, v interface{}) (*http.Request, *http.Response) {
	var requestClone *http.Request

	// /admin/get --> /admin/Get
	requestClone = ahttp.RequestClone(srcRequest)
	ahttp.ModifyURLPathCase(requestClone.URL)
	response := ahttp.DoRequest(requestClone)
	if d.judge(ctx, t, v, srcResponse, response) {
		return requestClone, response
	}

	// /admin/get --> /admin/./get
	for _, midPadding := range d.midPaddings {
		requestClone = ahttp.RequestClone(srcRequest)
		ahttp.ModifyURLPathMidPad(requestClone.URL, midPadding)
		response := ahttp.DoRequest(requestClone)
		if d.judge(ctx, t, v, srcResponse, response) {
			return requestClone, response
		}
	}

	// /admin/get --> /admin/get;.css
	for _, endPadding := range d.endPaddings {
		requestClone = ahttp.RequestClone(srcRequest)
		ahttp.ModifyURLPathEndPad(requestClone.URL, endPadding)
		response := ahttp.DoRequest(requestClone)
		if d.judge(ctx, t, v, srcResponse, response) {
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
func (d *A40xBypassModule) appendIpHeaders(req *http.Request) {
	// assemble request
	for _, header := range d.ipHeaders {
		ahttp.AppendHeader(req, header, d.ip)
	}
}

//
// apiVersionBypass
//  @Description:
//  @receiver a
//  @param ctx
//  @param srcRequest
//  @param srcResponse
//  @param t
//  @param v
//  @return *http.Request
//  @return *http.Response
//
func (d *A40xBypassModule) apiVersionBypass(ctx context.Context, srcRequest *http.Request, srcResponse *http.Response, t int, v interface{}) (*http.Request, *http.Response) {

	compiler, _ := regexp.Compile("/" + d.apiVersionFormat + "/")
	foundString := compiler.FindString(srcRequest.URL.Path)

	if foundString != "" {
		// get api version
		trimedString := strings.Trim(foundString, "/")
		version, _ := strconv.Atoi(strings.Trim(trimedString, d.apiVersionPrefix))

		for i := 1; i < version; i++ {
			requestClone := ahttp.RequestClone(srcRequest)
			ahttp.ModifyURLPathAPIVerion(requestClone.URL, foundString, fmt.Sprintf("/%s%d/", d.apiVersionPrefix, i))
			response := ahttp.DoRequest(requestClone)
			if d.judge(ctx, t, v, srcResponse, response) {
				return requestClone, response
			}
		}

	}

	return nil, nil
}

//
// judge
//  @Description: Judging whether bypass successfully
//  @receiver a
//  @param ctx
//  @param t type of auth failed flag(e.g. code:1, body:2)
//  @param value
//  @param srcResp
//  @param newResp
//  @return bool
//
func (d *A40xBypassModule) judge(ctx context.Context, t int, value interface{}, srcResp, newResp *http.Response) bool {
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

func NewA40xBypassModule(ctx context.Context) module.Detecter {
	if viper.GetInt("app.module.40xBypassModule.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] 40x bypass module")

	return &A40xBypassModule{
		typeFlag:         viper.GetString("app.module.40xBypassModule.typeFlag"),
		records:          nil,
		midPaddings:      viper.GetStringSlice("app.module.40xBypassModule.pathFuzz.midPadding"),
		endPaddings:      viper.GetStringSlice("app.module.40xBypassModule.pathFuzz.endPadding"),
		codeFlags:        viper.GetIntSlice("app.module.40xBypassModule.authFailedFlag.statusCode"),
		bodyFlags:        viper.GetStringSlice("app.module.40xBypassModule.authFailedFlag.body"),
		ipHeaders:        viper.GetStringSlice("app.module.40xBypassModule.ipHeader"),
		ip:               viper.GetString("app.module.40xBypassModule.ip"),
		apiVersionFormat: viper.GetString("app.module.40xBypassModule.apiVersion.format"),
		apiVersionPrefix: viper.GetString("app.module.40xBypassModule.apiVersion.prefix"),
		mu:               sync.Mutex{},
	}
}
