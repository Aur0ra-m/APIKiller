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
