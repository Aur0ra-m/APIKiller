package origin

import (
	"net/http"
)

type TransferItem struct {
	Req  *http.Request
	Resp *http.Response
}

var transferItemQueue = make(chan *TransferItem)

type Origin interface {
	//
	// LoadOriginRequest
	//  @Description: load request and transport via channel transferItemQueue
	//
	LoadOriginRequest()
}

//
// TransportOriginRequest
//  @Description: transport requests from origin via transferItemQueue
//  @param item
//
func TransportOriginRequest(item *TransferItem) {
	transferItemQueue <- item
}

//
// GetOriginRequest
//  @Description: get requests from transferItemQueue
//
func GetOriginRequest() *TransferItem {
	return <-transferItemQueue
}
