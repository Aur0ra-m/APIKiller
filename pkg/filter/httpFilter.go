package filter

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"net/http"
	"regexp"
)

type HttpFilter struct {
	hostsExp []string
}

func (f *HttpFilter) Filter(req *http.Request) bool {
	logger.Debug("[Filter] http filter\n")

	// match through RegExp
	if len(f.hostsExp) != 0 {
		reqHost := req.Host
		flag := FilterBlocked
		for _, hostExp := range f.hostsExp {
			if matched, _ := regexp.Match(hostExp, []byte(reqHost)); matched {
				flag = FilterPass
				break
			}
		}
		return flag
	}

	return FilterPass // default
}

func NewHttpFilter(cfg *config.FilterConfig) Filter {
	httpCfg := cfg.Http
	logger.Info("[Load Filter] http filter\n")

	return &HttpFilter{
		hostsExp: httpCfg.Host,
	}
}
