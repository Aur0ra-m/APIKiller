package filter

import (
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"net/http"
	"strings"
)

type StaticResourceFilter struct {
	forbidenExts []string
}

func (f *StaticResourceFilter) init(ctx context.Context) {
	// parse config
	extString := util.GetConfig(ctx, "app.filters.staticFileFilter.ext")
	f.forbidenExts = util.SplitConfigString(extString)
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

	f := &StaticResourceFilter{}

	f.init(ctx)

	return f
}
