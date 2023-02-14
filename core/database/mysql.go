package database

import (
	"APIKiller/core/data"
	log "APIKiller/log"
	"APIKiller/util"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type Mysql struct {
	db       *gorm.DB
	MaxCount int //the max num of per-query
}

// ListAllInfo fetch all results and return
func (m *Mysql) ListAllInfo() []data.DataItemStr {
	items := make([]data.DataItemStr, m.MaxCount) //需要动态设置，能先查有多少条记录，再创建？

	m.db.Find(&items)

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

// AddInfo append new result
func (m *Mysql) AddInfo(item *data.DataItem) {
	// transfer DataItem to DataItemStr
	itemStr := data.DataItemStr{
		Id:             item.Id,
		Domain:         item.Domain,
		Url:            item.Url,
		Https:          item.Https,
		Method:         item.Method,
		SourceRequest:  util.DumpRequest(item.SourceRequest),
		SourceResponse: util.DumpResponse(item.SourceResponse),
		VulnType:       strings.Join(item.VulnType, " "),
		VulnRequest:    util.DumpRequests(item.VulnRequest),
		VulnResponse:   util.DumpResponses(item.VulnResponse),
		ReportTime:     item.ReportTime,
		CheckState:     item.CheckState,
	}

	// store DataItemStr
	m.db.Create(&itemStr)
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

func NewMysqlClient(ctx context.Context) *Mysql {
	mysqlcli := &Mysql{}
	// init mysql

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
