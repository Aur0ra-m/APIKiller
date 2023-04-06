package filter

import (
	"net/http"
)

const (
	FilterPass    = true
	FilterBlocked = false
)

var (
	Filters []Filter
)

type Filter interface {
	Filter(*http.Request) bool
}

func RegisterFilter(filter Filter) {
	if filter == nil {
		return
	}

	Filters = append(Filters, filter)
}
