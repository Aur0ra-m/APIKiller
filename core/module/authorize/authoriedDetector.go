package authorize

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"APIKiller/util"
	"github.com/spf13/viper"
)

type AuthorizedDetector struct {
	authHeader       string
	Roles            []string
	blackStatusCodes []int
	blackKeywords    []string
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

	if d.judge(item.SourceResponse, newResp) == Bypass {
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
	if d.judge(item.SourceResponse, newResp) == Bypass {
		return util.BuildResult(item, "authorized-multiRoles", newReq, newResp)
	}

	return nil
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
	}
}
