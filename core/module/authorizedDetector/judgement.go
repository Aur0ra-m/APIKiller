package authorizedDetector

import (
	"APIKiller/util"
	"context"
	"github.com/antlabs/strsim"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	Bypass = true
	Fail   = false
)

//
// Judge
//  @Description: Judging whether there is an ultra vires
//  @param sourceResp
//  @param newResp
//  @return bool true-->bypass, false-->fail
//
func Judge(ctx context.Context, srcResp, newResp *http.Response) bool {
	// status code
	blackStatusCodes := util.GetConfig(ctx, "app.detectors.authorizedDetector.judgement.blackStatusCodes")
	codes := util.SplitConfigString(blackStatusCodes)
	for _, code := range codes {
		if strings.Index(newResp.Status, code) != -1 {
			return Fail
		}
	}

	// get body string
	newBody, _ := ioutil.ReadAll(newResp.Body)

	// keywords matching on the response body
	blackKeywords := util.GetConfig(ctx, "app.detectors.authorizedDetector.judgement.blackKeywords")
	splits := util.SplitConfigString(blackKeywords)
	for _, split := range splits {
		if strings.Index(string(newBody), split) != -1 {
			return Fail
		}
	}

	// textual similarity
	srcBody, _ := ioutil.ReadAll(srcResp.Body)
	sim := strsim.Compare(string(srcBody), string(newBody))
	if sim > 0.9 {
		return Bypass
	}

	return Fail
}
