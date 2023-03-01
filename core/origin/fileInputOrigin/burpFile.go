package fileInputOrigin

import (
	"APIKiller/core/origin"
	logger "APIKiller/log"
	"encoding/base64"
	"github.com/beevik/etree"
)

// parseData
//
//	@Description: parse data from burpsuite file
//	@receiver o
func (o *FileInputOrigin) parseDataFromBurpFile(httpItemQueue chan *origin.TransferItem) {

	doc := etree.NewDocument()

	if err := doc.ReadFromFile(o.path); err != nil {
		panic(err)
	}

	root := doc.SelectElement("items")
	for _, item := range root.SelectElements("item") {
		url := item.SelectElement("url")
		//fmt.Println(url.Text())
		rawUrl := url.Text()

		request := item.SelectElement("request")
		rawRequestBytes, err2 := base64.StdEncoding.DecodeString(request.Text())
		if err2 != nil {
			logger.Errorln("base64 decode error", err2)
			panic(err2)
		}

		response := item.SelectElement("response")
		rawResponseBytes, err3 := base64.StdEncoding.DecodeString(response.Text())
		if err2 != nil {
			logger.Errorln("base64 decode error", err3)
			panic(err3)
		}

		req, resp := RecoverHttpRequest(string(rawRequestBytes), rawUrl, string(rawResponseBytes))

		//transport via channel
		httpItemQueue <- &origin.TransferItem{
			Req:  req,
			Resp: resp,
		}
	}

	return
}
