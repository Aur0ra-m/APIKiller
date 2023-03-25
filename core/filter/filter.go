package filter

import (
	"context"
	"net/http"
)

const (
	FilterPass    = true
	FilterBlocked = false
)

type Filter interface {
	Filter(context.Context, *http.Request) bool
}

func RegisterFilter(filters []Filter, filter Filter) []Filter {
	if filter == nil {
		return filters
	}

	return append(filters, filter)
}
