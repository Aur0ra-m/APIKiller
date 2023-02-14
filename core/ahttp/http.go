package ahttp

import (
	"APIKiller/core/aio"
	logger "APIKiller/log"
	"APIKiller/util"
	"bufio"
	"net/http"
	"net/url"
	"strings"
)

//
// DoRequest
//  @Description: make an http request, and transform body before return response
//  @param r
//  @return *http.Response
//
func DoRequest(r *http.Request) *http.Response {
	Client := &http.Client{}
	response, err := Client.Do(r)
	if err != nil {
		logger.Errorln(err)
		return nil
	}

	// transform aio.Reader
	if response.Body != nil {
		response.Body = aio.TransformReadCloser(response.Body)
	}

	return response
}

func RequestClone(src *http.Request) *http.Request {
	// dump request
	reqStr := util.DumpRequest(src)
	// http.ReadRequest
	request, err := http.ReadRequest(bufio.NewReader(strings.NewReader(reqStr)))
	if err != nil {
		logger.Errorln("read request error: ", err)
	}
	// we can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	// set url
	u, err := url.Parse(src.URL.String())
	if err != nil {
		logger.Errorln("parse url error: ", err)
	}
	request.URL = u
	// transform body
	request.Body = aio.TransformReadCloser(request.Body)

	return request
}

func ResponseClone(src *http.Response, req *http.Request) (dst *http.Response) {
	// dump response
	respStr := util.DumpResponse(src)

	// http.ReadResponse
	response, err := http.ReadResponse(bufio.NewReader(strings.NewReader(respStr)), req)
	if err != nil {
		logger.Errorln("read response error: ", err)
	}

	// transform body
	response.Body = aio.TransformReadCloser(response.Body)

	return response
}
