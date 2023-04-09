package async

import (
	"APIKiller/core/data"
	"APIKiller/core/database"
	"APIKiller/core/notify"
	"github.com/tidwall/gjson"
	"strings"
	"time"

	"APIKiller/logger"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AsyncCheckEngine struct {
	httpAPI      string
	lastRecordId string
}

func NewAsyncCheckEngine() *AsyncCheckEngine {
	return &AsyncCheckEngine{
		httpAPI:      "http://api.ceye.io/v1/records?token=0920449a5ed8b9db7a287a66a6632498&type=http",
		lastRecordId: "0",
	}
}

func (e *AsyncCheckEngine) Start() {
	// build a request
	request, _ := http.NewRequest("GET", e.httpAPI, nil)
	client := http.Client{}
	// polling API interface
	for {
		// make a http request
		response, _ := client.Do(request)

		// get data from json body
		if response.Body != nil {
			all, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logger.Errorln(err)
			}
			results := gjson.Get(string(all), "data").Array()

			for _, result := range results {
				if result.Get("id").String() <= e.lastRecordId {
					continue
				}

				name := result.Get("name")

				token := strings.Replace(name.String(), "http://zpysri.ceye.io/", "", 1)

				go e.check(token)
			}

			if len(results) > 0 {
				e.lastRecordId = results[0].Get("id").String()
			}

		}

		// sleep
		time.Sleep(5 * 1000 * time.Millisecond)
	}
}

func (e *AsyncCheckEngine) check(token string) {
	logger.Infoln(fmt.Sprintf("[async check] token: %s", token))

	// notify
	notify.CreateNotification(&data.DataItem{
		Id:             "",
		Domain:         "异步检测",
		Url:            "",
		Method:         "",
		Https:          false,
		SourceRequest:  nil,
		SourceResponse: nil,
		VulnType:       token,
		VulnRequest:    nil,
		VulnResponse:   nil,
		ReportTime:     "",
		CheckState:     false,
	})

	// update database
	database.CreateUpdateTask(token)
}
