package authorizedDetector

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/log"
	"context"
	"github.com/spf13/viper"
)

type AuthorizedDetector struct {
	authHeader       string
	Roles            []string
	blackStatusCodes []int
	blackKeywords    []string
}

func (d *AuthorizedDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] authorized detect")

	d.unauthorizedDetect(ctx, item)
	for _, t := range item.VulnType {
		if t == "unauthorized" {
			return
		}
	}

	d.multiRolesDetect(ctx, item)

}

// unauthorizedDetect
//
//	@Description: unauthorized header detect
//	@receiver d
//	@param ctx
//	@param item
func (d *AuthorizedDetector) unauthorizedDetect(ctx context.Context, item *data.DataItem) {
	req := http2.RequestClone(item.SourceRequest)

	// delete auth header
	req.Header.Del(d.authHeader)

	// make request and judge
	response := http2.DoRequest(req)

	if d.Judge(ctx, item.SourceResponse, response) == Bypass {
		item.VulnType = append(item.VulnType, "unauthorized")
		item.VulnRequest = append(item.VulnRequest, req)
		item.VulnResponse = append(item.VulnResponse, response)
	}
}

// multiRolesDetect
//
//	@Description: multiple roles detect
//	@receiver d
//	@param ctx
//	@param item
func (d *AuthorizedDetector) multiRolesDetect(ctx context.Context, item *data.DataItem) {
	req := http2.RequestClone(item.SourceRequest)

	// default support one role
	newRole := d.Roles[0]

	// change auth header
	req.Header.Set(d.authHeader, newRole)

	// do request
	response := http2.DoRequest(req)

	// judge
	if d.Judge(ctx, item.SourceResponse, response) == Bypass {
		item.VulnType = append(item.VulnType, "authorized-multiRoles")
		item.VulnRequest = append(item.VulnRequest, req)
		item.VulnResponse = append(item.VulnResponse, response)
	}
}

func NewAuthorizedDetector(ctx context.Context) module.Detecter {
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
	}
}
