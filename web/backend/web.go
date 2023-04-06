package backend

import (
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

type APIServer struct {
	Page int   `form:"page"`
	Size int   `form:"size"`
	Ids  []int `form:"ids"`
	db   *gorm.DB
}

func (s *APIServer) init(ipaddr, port string) {
	// load database
	s.loadDatabase()

	server := gin.Default()

	// append route
	s.route(server)

	// start server
	server.Run(fmt.Sprintf("%s:%s", ipaddr, port))
}

func (s *APIServer) route(server *gin.Engine) {

	group := server.Group("/")
	group.GET("/test", s.test)
	group.GET("/list", s.list)
	group.GET("/check", s.updateCheckState)
}

func (s *APIServer) loadDatabase() {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("root:123456@tcp(192.168.52.153:3306)/apikiller?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// disable logging
	db.Logger.LogMode(1)

	if err != nil {
		log.Errorln("Connect database error", err)
		panic(err)
	}

	s.db = db
}

func (s *APIServer) test(c *gin.Context) {

	c.JSON(http.StatusOK, "Test api")
}

func (s *APIServer) getHttpItembyId(Id string) string {
	// convert string to id
	id, _ := strconv.Atoi(Id)

	item := &data.HttpItem{
		Id: int64(id),
	}

	s.db.Find(item)

	// decode base64
	decodeString, _ := base64.StdEncoding.DecodeString(item.Item)

	return string(decodeString)
}

func (s *APIServer) updateCheckState(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ignore CORS

	logger.Debugln(c.Query("Id"))

	tx := s.db.Model(&data.DataItemStr{}).Where("Id=?", c.Query("Id")).Update("check_state", true)
	if tx.Error != nil {
		logger.Errorln(tx.Error.Error())
	}
	c.JSON(http.StatusOK, "success！")
}

func (s *APIServer) list(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ignore CORS

	_ = c.Bind(&s)

	items := make([]data.DataItemStr, 1024)

	s.db.Where("vuln_type not like ?", "%"+module.AsyncDetectVulnTypeSeperator+"%").Order("domain").Order("url").Find(&items)

	// recover http item string from id
	for i, item := range items {
		// item.SourceRequest
		items[i].SourceRequest = s.getHttpItembyId(item.SourceRequest)

		//item.SourceResponse
		items[i].SourceResponse = s.getHttpItembyId(item.SourceResponse)

		//item.VulnRequest
		items[i].VulnRequest = s.getHttpItembyId(item.VulnRequest)

		//item.VulnResponse
		items[i].VulnResponse = s.getHttpItembyId(item.VulnResponse)
	}

	data := make(map[string]interface{})

	data["list"] = items
	//data["total"] = total
	c.JSON(http.StatusOK, items)
}

func NewAPIServer() {
	server := APIServer{}

	// disable logging
	gin.DefaultWriter = ioutil.Discard

	ipaddr := viper.GetString("app.web.ipaddr")
	port := viper.GetString("app.web.port")

	server.init(ipaddr, port)
}
