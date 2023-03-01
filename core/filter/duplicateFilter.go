package filter

import (
	logger "APIKiller/log"
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"
)

type DuplicateFilter struct {
	history []string // []string{"GET domain /admin/index",}
}

func (f *DuplicateFilter) Filter(ctx context.Context, req *http.Request) bool {
	logger.Debugln("[Filter] duplicate")

	// format
	curr := fmt.Sprintf("%s %s %s", req.Method, req.Host, req.URL.Path)

	// duplication
	if slices.Contains(f.history, curr) {
		logger.Infoln("duplicate data")
		return FilterBlocked
	}

	// append to history
	f.history = append(f.history, curr)

	return FilterPass
}

func NewDuplicateFilter() *DuplicateFilter {
	logger.Infoln("[Load Filter] duplicate filter")
	return &DuplicateFilter{}
}
