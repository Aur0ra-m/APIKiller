package notify

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	logger "APIKiller/log"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type Lark struct {
	webhookUrl string
	secret     string
	signature  string
	timestamp  int64
}

func (l *Lark) genSign() {
	//get timestamp
	l.timestamp = time.Now().Unix()

	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", l.timestamp) + "\n" + l.secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		logger.Errorln("lark generate signature error")
		panic("Lark generate signature error")
	}

	l.signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (l *Lark) init() {
	//generate signature
	if l.secret != "" {
		l.genSign()
	}
}

//
// NewLarkNotifier
//  @Description: create a lark object
//  @param webhook lark webhook url
//  @param signature lark webhook authorize parameter(optional)
//  @return *Lark
//
func NewLarkNotifier(webhookUrl, secret string) *Lark {
	lark := &Lark{
		webhookUrl: webhookUrl,
		signature:  secret,
	}

	// init object
	lark.init()

	return lark
}

func (l *Lark) Notify(ctx context.Context, item *data.DataItem) {
	logger.Infoln("notify lark robot")

	var jsonData []byte

	// Message format setting
	MessageFormat := fmt.Sprintf("%s-%s exists %s", item.Domain, item.Url, item.VulnType)

	if l.secret != "" {
		jsonData = []byte(fmt.Sprintf(`
		{
				"timestamp": "%v",
				"sign": "%v",
				"msg_type": "text",
				"content": {
						"text": "%v"
				}
		}`, l.timestamp, l.signature, MessageFormat))
	} else {
		jsonData = []byte(fmt.Sprintf(`{"msg_type":"text","content":{"text":"%v"}}`, MessageFormat))
	}

	request, _ := http.NewRequest("POST", l.webhookUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response := http2.DoRequest(request)

	defer response.Body.Close()
}
