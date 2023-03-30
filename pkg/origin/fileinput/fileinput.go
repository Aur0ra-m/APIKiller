package fileinput

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/origin"
	"bufio"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type FileInput struct {
	path string
}

func (f *FileInput) LoadOriginRequest(cfg *config.OriginConfig, httpItemQueue chan *origin.TransferItem) {
	logger.Info("[Load Request] load request from file input origin")
	if stat, _ := os.Stat(f.path); stat.IsDir() {

	} else {
		// load origin from target file[eg. burp file]
		f.parseDataFromBurpFile(httpItemQueue)
	}
}

// RecoverHttpRequest
//
//	@Description: create one new http.Request with rawRequest and rawURL
//	@param rawRequest
//	@param rawURL
//	@return *http.Request
func RecoverHttpRequest(rawRequest, rawURL, rawResponse string) (*http.Request, *http.Response) {
	reqWrapper := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := http.ReadRequest(reqWrapper)
	if err != nil {
		panic(err)
	}

	req.RequestURI = ""
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	req.URL = u

	respWrapper := bufio.NewReader(strings.NewReader(rawResponse))
	resp, _ := http.ReadResponse(respWrapper, req)
	return req, resp
}

func NewFileInputOrigin(path string) *FileInput {
	logger.Info("[Origin] file input origin")
	// determine whether the path is a file or a directory
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Errorf("%s does not exist", err)
		panic(err)
	}

	// instantiate FileInputOrigin
	return &FileInput{path: path}
}
