package ahttp

import (
	"APIKiller/core/ahttp/hook"
	"APIKiller/core/aio"
	logger "APIKiller/logger"
	"bufio"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
)

// DoRequest
//
//	@Description: make a http request without auto 30x redirect
//	@param r
//	@return *http.Response
func DoRequest(r *http.Request) *http.Response {
	var Client http.Client

	//fmt.Println(r.URL.String())

	// https request
	if r.URL.Scheme == "https" {
		// ignore certificate verification
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// https client
		Client = http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: tr,
		}
	} else {
		// http client
		Client = http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	// hook before initiating http request
	for _, requestHook := range hook.Hooks {
		requestHook.HookBefore(r)
	}

	response, err := Client.Do(r)
	if err != nil {
		logger.Errorln(err)
		return nil
	}

	// hook after finishing http request
	for _, requestHook := range hook.Hooks {
		requestHook.HookAfter(r)
	}

	// transform aio.Reader
	if response.Body != nil {
		response.Body = aio.TransformReadCloser(response.Body)
	}

	return response
}

func RequestClone(src *http.Request) *http.Request {
	// dump request
	reqStr := DumpRequest(src)
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
	if request.Body != nil {
		request.Body = aio.TransformReadCloser(request.Body)
	}

	return request
}

func ResponseClone(src *http.Response, req *http.Request) (dst *http.Response) {
	// dump response
	respStr := DumpResponse(src)

	// http.ReadResponse
	response, err := http.ReadResponse(bufio.NewReader(strings.NewReader(respStr)), req)
	if err != nil {
		logger.Errorln("read response error: ", err)
	}

	// transform body
	response.Body = aio.TransformReadCloser(response.Body)

	return response
}

//
// ExistsParam
//  @Description:
//  @param req
//  @param paramName
//  @return bool
//
func ExistsParam(req *http.Request, paramName string) bool {
	return false
}
