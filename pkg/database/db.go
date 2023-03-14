package database

import (
	"APIKiller/pkg/types"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	ProtocolMysql     = "mysql"
	ProtocolRedis     = "redis"
	ProtocolSQLServer = "sqlserver"
	ProtocolMariadb   = "mariadb"
)

type ErrNoClient struct {
	Name string
}

func (e ErrNoClient) Error() string {
	return fmt.Sprintf("not found %s client", e.Name)
}

var (
	ErrUnSupportedProtocol = errors.New("unsupported protocol")

	ErrRedisClient     = ErrNoClient{"Redis"}
	ErrMySQLClient     = ErrNoClient{"MySQL"}
	ErrSQLServerClient = ErrNoClient{"SQLServer"}
)

type supportedChecker func() error

var supportedMap = map[string]supportedChecker{
	ProtocolMysql:     mysqlSupported,
	ProtocolRedis:     redisSupported,
	ProtocolSQLServer: sqlServerSupported,
	ProtocolMariadb:   mysqlSupported,
}

type Database interface {
	ListAllInfo() []types.DataItemStr
	AddInfo(item *types.DataItem)
	Exist(domain, url, method string) bool
	ItemAddQueue() chan *types.DataItem
	SetItemAddQueue(chan *types.DataItem)
}

func IsSupportedProtocol(p string) error {
	if checker, ok := supportedMap[p]; ok {
		return checker()
	}

	return ErrUnSupportedProtocol
}

func mysqlSupported() error {
	checkLine := "mysql -V"
	cmd := exec.Command("bash", "-c", checkLine)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return fmt.Errorf("%w: %s", ErrMySQLClient, err)
	}
	if bytes.HasPrefix(out, []byte("mysql")) {
		return nil
	}
	return ErrMySQLClient
}

func redisSupported() error {
	checkLine := "redis-cli -v"
	cmd := exec.Command("bash", "-c", checkLine)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return fmt.Errorf("%w: %s", ErrRedisClient, err)
	}
	if bytes.HasPrefix(out, []byte("redis-cli")) {
		return nil
	}
	return ErrRedisClient
}

func sqlServerSupported() error {
	checkLine := "tsql -C"
	cmd := exec.Command("bash", "-c", checkLine)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return fmt.Errorf("%w: %s", ErrSQLServerClient, err)
	}
	if strings.Contains(string(out), "freetds") {
		return nil
	}
	return ErrSQLServerClient
}
