package authorizedDetector

import (
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

// Judge
//
//	@Description: Judging whether there is an ultra vires
//	@param sourceResp
//	@param newResp
//	@return bool true-->bypass, false-->fail
func (d *AuthorizedDetector) Judge(ctx context.Context, srcResp, newResp *http.Response) bool {

	for _, code := range d.blackStatusCodes {
		if newResp.StatusCode == code {
			return Fail
		}
	}

	// get body string
	newBody, _ := ioutil.ReadAll(newResp.Body)

	// keywords matching on the response body
	for _, split := range d.blackKeywords {
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
