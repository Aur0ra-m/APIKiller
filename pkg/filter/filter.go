package filter

import (
	"APIKiller/pkg/config"
	"net/http"
)

const (
	FilterPass    = true
	FilterBlocked = false
)

type Filter interface {
	Filter(*http.Request) bool
}

func NewFilter(cfg *config.Config) []Filter {
	var filters []Filter

	filters = append(filters, NewHttpFilter(cfg))
	filters = append(filters, NewStaticFileFilter(cfg))
	filters = append(filters, NewDuplicateFilter(cfg))
	return filters
}
