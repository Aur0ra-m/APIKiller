package ahttp

import (
	"APIKiller/internal/core/aio"
	"APIKiller/pkg/logger"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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

//
// ModifyParameterFormat
//  @Description: modify http parameter format. E.g., switch json to www-urlencoded
//  @param req
//  @return *http.Request
//
func ModifyParameterFormat(req *http.Request) *http.Request {
	return nil
}

//
// ModifyURLPathCase
//  @Description: modify the first letter of path to lowercase and the rest are uppercase. E.g. "/xxxx-->/xXXX"
//  @param Url
//
func ModifyURLPathCase(Url *url.URL) {
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
// ModifyURLPathMidPad
//  @Description: modify the end of url path with padding
//  @param Url
//  @param padding
//  @return *url.URL
//
func ModifyURLPathMidPad(Url *url.URL, padding string) {
	srcPath := Url.Path

	// trim the last separator
	if strings.HasSuffix(Url.Path, "/") {
		srcPath = strings.TrimRight(srcPath, "/")
	}

	// trim the first separator
	srcPath = strings.TrimLeft(srcPath, "/")

	// split path into fragments
	splits := strings.Split(srcPath, "/")

	firstPart := splits[0]
	thirdPart := ""
	newPath := ""

	if len(splits) > 1 {
		thirdPart = strings.Join(splits[1:], "/")

		newPath = firstPart + "/" + padding + "/" + thirdPart
	} else {
		newPath = "/" + padding + "/" + firstPart
	}
	Url.Path = newPath
}

//
// ModifyURLPathEndPad
//  @Description: modify the end of url path with padding
//  @param Url
//  @param padding
//  @return *url.URL
//
func ModifyURLPathEndPad(Url *url.URL, padding string) {
	srcPath := Url.Path

	// trim the last separator
	if strings.HasSuffix(Url.Path, "/") {
		srcPath = strings.TrimRight(srcPath, "/")
	}

	Url.Path = srcPath + padding
}

//
// ModifyURLPathAPIVerion
//  @Description: replace API version in path with target string
//  @param Url
//  @param srcString
//  @param targetString
//
func ModifyURLPathAPIVerion(Url *url.URL, srcString, targetString string) {
	Url.Path = strings.Replace(Url.Path, srcString, targetString, 1)
}

//
// ModifyParam
//  @Description: modify all positions parameter
//  @param req
//  @param paramName
//  @param newValue
//  @return bool
//
func ModifyParam(req *http.Request, paramName string, newValue string) *http.Request {

	newReq := ModifyQueryParam(req, paramName, newValue)

	if newReq != nil {
		return newReq
	}

	newReq = ModifyPostParam(req, paramName, newValue)

	if newReq != nil {
		return newReq
	}

	return nil
}

//
// ModifyQueryParam
//  @Description: modify url query parameter
//  @param req
//  @param paramName
//  @param newValue
//  @return *http.Request
//
func ModifyQueryParam(req *http.Request, paramName string, newValue string) *http.Request {
	queryParams := req.URL.Query()
	if queryParams.Get(paramName) != "" {
		
		// clone a new newReq
		newReq := RequestClone(req)

		newReq.URL.RawQuery = strings.Replace(newReq.URL.RawQuery, paramName+"="+newReq.URL.Query()[paramName][0], paramName+"="+newValue, 1)
		return newReq
	}

	return nil
}

//
// ModifyPostParam
//  @Description: modify post body parameter through Content-Type in request
//  @param req
//  @param paramName
//  @param newValue
//  @return *http.Request
//
func ModifyPostParam(req *http.Request, paramName string, newValue string) *http.Request {
	ct := req.Header.Get("Content-Type")

	// determine whether the request does have a body or not
	if ct == "" {
		logger.Debugln("target request does not have Content-Type header")
		return nil
	}

	readAll, _ := ioutil.ReadAll(req.Body)
	bodyStr := string(readAll)

	var newReq *http.Request
	if ct == "application/x-www-form-urlencoded" {
		if strings.Contains(bodyStr, paramName+"=") {
			newReq = RequestClone(req)

			modifyPostFormParam(newReq, paramName, newValue)
		}
	} else if ct == "application/json" {
		if strings.Contains(bodyStr, "\""+paramName+"\"") {
			newReq = RequestClone(req)

			modifyPostJsonParam(newReq, paramName, newValue)
		}
	} else if ct == "application/xml" {
		if strings.Contains(bodyStr, paramName+">") {
			newReq = RequestClone(req)

			modifyPostXMLParam(newReq, paramName, newValue)
		}
	} else if strings.Contains(ct, "multipart/form-data") {
		if strings.Contains(bodyStr, ";name=\""+paramName) {
			newReq = RequestClone(req)

			modifyPostMultiDataParam(newReq, paramName, newValue)
		}
	} else {
		logger.Errorln("Not support other Content-Type")
		return nil
	}

	return newReq
}

//
// modifyPostXMLParam
//  @Description: modify xml data with newValue
//  @param req
//  @param name
//  @param value
//
func modifyPostXMLParam(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := `<` + paramName + `>(.*?)<`

	modifyPostBody(req, paramItemRegExp, ">", newValue)
}

//
// modifyPostJsonParam
//  @Description: modify simple post json data with newValue
//  @param req
//  @param paramName
//  @param newValue
//
func modifyPostJsonParam(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := `"` + paramName + `"\s*?:\s*"?(.*?)"?[\s,\}]`

	modifyPostBody(req, paramItemRegExp, ":", newValue)
}

//{"size": "10000"
//
// modifyPostFormParam
//  @Description: modify simple post form data with newValue
//  @param req
//  @param paramName
//  @param newValue
//
func modifyPostFormParam(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := paramName + `=([^&]*)`

	modifyPostBody(req, paramItemRegExp, "=", newValue)
}

//
// modifyPostMultiDataParam
//  @Description: modify post parameter in form of multi-data
//  @param req
//  @param paramName
//  @param newValue
//
func modifyPostMultiDataParam(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := "name=\"" + paramName + "\".*?" + "\r\n\r\n" + "([^\r\n]*)"

	modifyPostBody(req, paramItemRegExp, "\r\n\r\n", newValue)
}

//
// modifyPostBody
//  @Description: modify parameter item in post body
//  @param req
//  @param paramItemRegExp
//  @param paramKVSeparator the separator between parameter key and parameter value
//  @param newValue
//
func modifyPostBody(req *http.Request, paramItemRegExp string, paramKVSeparator string, newValue string) {
	// read data from body
	all, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorln(fmt.Sprintf("modify post data error:%v", err))
		panic(err)
	}
	body := string(all)

	// extract paramItem and value that will be replaced from body with regular expression
	r := regexp.MustCompile(paramItemRegExp)
	submatchs := r.FindStringSubmatch(body)
	paramItem := submatchs[0]
	srcValue := submatchs[1]

	if srcValue == "" {
		logger.Infoln("original parameter has no value")
		return
	}

	// split paramItem into two parts, key-part and value-part
	splits := strings.Split(paramItem, paramKVSeparator)

	// replace source value with new value
	newValuePart := strings.Replace(splits[1], srcValue, newValue, 1)

	// join key-part and value-part
	newParamItem := splits[0] + paramKVSeparator + newValuePart

	// replace paramItem in body
	newBody := strings.Replace(body, paramItem, newParamItem, 1)

	// refill body
	req.Body = aio.TransformReadCloser(bytes.NewReader([]byte(newBody)))

	// update Content-Length
	req.ContentLength = int64(len(newBody))
}

//
// ModifyParamByRegExp
//  @Description: modify parameter value everywhere through specified value format
//  @param req
//  @param paramName
//  @param newValue
//  @return *http.Request
//
func ModifyParamByRegExp(req *http.Request, paramName string, newValue string) *http.Request {

	newReq := ModifyQueryParam(req, paramName, newValue)

	if newReq != nil {
		return newReq
	}

	newReq = ModifyPostParam(req, paramName, newValue)

	if newReq != nil {
		return newReq
	}

	return nil
}

//
// ModifyQueryParamByRegExp
//  @Description: modify query parameter value  through specified value format
//  @param req
//  @param valueRegExp
//  @param newValue
//  @return *http.Request
//
func ModifyQueryParamByRegExp(req *http.Request, valueRegExp string, newValue string) *http.Request {
	re := regexp.MustCompile(valueRegExp)

	findAllString := re.FindAllString(req.URL.RawQuery, -1)
	if len(findAllString) > 0 {
		// clone request
		newReq := RequestClone(req)

		newReq.URL.RawQuery = strings.Replace(req.URL.RawQuery, findAllString[0], newValue, 1)
		return newReq
	}

	return nil
}

//
// ModifyPostParamByRegExp
//  @Description: modify parameter value in post body through specified value format
//  @param req
//  @param valueRegExp value format
//  @param newValue
//  @return *http.Request
//
func ModifyPostParamByRegExp(req *http.Request, valueRegExp string, newValue string) *http.Request {

	re := regexp.MustCompile(valueRegExp)

	findAllString := re.FindAllString(req.URL.RawQuery, -1)
	if len(findAllString) > 0 {
		// clone request
		newReq := RequestClone(req)

		newReq.URL.RawQuery = strings.Replace(req.URL.RawQuery, findAllString[0], newValue, 1)
		return newReq
	}

	ct := req.Header.Get("Content-Type")

	// determine whether the request does have a body or not
	if ct == "" {
		logger.Debugln("target request does not have Content-Type header")
		return nil
	}

	readAll, _ := ioutil.ReadAll(req.Body)
	bodyStr := string(readAll)

	matches := re.FindAllString(bodyStr, -1)

	if len(matches) <= 0 {
		return nil
	}

	// clone request
	newReq := RequestClone(req)

	// replace parameter matching the value format
	newBody := strings.Replace(bodyStr, matches[0], newValue, 1)

	// refill body
	req.Body = aio.TransformReadCloser(bytes.NewReader([]byte(newBody)))

	// update Content-Length
	req.ContentLength = int64(int(len(newBody)))

	return newReq
}

//
// AppendHeader
//  @Description: append new header to request
//  @param req
//  @param header
//  @param value
//
func AppendHeader(req *http.Request, header string, value string) {
	req.Header.Add(header, value)
}

//
// RemoveHeader
//  @Description: remove the specified header from the request
//  @param req
//  @param header
//
func RemoveHeader(req *http.Request, header string) {
	if req.Header.Get(header) != "" {
		req.Header.Del(header)
	} else {
		logger.Errorln("no specified header which will be removed in request")
	}
}

//
// UpdateHeader
//  @Description: update the specified header in the request
//  @param req
//  @param header
//  @param value
//
func UpdateHeader(req *http.Request, header string, value string) {
	if req.Header.Get(header) != "" {
		req.Header.Set(header, value)
	} else {
		logger.Errorln("no specified header which will be updated in request")
	}
}
