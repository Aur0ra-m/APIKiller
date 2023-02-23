package filter

import (
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"net/http"
	"regexp"
)

type HttpFilter struct {
}

func (f HttpFilter) Filter(ctx context.Context, req *http.Request) bool {
	logger.Debugln("[Filter] ahttp filter")

	// get config and match through RegExp
	hostExp := util.GetConfig(ctx, "app.filters.httpFilter.host")
	hostsExp := util.SplitConfigString(hostExp)
	if len(hostsExp) != 0 {
		reqHost := req.Host
		flag := FilterBlocked
		for _, hostExp := range hostsExp {
			if matched, _ := regexp.Match(hostExp, []byte(reqHost)); matched {
				flag = FilterPass
				break
			}
		}
		return flag
	}

	return FilterPass // default
}

func NewHttpFilter() Filter {
	logger.Infoln("[Load Filter] http filter")

	return &HttpFilter{}
}
