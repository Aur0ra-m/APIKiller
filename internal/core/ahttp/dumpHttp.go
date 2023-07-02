package ahttp

import (
	"APIKiller/pkg/logger"
	"net/http"
	"net/http/httputil"
)

//
// DumpRequest
//  @Description: convert *http.Request to string
//  @param request
//  @return string
//
func DumpRequest(request *http.Request) string {
	if request == nil {
		logger.Debugln("dump response error: request is nil")
		return ""
	}

	dumpRequest, _ := httputil.DumpRequest(request, request.Body != nil)

	return string(dumpRequest)
}

//
// DumpResponse
//  @Description: convert *http.Response to string
//  @param response
//  @return string
//
func DumpResponse(response *http.Response) string {
	if response == nil {
		logger.Debugln("dump response error: response is nil")
		return ""
	}

	dumpRequest, _ := httputil.DumpResponse(response, response.Body != nil)

	return string(dumpRequest)
}

//
// DumpRequests
//  @Description: convert []*http.Request to []string
//  @param requests
//  @return []string
//
func DumpRequests(requests []*http.Request) []string {
	var result []string
	for _, request := range requests {
		dumpRequest := DumpRequest(request)
		result = append(result, dumpRequest)
	}

	return result
}

//
// DumpResponses
//  @Description: convert []*http.Response to []string
//  @param requests
//  @return []string
//
func DumpResponses(responses []*http.Response) []string {
	var result []string
	for _, response := range responses {
		dumpResponse := DumpResponse(response)
		result = append(result, dumpResponse)
	}

	return result
}
