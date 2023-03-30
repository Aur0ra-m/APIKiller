package notifier

import (
	"APIKiller/pkg/config"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type Lark struct {
	webhookUrl  string
	secret      string
	signature   string
	timestamp   int64
	notifyQueue chan *types.DataItem
}

func (l *Lark) NotifyQueue() chan *types.DataItem {
	return l.notifyQueue
}

func (l *Lark) SetNotifyQueue(NotifyQueue chan *types.DataItem) {
	l.notifyQueue = NotifyQueue
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
		logger.Error("lark generate signature error\n")
		panic("Lark generate signature error")
	}

	l.signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// NewLarkNotifier
//
//	@Description: create a lark object
//	@param webhook lark webhook url
//	@param signature lark webhook authorize parameter(optional)
//	@return *Lark
func NewLarkNotifier(cfg *config.NotifierConfig) *Lark {
	// get config
	larkCfg := cfg.Lark
	webhookUrl := larkCfg["webhookUrl"]
	secret := larkCfg["secret"]

	// create
	lark := &Lark{
		webhookUrl: webhookUrl,
		signature:  secret,
	}

	//generate signature
	if lark.secret != "" {
		lark.genSign()
	}

	return lark
}

func (l *Lark) GetQueue() chan *types.DataItem {
	return l.notifyQueue
}

func (l *Lark) Notify(item *types.DataItem) {
	logger.Info("notify lark robot\n")

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

	response := gohttp.DoRequest(request)

	defer response.Body.Close()
}
