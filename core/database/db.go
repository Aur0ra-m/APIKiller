package database

import (
	"APIKiller/core/data"
)

type Database interface {
	ListAllInfo() []data.DataItemStr
	AddInfo(item *data.DataItem)
	Exist(domain, url, method string) bool
	ItemAddQueue() chan *data.DataItem
	SetItemAddQueue(chan *data.DataItem)
}
