package origin

import (
	"net/http"
)

type TransferItem struct {
	Req   *http.Request
	Resp  *http.Response
	Https bool
}

var TransferItemQueue = make(chan *TransferItem)

type Origin interface {
	LoadOriginRequest()
}
