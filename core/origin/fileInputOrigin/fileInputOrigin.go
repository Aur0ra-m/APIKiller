package fileInputOrigin

import (
	"APIKiller/core/origin"
	logger "APIKiller/log"
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type FileInputOrigin struct {
	path string
}

func (o *FileInputOrigin) LoadOriginRequest(ctx context.Context, httpItemQueue chan *origin.TransferItem) {
	logger.Infoln("[Load Request] load request from file input origin")

	if stat, _ := os.Stat(o.path); stat.IsDir() {
		// load origin from target directory

	} else {
		// load origin from target file[eg. burp file]
		o.parseDataFromBurpFile(httpItemQueue)
	}

}

// RecoverHttpRequest
//
//	@Description: create one new http.Request with rawRequest and rawURL
//	@param rawRequest
//	@param rawURL
//	@return *http.Request
func RecoverHttpRequest(rawRequest, rawURL, rawResponse string) (*http.Request, *http.Response) {
	b := bufio.NewReader(strings.NewReader(rawRequest))

	req, err := http.ReadRequest(b)
	if err != nil {
		panic(err)
	}

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	req.RequestURI = ""

	// Since the req.URL will not have all the information set,
	// such as protocol scheme and host, we create a new URL
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	req.URL = u

	b2 := bufio.NewReader(strings.NewReader(rawResponse))

	resp, _ := http.ReadResponse(b2, req)

	return req, resp
}

func NewFileInputOrigin(path string) *FileInputOrigin {
	logger.Infoln("[Origin] file input origin")

	// determine whether the path is a file or a directory
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Errorln(fmt.Sprintf("%s does not exist", path))
		panic(fmt.Sprintf("%s does not exist", path))
	}

	// instantiate FileInputOrigin
	return &FileInputOrigin{path: path}
}
