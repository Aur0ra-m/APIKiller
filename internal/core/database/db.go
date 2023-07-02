package database

import (
	"APIKiller/internal/core/data"
)

type Database interface {
	ListAllInfo() []data.DataItemStr
	AddInfo(item *data.DataItem)
	Exist(domain, url, method string) bool
	UpdateVulnType(vulnType string)
}

var (
	saveTaskQueue   chan *data.DataItem
	updateTaskQueue chan string
	db              Database
)

//
// CreateSaveTask
//  @Description: create save result task
//  @param item
//
func CreateSaveTask(item *data.DataItem) {
	saveTaskQueue <- item
}

func CreateUpdateTask(vulnType string) {
	updateTaskQueue <- vulnType
}

//
// BindDatabase
//  @Description: bind global database with provided db object
//  @param database
//
func BindDatabase(database Database) {
	db = database

	// create result save task system
	saveTaskQueue = make(chan *data.DataItem, 1024)
	// result-save queue
	go func() {
		var item *data.DataItem
		for {
			item = <-saveTaskQueue
			db.AddInfo(item)
		}
	}()

	// create update task system
	updateTaskQueue = make(chan string, 1024)
	// update vulnType queue
	go func() {
		var vulnType string
		for {
			vulnType = <-updateTaskQueue
			// update vuln type in db
			db.UpdateVulnType(vulnType)
		}
	}()
}
