package filter

import (
	logger "APIKiller/log"
	"context"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type StaticResourceFilter struct {
	forbidenExts []string
}

func (f *StaticResourceFilter) Filter(ctx context.Context, req *http.Request) bool {
	logger.Debugln("[Filter] static file filter")

	// get request path extension
	lastIndex := strings.LastIndex(req.URL.Path, ".")
	if lastIndex == -1 {
		return FilterPass
	}
	ext := req.URL.Path[lastIndex+1:]

	// filter
	for _, forbidenExt := range f.forbidenExts {
		if forbidenExt == ext {
			return FilterBlocked
		}
	}

	return FilterPass
}

func NewStaticFileFilter(ctx context.Context) Filter {
	logger.Infoln("[Load Filter] static file filter")

	return &StaticResourceFilter{
		forbidenExts: viper.GetStringSlice("app.filter.staticFileFilter.ext"),
	}
}
