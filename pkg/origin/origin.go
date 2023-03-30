package origin

import (
	"APIKiller/pkg/config"
	"net/http"
)

type TransferItem struct {
	Req   *http.Request
	Resp  *http.Response
	Https bool
}

type Origin interface {
	LoadOriginRequest(cfg *config.OriginConfig, httpItemQueue chan *TransferItem)
}
