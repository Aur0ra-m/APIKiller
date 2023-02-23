package notify

import (
	http2 "APIKiller/core/ahttp"
	"APIKiller/core/data"
	logger "APIKiller/log"
	"APIKiller/util"
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Dingding struct {
	webhookUrl  string
	NotifyQueue chan *data.DataItem
}

func (d *Dingding) Notify(item *data.DataItem) {
	logger.Infoln("notify dingding robot")

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

	response := http2.DoRequest(request, false)

	defer response.Body.Close()
}

func (d *Dingding) GetQueue() chan *data.DataItem {
	return d.NotifyQueue
}

func NewDingdingNotifier(ctx context.Context) *Dingding {
	// get config
	webhookUrl := util.GetConfig(ctx, "app.notifier.Dingding.webhookUrl")
	// create
	return &Dingding{webhookUrl: webhookUrl}
}
