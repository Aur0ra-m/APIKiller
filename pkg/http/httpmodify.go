package http

import (
	"APIKiller/pkg/logger"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// ModifyMethod
//
//	@Description: modify to target method
//	@param method
//	@param req
//	@return *http.Request
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
		request.Body = TransformReadCloser(request.Body)
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

// ModifyURLPathCase
//
//	@Description: modify the first letter of path to lowercase and the rest are uppercase. E.g. "/xxxx-->/xXXX"
//	@param url
//	@return *url.URL
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

// ModifyURLPathMidPad
//
//	@Description: modify the end of url path with padding
//	@param Url
//	@param padding
//	@return *url.URL
func ModifyURLPathMidPad(Url *url.URL, padding string) {
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

// ModifyURLPathEndPad
//
//	@Description: modify the end of url path with padding
//	@param Url
//	@param padding
//	@return *url.URL
func ModifyURLPathEndPad(Url *url.URL, padding string) {
	srcPath := Url.Path

	// trim the last separator
	if strings.HasSuffix(Url.Path, "/") {
		srcPath = strings.TrimRight(srcPath, "/")
	}

	Url.Path = srcPath + padding
}

// ModifyURLPathAPIVerion
//
//	@Description: replace API version in path with target string
//	@param Url
//	@param srcString
//	@param targetString
func ModifyURLPathAPIVerion(Url *url.URL, srcString, targetString string) {
	Url.Path = strings.Replace(Url.Path, srcString, targetString, 1)
}

// ModifyURLQueryParameter
//
//	@Description: modify parameters in form of url query
//	@param req
//	@param paramName
//	@param newValue
func ModifyURLQueryParameter(req *http.Request, paramName string, newValue []string) {
	req.URL.RawQuery = strings.Replace(req.URL.RawQuery, paramName+"="+req.URL.Query()[paramName][0], paramName+"="+newValue[0], 1)
}

//
// ModifyPostFormParameter
//  @Description: modify simple post form data with newValue
//  @param req
//  @param paramName
//  @param newValue
//
//func ModifyPostFormParameter(req *http.Request, paramName string, newValue string) {
//	// read data from body
//	all, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		logger.Errorln("modify post data error:%v", err)
//		panic(err)
//	}
//	body := string(all)
//
//	// split data with &
//	splits := strings.Split(body, "&")
//
//	// replace data and join again
//	for i, _ := range splits {
//		kv := strings.Split(splits[i], "=")
//		if kv[0] == paramName {
//			splits[i] = kv[0] + "=" + newValue
//			break
//		}
//	}
//
//	newBody := strings.Join(splits, "&")
//
//	// back-fill body
//	req.Body = aio.TransformReadCloser(bytes.NewReader([]byte(newBody)))
//
//	// update Content-Length
//	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(newBody)))
//}

// ModifyPostJsonParameter
//
//	@Description: modify simple post json data with newValue
//	@param req
//	@param paramName
//	@param newValue
func ModifyPostJsonParameter(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := `"` + paramName + `"\s*?:\s*?"?(.*?)?"?,?\s`

	modifyPostBody(req, paramItemRegExp, ":", newValue)
}

// ModifyPostFormParameter
//
//	@Description: modify simple post form data with newValue
//	@param req
//	@param paramName
//	@param newValue
func ModifyPostFormParameter(req *http.Request, paramName string, newValue string) {
	paramItemRegExp := paramName + `=(.*?)&?\s?`

	modifyPostBody(req, paramItemRegExp, "=", newValue)
}

// ModifyPostMultiDataParameter
//
//	@Description: modify post parameter in form of multi-data
//	@param req
//	@param paramName
//	@param newValue
func ModifyPostMultiDataParameter(req *http.Request, paramName string, newValue string) {
	//
}

// modifyPostBody
//
//	@Description: modify parameter item in post body
//	@param req
//	@param paramItemRegExp
//	@param paramKVSeparator the separator between parameter key and parameter value
//	@param newValue
func modifyPostBody(req *http.Request, paramItemRegExp string, paramKVSeparator string, newValue string) {
	// read data from body
	all, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf(fmt.Sprintf("modify post data error:%v", err))
		panic(err)
	}
	body := string(all)

	// extract paramItem and value that will be replaced from body with regular expression
	r := regexp.MustCompile(paramItemRegExp)
	submatchs := r.FindStringSubmatch(body)
	paramItem := submatchs[0]
	srcValue := submatchs[1]

	// split paramItem into two parts, key-part and value-part
	splits := strings.Split(paramItem, paramKVSeparator)

	// replace source value with new value
	newValuePart := strings.Replace(splits[1], srcValue, newValue, 1)

	// join key-part and value-part
	newParamItem := splits[0] + paramKVSeparator + newValuePart

	// replace paramItem in body
	newBody := strings.Replace(body, paramItem, newParamItem, 1)

	// refill body
	req.Body = TransformReadCloser(bytes.NewReader([]byte(newBody)))

	// update Content-Length
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(newBody)))
}
