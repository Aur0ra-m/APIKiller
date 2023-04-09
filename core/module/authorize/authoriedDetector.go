package authorize

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"APIKiller/util"
	"fmt"
	"github.com/antlabs/strsim"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	BYPASSED = true
	FAILED   = false
)

type AuthorizedDetector struct {
	authHeader       string
	Roles            []string
	blackStatusCodes []int
	blackKeywords    []string
	records          []string
	midPaddings      []string
	endPaddings      []string
	ipHeaders        []string
	ip               string
	apiVersionFormat string
	apiVersionPrefix string

	mu sync.Mutex
}

func (d *AuthorizedDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	logger.Debugln("[Detect] authorized detect")

	resultDataItem := d.unauthorizedDetect(item)
	if resultDataItem != nil {
		return resultDataItem
	}

	resultDataItem = d.multiRolesDetect(item)
	if resultDataItem != nil {
		return resultDataItem
	}

	return nil
}

// unauthorizedDetect
//
//	@Description: unauthorized header detect
//	@receiver d
//	@param
//	@param item
func (d *AuthorizedDetector) unauthorizedDetect(item *data.DataItem) (result *data.DataItem) {
	newReq := http2.RequestClone(item.SourceRequest)

	// delete auth header
	http2.RemoveHeader(newReq, d.authHeader)

	// make request and judge
	newResp := http2.DoRequest(newReq)

	if d.judge(item.SourceResponse, newResp) == BYPASSED {
		return util.BuildResult(item, "unauthorized", newReq, newResp)
	}

	return nil
}

// multiRolesDetect
//
//	@Description: multiple roles detect
//	@receiver d
//	@param
//	@param item
func (d *AuthorizedDetector) multiRolesDetect(item *data.DataItem) (result *data.DataItem) {
	newReq := http2.RequestClone(item.SourceRequest)

	// default support one role
	newRole := d.Roles[0]

	// change auth header
	newReq.Header.Set(d.authHeader, newRole)

	// do request
	newResp := http2.DoRequest(newReq)

	// judge
	if d.judge(item.SourceResponse, newResp) == BYPASSED {
		return util.BuildResult(item, "authorize-multiRoles", newReq, newResp)
	}

	// bypass
	result = d.bypass(item)
	if result != nil {
		return
	}

	return nil
}

func (d *AuthorizedDetector) bypass(item *data.DataItem) (result *data.DataItem) {
	srcReq := item.SourceRequest
	srcResp := item.SourceResponse

	request := http2.RequestClone(srcReq)
	var (
		vulnRequest  *http.Request
		vulnResponse *http.Response
	)

	// add ip headers and assign d value of 127.0.0.1 for each header
	for _, header := range d.ipHeaders {
		http2.AppendHeader(request, header, d.ip)
	}

	// architecture layer detect
	d.mu.Lock()

	if slices.Contains(d.records, srcReq.Host) == false {
		// record
		d.records = append(d.records, srcReq.Host)

		// path bypass
		if vulnRequest == nil {
			vulnRequest, vulnResponse = d.pathBypass(request, srcResp)
		}

		// protocol version
	}

	d.mu.Unlock()

	// api layer
	if vulnRequest == nil {
		// api version bypass
		vulnRequest, vulnResponse = d.apiVersionBypass(request, srcResp)
	}

	if vulnRequest != nil {
		return util.BuildResult(item, "authorize-bypass", vulnRequest, vulnResponse)
	}

	return nil
}

//
// apiVersionBypass
//  @Description:
//  @receiver a
//  @param
//  @param srcRequest
//  @param srcResponse
//  @param t
//  @param v
//  @return *http.Request
//  @return *http.Response
//
func (d *AuthorizedDetector) apiVersionBypass(srcRequest *http.Request, srcResponse *http.Response) (*http.Request, *http.Response) {

	compiler, _ := regexp.Compile("/" + d.apiVersionFormat + "/")
	foundString := compiler.FindString(srcRequest.URL.Path)

	if foundString != "" {
		// get api version
		trimedString := strings.Trim(foundString, "/")
		version, _ := strconv.Atoi(strings.Trim(trimedString, d.apiVersionPrefix))

		for i := 1; i < version; i++ {
			requestClone := http2.RequestClone(srcRequest)
			http2.ModifyURLPathAPIVerion(requestClone.URL, foundString, fmt.Sprintf("/%s%d/", d.apiVersionPrefix, i))
			response := http2.DoRequest(requestClone)
			if d.judge(srcResponse, response) {
				return requestClone, response
			}
		}

	}

	return nil, nil
}

//
// pathBypass
//  @Description:
//  @receiver a
//  @param
//  @param srcRequest
//  @param srcResponse
//  @param t
//  @param v
//  @return *http.Request
//  @return *http.Response
//
func (d *AuthorizedDetector) pathBypass(srcReq *http.Request, srcResp *http.Response) (*http.Request, *http.Response) {
	var requestClone *http.Request

	// filter url path equals "/"
	if srcReq.URL.Path == "/" {
		return nil, nil
	}

	// /admin/get --> /admin/Get
	requestClone = http2.RequestClone(srcReq)
	http2.ModifyURLPathCase(requestClone.URL)
	response := http2.DoRequest(requestClone)
	if d.judge(srcResp, response) {
		return requestClone, response
	}

	// /admin/get --> /admin/./get
	for _, midPadding := range d.midPaddings {
		requestClone = http2.RequestClone(srcReq)
		http2.ModifyURLPathMidPad(requestClone.URL, midPadding)
		response := http2.DoRequest(requestClone)
		if d.judge(srcResp, response) {
			return requestClone, response
		}
	}

	// /admin/get --> /admin/get;.css
	for _, endPadding := range d.endPaddings {
		requestClone = http2.RequestClone(srcReq)
		http2.ModifyURLPathEndPad(requestClone.URL, endPadding)
		response := http2.DoRequest(requestClone)
		if d.judge(srcResp, response) {
			return requestClone, response
		}
	}
	return nil, nil
}

// judge
//
//	@Description: Judging whether there is an ultra vires
//	@param sourceResp
//	@param newResp
//	@return bool true-->bypass, false-->fail
func (d *AuthorizedDetector) judge(srcResp, newResp *http.Response) bool {

	for _, code := range d.blackStatusCodes {
		if newResp.StatusCode == code {
			return FAILED
		}
	}

	// get body string
	newBody, _ := ioutil.ReadAll(newResp.Body)

	// keywords matching on the response body
	for _, split := range d.blackKeywords {
		if strings.Index(string(newBody), split) != -1 {
			return FAILED
		}
	}

	// textual similarity
	srcBody, _ := ioutil.ReadAll(srcResp.Body)
	sim := strsim.Compare(string(srcBody), string(newBody))
	if sim > 0.9 {
		return BYPASSED
	}

	return FAILED
}

func NewAuthorizedDetector() module.Detecter {
	if viper.GetInt("app.module.authorizedDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] authorized module")

	if len(viper.GetStringSlice("app.module.authorizedDetector.roles")) == 0 {
		logger.Errorln("no role set")
		panic("no role set")
	}

	return &AuthorizedDetector{
		authHeader:       viper.GetString("app.module.authorizedDetector.authHeader"),
		Roles:            viper.GetStringSlice("app.module.authorizedDetector.roles"),
		blackStatusCodes: viper.GetIntSlice("app.module.authorizedDetector.judgement.blackStatusCodes"),
		blackKeywords:    viper.GetStringSlice("app.module.authorizedDetector.judgement.blackKeywords"),
		records:          nil,
		midPaddings:      viper.GetStringSlice("app.module.authorizedDetector.pathFuzz.midPadding"),
		endPaddings:      viper.GetStringSlice("app.module.authorizedDetector.pathFuzz.endPadding"),
		ipHeaders:        viper.GetStringSlice("app.module.authorizedDetector.ipHeader"),
		ip:               viper.GetString("app.module.authorizedDetector.ip"),
		apiVersionFormat: viper.GetString("app.module.authorizedDetector.apiVersion.format"),
		apiVersionPrefix: viper.GetString("app.module.authorizedDetector.apiVersion.prefix"),
		mu:               sync.Mutex{},
	}
}
