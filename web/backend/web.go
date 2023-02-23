package backend

import (
	"APIKiller/core/data"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type APIServer struct {
	Page int   `form:"page"`
	Size int   `form:"size"`
	Ids  []int `form:"ids"`
}

func (s *APIServer) init(ipaddr, port string) {
	server := gin.Default()

	group := server.Group("/")
	group.GET("/test", s.test)
	group.GET("/list", s.list)
	group.GET("/check", s.updateCheckState)

	server.Run(fmt.Sprintf("%s:%s", ipaddr, port))
}

func (s *APIServer) db() *gorm.DB {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("root:123456@tcp(192.168.52.153:3306)/apikiller?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Errorln("Connect database error", err)
		panic(err)
	}
	return db
}

func (s *APIServer) test(c *gin.Context) {

	c.JSON(http.StatusOK, "Test api")
}

func (s *APIServer) getHttpItembyId(Id string) string {
	// convert string to id
	id, _ := strconv.Atoi(Id)

	item := &data.HttpItem{}

	s.db().Find(item).Where("id = ?", id)

	// decode base64
	decodeString, _ := base64.StdEncoding.DecodeString(item.Item)

	return string(decodeString)
}

func (s *APIServer) updateCheckState(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ignore CORS

	var v data.DataItemStr
	_ = c.ShouldBindJSON(&v)
	tx := s.db().Model(&v).Where("Id=?", c.PostForm("Id")).Update("CheckState", false)
	if tx.Error != nil {
		logger.Errorln(tx.Error.Error())
	}
	c.JSON(http.StatusOK, "success！")
}

func (s *APIServer) list(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ignore CORS

	_ = c.Bind(&s)

	items := make([]data.DataItemStr, 1024)

	s.db().Find(&items)

	// recover http item string from id
	for i, item := range items {
		// item.SourceRequest
		items[i].SourceRequest = s.getHttpItembyId(item.SourceRequest)

		//item.SourceResponse
		items[i].SourceResponse = s.getHttpItembyId(item.SourceResponse)

		//item.VulnRequest
		ids := strings.Split(item.VulnRequest, ",")
		if len(ids) != 0 {
			result := ""
			for _, id := range ids {
				result += "**************************************************\n"
				result += s.getHttpItembyId(id)
			}
			items[i].VulnRequest = result
		}

		//item.VulnResponse
		ids2 := strings.Split(item.VulnResponse, ",")
		if len(ids2) != 0 {
			result := ""
			for _, id := range ids2 {
				result += "**************************************************\n"
				result += s.getHttpItembyId(id)
			}
			items[i].VulnResponse = result
		}
	}

	data := make(map[string]interface{})

	data["list"] = items
	//data["total"] = total
	c.JSON(http.StatusOK, items)
}

func Server() {
	server := APIServer{}

	server.init("127.0.0.1", "80")
}

func NewAPIServer(ctx context.Context) {
	server := APIServer{}

	ipaddr := util.GetConfig(ctx, "app.web.ipaddr")
	port := util.GetConfig(ctx, "app.web.port")

	server.init(ipaddr, port)
}
