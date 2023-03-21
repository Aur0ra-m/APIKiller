package filter

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"net/http"
	"strings"
)

type StaticResourceFilter struct {
	forbidenExts []string
}

func (f *StaticResourceFilter) Filter(req *http.Request) bool {
	logger.Debug("[Filter] static file filter\n")

	// get request path extension
	lastIndex := strings.LastIndex(req.URL.Path, ".")
	if lastIndex < 0 {
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

func NewStaticFileFilter(cfg *config.Config) Filter {
	staticCfg := cfg.Filter.StaticFile
	logger.Debug("[Load Filter] static file filter\n")

	return &StaticResourceFilter{
		forbidenExts: staticCfg.Ext,
	}
}
