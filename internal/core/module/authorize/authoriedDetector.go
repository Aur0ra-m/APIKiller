package authorize

import (
	"APIKiller/internal/core/ahttp"
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/module"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/util"
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
	authGroups       []authGroup
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

	// match auth group
	var group = authGroup{}
	for _, group = range d.authGroups {
		if slices.Contains(group.Domain, item.SourceRequest.Host) {
			break
		}
	}

	resultDataItem := d.unauthorizedDetect(item, &group)
	if resultDataItem != nil {
		return resultDataItem
	}

	resultDataItem = d.multiRolesDetect(item, &group)
	if resultDataItem != nil {
		return resultDataItem
	}

	return nil
}

//
// deleteAuthIdentifier
//  @Description: delete auth identifier according to the first replace group in auth group
//  @receiver d
//  @param group
//
func (d *AuthorizedDetector) deleteAuthIdentifier(group *authGroup, req *http.Request) {
	rg := group.ReplaceGroups[0]

	switch rg.Position {
	case Replace_Position_Code_Header:
		ahttp.RemoveHeader(req, rg.Key)
		break
	case Replace_Position_Code_Query:
		// ...
		break
	case Replace_Position_Code_Body:
		// ...
		break
	}
}

//
// replaceAuthIdentifierandNecessaryParms
//  @Description: replace auth identifier and other necessary parameters(e.g. csrf token) with replace group items
//  @receiver d
//  @param group
//  @param req
//
func (d *AuthorizedDetector) replaceAuthIdentifierandNecessaryParms(group *authGroup, req *http.Request) {
	for _, rg := range group.ReplaceGroups {
		switch rg.Position {
		case Replace_Position_Code_Header:
			ahttp.UpdateHeader(req, rg.Key, rg.Value)
			break
		case Replace_Position_Code_Query:
			ahttp.ModifyQueryParam(req, rg.Key, rg.Value)
			break
		case Replace_Position_Code_Body:
			ahttp.ModifyPostParam(req, rg.Key, rg.Value)
			break
		default:
			logger.Errorln("cannot recognize the position code")
		}
	}
}

// unauthorizedDetect
//
//	@Description: unauthorized header detect
//	@receiver d
//	@param
//	@param item
func (d *AuthorizedDetector) unauthorizedDetect(item *data.DataItem, group *authGroup) (result *data.DataItem) {
	newReq := ahttp.RequestClone(item.SourceRequest)

	// delete auth identifier
	d.deleteAuthIdentifier(group, newReq)

	// make request and judge
	newResp := ahttp.DoRequest(newReq)

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
func (d *AuthorizedDetector) multiRolesDetect(item *data.DataItem, group *authGroup) (result *data.DataItem) {
	newReq := ahttp.RequestClone(item.SourceRequest)

	// replace necessary data
	d.replaceAuthIdentifierandNecessaryParms(group, newReq)

	// do request
	newResp := ahttp.DoRequest(newReq)

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

	request := ahttp.RequestClone(srcReq)
	var (
		vulnRequest  *http.Request
		vulnResponse *http.Response
	)

	// add ip headers and assign d value of 127.0.0.1 for each header
	for _, header := range d.ipHeaders {
		ahttp.AppendHeader(request, header, d.ip)
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
			requestClone := ahttp.RequestClone(srcRequest)
			ahttp.ModifyURLPathAPIVerion(requestClone.URL, foundString, fmt.Sprintf("/%s%d/", d.apiVersionPrefix, i))
			response := ahttp.DoRequest(requestClone)
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
	requestClone = ahttp.RequestClone(srcReq)
	ahttp.ModifyURLPathCase(requestClone.URL)
	response := ahttp.DoRequest(requestClone)
	if d.judge(srcResp, response) {
		return requestClone, response
	}

	// /admin/get --> /admin/./get
	for _, midPadding := range d.midPaddings {
		requestClone = ahttp.RequestClone(srcReq)
		ahttp.ModifyURLPathMidPad(requestClone.URL, midPadding)
		response := ahttp.DoRequest(requestClone)
		if d.judge(srcResp, response) {
			return requestClone, response
		}
	}

	// /admin/get --> /admin/get;.css
	for _, endPadding := range d.endPaddings {
		requestClone = ahttp.RequestClone(srcReq)
		ahttp.ModifyURLPathEndPad(requestClone.URL, endPadding)
		response := ahttp.DoRequest(requestClone)
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
	
	detector := &AuthorizedDetector{
		authGroups:       []authGroup{},
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

	viper.UnmarshalKey("app.module.authorizedDetector.authGroup", &detector.authGroups)

	return detector
}
