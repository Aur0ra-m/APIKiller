package fileinput

import (
	"APIKiller/pkg/logger"
	"APIKiller/pkg/origin"
	"encoding/base64"
	"github.com/beevik/etree"
)

// parseData
//
//	@Description: parse data from burpsuite file
//	@receiver o
func (f *FileInput) parseDataFromBurpFile(httpItemQueue chan *origin.TransferItem) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(f.path); err != nil {
		panic(err)
	}
	root := doc.SelectElement("items")
	for _, item := range root.SelectElements("item") {
		url := item.SelectElement("url")
		rawURL := url.Text()

		request := item.SelectElement("request")
		rawRequestBytes, err := base64.StdEncoding.DecodeString(request.Text())
		if err != nil {
			logger.Error("base64 decode request error\n", err)
			panic(err)
		}

		response := item.SelectElement("response")
		rawResponseBytes, err := base64.StdEncoding.DecodeString(response.Text())
		if err != nil {
			logger.Error("base64 deocde response error\n", err)
			panic(err)
		}

		req, resp := RecoverHttpRequest(string(rawRequestBytes), rawURL, string(rawResponseBytes))

		// transport via channel
		httpItemQueue <- &origin.TransferItem{
			Req:  req,
			Resp: resp,
		}
	}
}
