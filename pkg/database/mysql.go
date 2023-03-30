package database

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/types"
	"APIKiller/pkg/util"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type MysqlConn struct {
	db *gorm.DB
	// the max num of pre-query
	MaxCount     int
	itemAddQueue chan *types.DataItem
}

func NewMysqlConnection(config *config.Config) (*MysqlConn, error) {
	dbCfg := config.Db.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Schema)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorln("Connect database error", err)
		return nil, err
	}

	mysqlConn := &MysqlConn{
		db:           conn,
		itemAddQueue: make(chan *types.DataItem, 100),
	}
	go func() {
		var item *types.DataItem
		for {
			item = <-mysqlConn.ItemAddQueue()
			mysqlConn.AddInfo(item)
		}
	}()

	return mysqlConn, nil
}

func (m *MysqlConn) ItemAddQueue() chan *types.DataItem {
	return m.itemAddQueue
}

func (m *MysqlConn) SetItemAddQueue(itemAddQueue chan *types.DataItem) {
	m.itemAddQueue = itemAddQueue
}

func (m *MysqlConn) ListAllInfo() []types.DataItemStr {
	items := make([]types.DataItemStr, m.MaxCount)
	m.db.Find(&items)

	// recover http item string from id
	for i, item := range items {
		// item.SourceRequest
		items[i].SourceRequest = m.getHttpItemById(item.SourceRequest)

		//item.SourceResponse
		items[i].SourceResponse = m.getHttpItemById(item.SourceResponse)

		//item.VulnRequest
		ids := strings.Split(item.VulnRequest, ",")
		if len(ids) != 0 {
			result := ""
			for _, id := range ids {
				result += "**************************************************\n"
				result += m.getHttpItemById(id)
			}
			items[i].VulnRequest = result
		}

		//item.VulnResponse
		ids2 := strings.Split(item.VulnResponse, ",")
		if len(ids2) != 0 {
			result := ""
			for _, id := range ids2 {
				result += "**************************************************\n"
				result += m.getHttpItemById(id)
			}
			items[i].VulnResponse = result
		}
	}

	return items
}

func (m *MysqlConn) Exist(domain, url, method string) bool {
	var count int64
	v := &types.DataItemStr{}

	m.db.Model(&v).Where("url = ?", url).Where("domain = ?", domain).Where("method = ?", method).Count(&count)

	if count > 0 {
		return true
	}

	return false
}

// AddInfo append new result
func (m *MysqlConn) AddInfo(item *types.DataItem) {
	// transfer DataItem to DataItemStr
	itemStr := types.DataItemStr{
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

// addHttpItem
//
//	@Description: store request or response in form of string and return id
//	@receiver m
//	@param item
//	@return string
func (m *MysqlConn) addHttpItem(itemStr string) string {
	// substr if itemStr is too long
	if len(itemStr) > 10000 {
		itemStr = itemStr[:10000]
	}

	// base64 encode
	b64 := base64.StdEncoding.EncodeToString([]byte(itemStr))

	httpItem := &types.HttpItem{
		Item: b64,
	}

	m.db.Create(&httpItem)

	return fmt.Sprintf("%v", httpItem.Id)
}

func (m *MysqlConn) getHttpItemById(Id string) string {
	// convert string to id
	id, _ := strconv.Atoi(Id)

	item := &types.HttpItem{}

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
func (m *MysqlConn) addHttpItems(items []string) string {
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
