package authorize

import (
	"github.com/antlabs/strsim"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	SUCCESS = true
	FAIL    = false
)

func (d *AuthorizedDetector) Judge(oldResp, newResp *http.Response) bool {
	for _, code := range d.blackStatusCodes {
		if newResp.StatusCode == code {
			return FAIL
		}
	}

	// get body string
	newBody, _ := ioutil.ReadAll(newResp.Body)
	// keywords matching on the response body
	for _, split := range d.blackKeywords {
		if strings.Index(string(newBody), split) != -1 {
			return FAIL
		}
	}

	// textual similarity
	srcBody, _ := ioutil.ReadAll(oldResp.Body)
	sim := strsim.Compare(string(srcBody), string(newBody))
	if sim > 0.9 {
		return SUCCESS
	}

	return FAIL
}
