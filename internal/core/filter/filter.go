package filter

import (
	"net/http"
)

const (
	FilterPass    = true
	FilterBlocked = false
)

var (
	filters []Filter
)

type Filter interface {
	//
	// Filter
	//  @Description: filter out *http.Request that do not meet the conditions
	//  @param *http.Request
	//  @return bool
	//
	Filter(*http.Request) bool
}

func RegisterFilter(filter Filter) {
	if filter == nil {
		return
	}

	filters = append(filters, filter)
}

func GetFilters() []Filter {
	if filters != nil {
		return filters
	}

	return nil
}
