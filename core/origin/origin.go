package origin

import (
	"net/http"
)

type TransferItem struct {
	Req   *http.Request
	Resp  *http.Response
	Https bool
}

type Origin interface {
	LoadOriginRequest(httpItemQueue chan *TransferItem)
}
