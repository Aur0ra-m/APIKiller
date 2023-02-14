package origin

import (
	"context"
	"net/http"
)

type HttpItem struct {
	Req   *http.Request
	Resp  *http.Response
	Https bool
}

type Origin interface {
	LoadOriginRequest(ctx context.Context, httpItemQueue chan *HttpItem)
}
