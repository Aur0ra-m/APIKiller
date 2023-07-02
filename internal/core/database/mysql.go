package database

import (
	"APIKiller/internal/core/ahttp"
	"APIKiller/internal/core/data"
	"APIKiller/internal/core/module"
	log "APIKiller/pkg/logger"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Mysql struct {
	db       *gorm.DB
	MaxCount int //the max num of per-query
}

func (m *Mysql) UpdateVulnType(token string) {
	m.db.Model(&data.DataItemStr{}).Where("vuln_type = ?", token).Update("vuln_type", strings.Split(token, module.AsyncDetectVulnTypeSeperator)[0])
}

// ListAllInfo fetch all results and return
func (m *Mysql) ListAllInfo() []data.DataItemStr {
	items := make([]data.DataItemStr, m.MaxCount) //需要动态设置，能先查有多少条记录，再创建？

	m.db.Where("vuln_type not like ?", "%"+module.AsyncDetectVulnTypeSeperator+"%").Order("domain").Order("url").Find(&items)

	// recover http item string from id
	for i, item := range items {
		// item.SourceRequest
		items[i].SourceRequest = m.getHttpItembyId(item.SourceRequest)

		//item.SourceResponse
		items[i].SourceResponse = m.getHttpItembyId(item.SourceResponse)

		//item.VulnRequest
		items[i].VulnRequest = m.getHttpItembyId(item.VulnRequest)

		//item.VulnResponse
		items[i].VulnResponse = m.getHttpItembyId(item.VulnResponse)
	}

	return items
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
		SourceRequest:  m.addHttpItem(ahttp.DumpRequest(item.SourceRequest)),
		SourceResponse: m.addHttpItem(ahttp.DumpResponse(item.SourceResponse)),
		VulnType:       item.VulnType,
		VulnRequest:    m.addHttpItem(ahttp.DumpRequest(item.VulnRequest)),
		VulnResponse:   m.addHttpItem(ahttp.DumpResponse(item.VulnResponse)),
		ReportTime:     item.ReportTime,
		CheckState:     item.CheckState,
	}

	// store DataItemStr
	m.db.Create(&itemStr)
}

// addHttpItem
//
//	@Description: store request or response in form of string and return id
//	@receiver m
//	@param item
//	@return string
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

// addHttpItems
//
//	@Description: store requests or responses in form of string and return ids seperated by comma
//	@receiver m
//	@param item
//	@return string
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

// test data: connect("192.168.52.153", "3306","apikiller", "root","123456")
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

// init
//
//	@Description:
//	@receiver m
func (m *Mysql) init() {
	// disable logging
	m.db.Logger.LogMode(1)
}

func NewMysqlClient() *Mysql {
	mysqlcli := &Mysql{}

	//parse config
	host := viper.GetString("app.db.mysql.host")
	port := viper.GetString("app.db.mysql.port")
	dbname := viper.GetString("app.db.mysql.dbname")
	username := viper.GetString("app.db.mysql.username")
	password := viper.GetString("app.db.mysql.password")

	//connect db and return DB object
	mysqlcli.connect(host, port, dbname, username, password)

	// init mysql
	mysqlcli.init()

	return mysqlcli
}
