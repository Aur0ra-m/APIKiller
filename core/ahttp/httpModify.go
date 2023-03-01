package ahttp

import (
	"APIKiller/core/aio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//
// ModifyMethod
//  @Description: modify to target method
//  @param method
//  @param req
//  @return *http.Request
//
func ModifyMethod(method string, req *http.Request) *http.Request {
	// get URL and read body
	url := req.URL.String()

	body := ""
	if req.Body != nil {
		bytes, _ := ioutil.ReadAll(req.Body)
		body = string(bytes)
	}

	var (
		request *http.Request
	)

	if body == "" {
		request, _ = http.NewRequest(method, url, nil)
	} else {
		request, _ = http.NewRequest(method, url, bytes.NewReader([]byte(body)))

		// transform body
		request.Body = aio.TransformReadCloser(request.Body)
	}

	// copy headers
	for s, _ := range req.Header {
		request.Header.Set(s, req.Header.Get(s))
	}

	return request
}

// ModifyParameterFormat
//
//	@Description: modify http parameter format. E.g., switch json to www-urlencoded
func ModifyParameterFormat(req *http.Request) *http.Request {
	return nil
}

// URLPathCase
//
//	@Description: modify the first letter of path to lowercase and the rest are uppercase. E.g. "/xxxx-->/xXXX"
//	@param url
//	@return *url.URL
func URLPathCase(Url *url.URL) {
	srcPath := Url.Path

	splits := strings.Split(srcPath, "/")

	// switch to upper case
	lastPart := splits[len(splits)-1]

	lastPart = strings.ToUpper(lastPart)
	lastPart = strings.ToLower(lastPart[0:1]) + lastPart[1:]

	splits[len(splits)-1] = lastPart

	// join them into a new path
	newPath := strings.Join(splits, "/")

	Url.Path = newPath
}

//
// URLPathMidPad
//  @Description: modify the end of url path with padding
//  @param Url
//  @param padding
//  @return *url.URL
//
func URLPathMidPad(Url *url.URL, padding string) {
	srcPath := Url.Path
	// trim the last separator
	if strings.HasSuffix(Url.Path, "/") {
		srcPath = strings.TrimRight(srcPath, "/")
	}

	// split path into fragments
	splits := strings.Split(srcPath, "/")

	firstPart := splits[0]
	thirdPart := ""

	newPath := ""

	if len(splits) > 1 {
		thirdPart = strings.Join(splits[1:], "/")

		newPath = firstPart + "/" + padding + "/" + thirdPart
	} else {
		newPath = padding + "/" + firstPart
	}
	Url.Path = newPath
}

//
// URLPathEndPad
//  @Description: modify the end of url path with padding
//  @param Url
//  @param padding
//  @return *url.URL
//
func URLPathEndPad(Url *url.URL, padding string) {
	srcPath := Url.Path

	// trim the last separator
	if strings.HasSuffix(Url.Path, "/") {
		srcPath = strings.TrimRight(srcPath, "/")
	}

	Url.Path = srcPath + padding
}

//
// URLPathAPIVerionModify
//  @Description: replace API version in path with target string
//  @param Url
//  @param srcString
//  @param targetString
//
func URLPathAPIVerionModify(Url *url.URL, srcString, targetString string) {
	Url.Path = strings.Replace(Url.Path, srcString, targetString, 1)
}
