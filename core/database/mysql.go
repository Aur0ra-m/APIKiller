package database

import (
	"APIKiller/core/data"
	log "APIKiller/log"
	"APIKiller/util"
	"context"
	"encoding/base64"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Mysql struct {
	db           *gorm.DB
	MaxCount     int //the max num of per-query
	ItemAddQueue chan *data.DataItem
}

// ListAllInfo fetch all results and return
func (m *Mysql) ListAllInfo() []data.DataItemStr {
	items := make([]data.DataItemStr, m.MaxCount) //需要动态设置，能先查有多少条记录，再创建？

	m.db.Find(&items)

	// recover http item string from id
	for i, item := range items {
		// item.SourceRequest
		items[i].SourceRequest = m.getHttpItembyId(item.SourceRequest)

		//item.SourceResponse
		items[i].SourceResponse = m.getHttpItembyId(item.SourceResponse)

		//item.VulnRequest
		ids := strings.Split(item.VulnRequest, ",")
		if len(ids) != 0 {
			result := ""
			for _, id := range ids {
				result += "**************************************************\n"
				result += m.getHttpItembyId(id)
			}
			items[i].VulnRequest = result
		}

		//item.VulnResponse
		ids2 := strings.Split(item.VulnResponse, ",")
		if len(ids2) != 0 {
			result := ""
			for _, id := range ids2 {
				result += "**************************************************\n"
				result += m.getHttpItembyId(id)
			}
			items[i].VulnResponse = result
		}
	}

	return items
}

func (m *Mysql) GetItemAddQueue() chan *data.DataItem {
	return m.ItemAddQueue
}

func (m *Mysql) Exist(domain, url, method string) bool {
	var count int64
	v := &data.DataItemStr{}

	m.db.Model(&v).Where("url = ?", url).Where("domain = ?", domain).Where("method = ?", method).Count(&count)

	if count > 0 {
		return true
	}

	return false
}

// addInfo append new result
func (m *Mysql) AddInfo(item *data.DataItem) {
	// transfer DataItem to DataItemStr
	itemStr := data.DataItemStr{
		Id:             item.Id,
		Domain:         item.Domain,
		Url:            item.Url,
		Https:          item.Https,
		Method:         item.Method,
		SourceRequest:  m.addHttpItem(util.DumpRequest(item.SourceRequest)),
		SourceResponse: m.addHttpItem(util.DumpResponse(item.SourceResponse)),
		VulnType:       strings.Join(item.VulnType, " "),
		VulnRequest:    m.addHttpItems(util.DumpRequests(item.VulnRequest)),
		VulnResponse:   m.addHttpItems(util.DumpResponses(item.VulnResponse)),
		ReportTime:     item.ReportTime,
		CheckState:     item.CheckState,
	}

	// store DataItemStr
	m.db.Create(&itemStr)
}

//
// addHttpItem
//  @Description: store request or response in form of string and return id
//  @receiver m
//  @param item
//  @return string
//
func (m *Mysql) addHttpItem(itemStr string) string {
	// substr if itemStr is too long
	if len(itemStr) > 10000 {
		itemStr = itemStr[:10000]
	}

	// base64 encode
	b64 := base64.StdEncoding.EncodeToString([]byte(itemStr))

	httpItem := &data.HttpItem{
		Item: b64,
	}

	m.db.Create(&httpItem)

	return fmt.Sprintf("%v", httpItem.Id)
}

func (m *Mysql) getHttpItembyId(Id string) string {
	// convert string to id
	id, _ := strconv.Atoi(Id)

	item := &data.HttpItem{}

	m.db.Find(item).Where("id = ?", id)

	// decode base64
	decodeString, _ := base64.StdEncoding.DecodeString(item.Item)

	return string(decodeString)
}

//
// addHttpItems
//  @Description: store requests or responses in form of string and return ids seperated by comma
//  @receiver m
//  @param item
//  @return string
//
func (m *Mysql) addHttpItems(items []string) string {
	if len(items) == 0 {
		return ""
	}

	var Ids []string

	for _, item := range items {
		Id := m.addHttpItem(item)
		Ids = append(Ids, Id)
	}

	return strings.Join(Ids, ",")
}

//test data: connect("192.168.52.153", "3306","apikiller", "root","123456")
func (m *Mysql) connect(host, port, dbname, username, password string) {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Errorln("Connect database error", err)
		panic(err)
	}

	m.db = db
}

//
// init
//  @Description:
//  @receiver m
//
func (m *Mysql) init() {

}

func NewMysqlClient(ctx context.Context) *Mysql {
	mysqlcli := &Mysql{}

	// init mysql
	mysqlcli.init()

	//parse config
	host := util.GetConfig(ctx, "app.db.mysql.host")
	port := util.GetConfig(ctx, "app.db.mysql.port")
	dbname := util.GetConfig(ctx, "app.db.mysql.dbname")
	username := util.GetConfig(ctx, "app.db.mysql.username")
	password := util.GetConfig(ctx, "app.db.mysql.password")

	//connect db and return DB object
	mysqlcli.connect(host, port, dbname, username, password)

	return mysqlcli
}
