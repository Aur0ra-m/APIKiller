package authorizedDetector

import (
	"APIKiller/core/data"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
)

type singleRoleDetector struct {
	pattern1 string // RESTFul API format:/post/:id and static file object format:/api/v1/huiqeb1230b32jv43h2b3k.jpg
	pattren2 string //
	role1    string
	role2    string
}

func (d *singleRoleDetector) init() {
	d.pattern1 = `/[0-9a-zA-Z]{10,25}`

	d.role1 = "\\d{5,25}"
	d.role2 = "[0-9a-zA-Z]{10,25}"
}

func (d *singleRoleDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect ] single role detect")
	//req := item.SourceRequest

	// handle restful and file object form
	//path := req.URL.path

	// query parameters

	// body parameters:json/url-encoded/
}

func newSingleRoleDetector(ctx context.Context) *singleRoleDetector {
	//whether to use the current module
	if util.GetConfig(ctx, "app.detectors.authorizedDetector.singleRoleDetector.option") != "1" {
		return nil
	}

	logger.Infoln("[Load Module] single role module")
	detector := &singleRoleDetector{}

	detector.init()

	return detector
}
