package filter

import (
	logger "APIKiller/log"
	"context"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
)

type HttpFilter struct {
	hostsExp []string
}

func (f *HttpFilter) Filter(ctx context.Context, req *http.Request) bool {
	logger.Debugln("[Filter] ahttp filter")

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

func NewHttpFilter() Filter {
	logger.Infoln("[Load Filter] http filter")

	return &HttpFilter{
		hostsExp: viper.GetStringSlice("app.filter.httpFilter.host"),
	}
}
