package authorize

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/detector"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
)

type AuthorizedDetector struct {
	authHeader       string
	Roles            []string
	blackStatusCodes []int
	blackKeywords    []string
}

func (d *AuthorizedDetector) Detect(item *types.DataItem) {

}

// unauthorizedDetect
//
//	@Description: unauthorized header detect
//	@receiver d
//	@param ctx
//	@param item
func (d *AuthorizedDetector) unauthorizedDetect(item *types.DataItem) {
	req := gohttp.RequestClone(item.SourceRequest)

	// delete auth header
	req.Header.Del(d.authHeader)

	// make request and judge
	response := gohttp.DoRequest(req)

	if d.Judge(item.SourceResponse, response) {
		item.VulnType = append(item.VulnType, "unauthorized")
		item.VulnRequest = append(item.VulnRequest, req)
		item.VulnResponse = append(item.VulnResponse, response)
	}
}

func NewUnauthorizedDetector(cfg *config.Config) detector.Detector {
	authConf := cfg.Detector.Authorize
	if !authConf.Enable {
		return nil
	}
	if len(authConf.Roles) == 0 {
		logger.Errorf("authorizedDetector did not set roles")
		return nil
	}

	return &AuthorizedDetector{
		authHeader:       authConf.AuthHeader,
		Roles:            authConf.Roles,
		blackStatusCodes: authConf.Judgement["blackStatusCodes"].([]int),
		blackKeywords:    authConf.Judgement["blackKeywords"].([]string),
	}
}
