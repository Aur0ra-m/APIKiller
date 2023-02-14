package database

import (
	"fmt"
	"testing"
)

func TestMysql_Init(t *testing.T) {

}

func TestMysql_ListAllInfo(t *testing.T) {
	fmt.Println("Test")
	m := new(Mysql)
	m.connect("192.168.52.153", "3306", "apikiller", "root", "123456")
	m.ListAllInfo()
}
