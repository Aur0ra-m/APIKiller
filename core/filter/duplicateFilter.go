package filter

import (
	"APIKiller/core/database"
	logger "APIKiller/log"
	"context"
	"net/http"
)

type DuplicateFilter struct {
}

func (f *DuplicateFilter) Filter(ctx context.Context, req *http.Request) bool {
	logger.Debugln("[Filter] duplicate")

	// get domain,url,Method
	domain := req.Host
	url := req.URL.Path
	method := req.Method

	// duplication
	db := ctx.Value("db").(database.Database)
	if db.Exist(domain, url, method) {
		logger.Infoln("duplicate data")
		return FilterBlocked
	}
	return FilterPass
}

func NewDuplicateFilter() *DuplicateFilter {
	logger.Infoln("[Load Filter] duplicate filter")
	return &DuplicateFilter{}
}
