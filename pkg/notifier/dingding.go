package notifier

import (
	"APIKiller/pkg/config"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
	"bytes"
	"fmt"
	"net/http"
)

type Dingding struct {
	webhookUrl  string
	notifyQueue chan *types.DataItem
}

func (d *Dingding) NotifyQueue() chan *types.DataItem {
	return d.notifyQueue
}

func (d *Dingding) SetNotifyQueue(NotifyQueue chan *types.DataItem) {
	d.notifyQueue = NotifyQueue
}

func (d *Dingding) Notify(item *types.DataItem) {
	logger.Info("notify dingding robot\n")

	var jsonData []byte

	// Message format setting
	MessageFormat := fmt.Sprintf("%s-%s exists %s", item.Domain, item.Url, item.VulnType)

	jsonData = []byte(fmt.Sprintf(`{
    "at": {
        "isAtAll": true
    },
    "text": {
        "content":"%s"
    },
    "msgtype":"text"
}`, MessageFormat))

	request, _ := http.NewRequest("POST", d.webhookUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response := gohttp.DoRequest(request)

	defer response.Body.Close()
}

func NewDingdingNotifer(cfg *config.NotifierConfig) *Dingding {
	// get config
	dingCfg := cfg.Dingding
	webhookUrl := dingCfg["webhookUrl"]
	// create
	notifer := &Dingding{webhookUrl: webhookUrl}

	return notifer
}
