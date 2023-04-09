package notify

import (
	"APIKiller/core/data"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type Dingding struct {
	webhookUrl string
}

func (d *Dingding) Notify(item *data.DataItem) {
	//logger.Infoln("notify dingding robot")

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

	client := http.Client{}
	response, _ := client.Do(request)

	defer response.Body.Close()
}

func NewDingdingNotifer() *Dingding {
	// get config
	webhookUrl := viper.GetString("app.notifier.Dingding.webhookUrl")
	// create
	notifer := &Dingding{webhookUrl: webhookUrl}

	return notifer
}
