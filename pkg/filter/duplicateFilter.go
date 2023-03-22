package filter

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"
)

type DuplicateFilter struct {
	history []string // []string{"GET domain /admin/index",}
}

func (f *DuplicateFilter) Filter(req *http.Request) bool {
	logger.Debug("[Filter] duplicate\n")

	// format
	curr := fmt.Sprintf("%s %s %s", req.Method, req.Host, req.URL.Path)

	// duplication
	if slices.Contains(f.history, curr) {
		logger.Info("duplicate data\n")
		return FilterBlocked
	}

	// append to history
	f.history = append(f.history, curr)

	return FilterPass
}

func NewDuplicateFilter(cfg *config.FilterConfig) *DuplicateFilter {
	logger.Info("[Load Filter] duplicate filter\n")
	return &DuplicateFilter{}
}
