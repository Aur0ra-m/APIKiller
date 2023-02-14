package authorizedDetector

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"encoding/json"
)

type unauthorizedDetector struct {
	commonAuthHeaders []string
	normalHeaders     []string
}

func (d *unauthorizedDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] unauthorized detect")

	req := item.SourceRequest

	// default flag
	f := false

	// enum all headers
	for key, header := range req.Header {
		// common auth headers match
		for _, authHeader := range d.commonAuthHeaders {
			if key == authHeader {
				// save old value
				oldvalue, _ := json.Marshal(header)

				if d.detect(ctx, item, key) {
					// update f
					f = true

					// recover old value
					req.Header.Add(key, string(oldvalue))

					break
				}

				// recover old value
				req.Header.Add(key, string(oldvalue))
				continue
			}
		}

		//// normal headers match
		//isNormal := false
		//for _, normalHeader := range d.normalHeaders {
		//	if key == normalHeader {
		//		isNormal = true
		//		break
		//	}
		//}
		//if isNormal {
		//	continue
		//}

		//// other headers detect
		//
		//// save old value
		//oldvalue, _ := json.Marshal(header)
		//
		//if d.detect(ctx, item, key) {
		//	// recover old value
		//	req.Header.Add(key, string(oldvalue))
		//
		//	// add header to auth headers
		//	d.commonAuthHeaders = append(d.commonAuthHeaders, key)
		//
		//	break
		//}
		//
		//// recover old value
		//req.Header.Add(key, string(oldvalue))

	}

	// if f is true ,there will be an unauthorized api
	if f {
		item.VulnType = append(item.VulnType, "unauthorized")
	}

}

//
// detect
//  @Description: detect whether target header is effective authority header
//  @receiver d
//  @param ctx
//  @param item
//  @param targetHeader
//  @return bool
//
func (d *unauthorizedDetector) detect(ctx context.Context, item *data.DataItem, targetHeader string) bool {
	req := item.SourceRequest

	// delete target header
	req.Header.Del(targetHeader)

	// do request
	response := http2.DoRequest(req)

	// judge
	if Judge(ctx, item.SourceResponse, response) == Bypass {
		return true
	}
	return false
}

func newUnauthorizedDetector(ctx context.Context) *unauthorizedDetector {
	if util.GetConfig(ctx, "app.detectors.authorizedDetector.unauthorizedDetector.commonAuthHeaders") == "0" {
		return nil
	}
	logger.Infoln("[Load Module] unauthorized module")

	detector := &unauthorizedDetector{}

	commonAuthHeaders := util.GetConfig(ctx, "app.detectors.authorizedDetector.unauthorizedDetector.commonAuthHeaders")
	normalHeaders := util.GetConfig(ctx, "app.detectors.authorizedDetector.unauthorizedDetector.normalHeaders")

	detector.commonAuthHeaders = util.SplitConfigString(commonAuthHeaders)
	detector.normalHeaders = util.SplitConfigString(normalHeaders)

	return detector
}
