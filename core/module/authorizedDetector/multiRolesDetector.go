package authorizedDetector

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"strings"
)

type multiRolesDetector struct {
	authHeader map[string]string
}

func (d *multiRolesDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] multiple roles detect")

	req := item.SourceRequest

	// save old value and replace header
	oldValue := req.Header.Get(d.authHeader["httpHeader"])
	req.Header.Set(d.authHeader["httpHeader"], d.authHeader["value"])

	// do request
	response := http2.DoRequest(req, item.Https)

	// judge
	if Judge(ctx, item.SourceResponse, response) == Bypass {
		// deep copy
		vulnRequest := http2.RequestClone(req)

		item.VulnType = append(item.VulnType, "authorized-multiRoles")
		item.VulnRequest = append(item.VulnRequest, vulnRequest)
		item.VulnResponse = append(item.VulnResponse, response)
	}

	// recover old value
	req.Header.Set(d.authHeader["httpHeader"], oldValue)
}

func newMultiRolesDetector(ctx context.Context) *multiRolesDetector {
	logger.Infoln("[Load Module] multiple roles module")

	// get config: app.modules.authorizedDetector.multiRolesDetector.role
	role := util.GetConfig(ctx, "app.modules.authorizedDetector.multiRolesDetector.role")
	// split role into header and value
	splits := strings.Split(role, ":")
	header := splits[0]
	value := splits[1]

	m := map[string]string{}
	m["httpHeader"] = header
	m["value"] = value

	detector := &multiRolesDetector{m}
	return detector
}
