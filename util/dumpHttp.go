package util

import (
	logger "APIKiller/log"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func DumpRequest(request *http.Request) string {
	if request == nil {
		logger.Debugln("dump response error: request is nil")
		return ""
	}

	dumpRequest, _ := httputil.DumpRequest(request, request.Body != nil)

	return string(dumpRequest)
}

func DumpResponse(response *http.Response) string {
	if response == nil {
		logger.Debugln("dump response error: response is nil")
		return ""
	}

	//all, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	logger.Errorln("read data from response body error", err)
	//	panic(err)
	//}
	//responseBody := bool(string(all) == "")

	dumpRequest, _ := httputil.DumpResponse(response, response.Body != nil)

	return string(dumpRequest)
}

func DumpRequests(requests []*http.Request) string {
	result := ""
	for _, request := range requests {
		dumpRequest := DumpRequest(request)

		result += fmt.Sprintf("\n**********************************************\n%s", dumpRequest)
	}

	return result
}

func DumpResponses(responses []*http.Response) string {
	result := ""

	for _, response := range responses {

		dumpResponse := DumpResponse(response)

		result += fmt.Sprintf("\n*****************************************\n%s", dumpResponse)
	}

	return result
}
