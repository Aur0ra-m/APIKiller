package origin

import (
	"net/http"
)

type TransferItem struct {
	Req  *http.Request
	Resp *http.Response
}

var TransferItemQueue = make(chan *TransferItem)

type Origin interface {
	LoadOriginRequest()
}
