package main

import (
	"APIKiller/v1/pkg/types"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const VERSION = "0.0.4"

type Config struct {
	Db       DbConfig       `mapstructure:"db"`
	Origin   OriginConfig   `mapstructure:"origin"`
	Detector DetectorConfig `mapstructure:"detector"`
	Filter   FilterConfig   `mapstructure:"filter"`
	Notifier NotifierConfig `mapstructure:"notifier"`
	Web      WebConfig      `mapstructure:"web"`
}

type DbConfig struct {
	Mysql MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Schema   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type DetectorConfig struct {
	Authorize    AuthorizedConfig   `mapstructure:"authorizedDetector"`
	A40xBypass   A40xBypassConfig   `mapstructure:"40xBypassDetector"`
	Csrf         CsrfConfig         `mapstructure:"csrfDetector"`
	OpenRedirect OpenRedirectConfig `mapstructure:"openRedirectDetector"`
	Dos          DosDetector        `mapstructure:"dosDetector"`
}

type AuthorizedConfig struct {
	Enable     bool                   `mapstructure:"enable"`
	AuthHeader string                 `mapstructure:"authHeader"`
	Roles      []string               `mapstructure:"roles"`
	Judgement  map[string]interface{} `mapstructure:"judgement"`
}

type A40xBypassConfig struct {
	Enable       bool                   `mapstructure:"enable"`
	AuthFailFlag map[string]interface{} `mapstructure:"authFailedFlag"`
	IpHeader     []string               `mapstructure:"ipHeader"`
	Ip           string                 `mapstructure:"ip"`
	ApiVersion   map[string]string      `mapstructure:"apiVersion"`
	PathFuzz     map[string]interface{} `mapstructure:"pathFuzz"`
}

type CsrfConfig struct {
	Enable             bool   `mapstructure:"enable"`
	CsrfToken          string `mapstructure:"csrfToken"`
	CsrfInvalidPattern string `mapstructure:"csrfInvalidPattern"`
}

type OpenRedirectConfig struct {
	Enable         bool     `mapstructure:"enable"`
	RawQueryParams []string `mapstructure:"rawQueryParams"`
	FailFlag       []string `mapstructure:"failFlag"`
}

type DosDetector struct {
	Enable    bool                   `mapstructure:"enable"`
	SizeParam []string               `mapstructure:"sizeParam"`
	RateLimit map[string]interface{} `mapstructure:"rateLimit"`
}

type OriginConfig struct {
	RealTime map[string]string `mapstructure:"realTime"`
}

type FilterConfig struct {
	Http       HttpFilterConfig       `mapstructure:"httpFilter"`
	StaticFile StaticFileFilterConfig `mapstructure:"staticFileFilter"`
}

type HttpFilterConfig struct {
	Host []string `mapstructure:"host"`
}

type StaticFileFilterConfig struct {
	Ext []string `mapstructure:"ext"`
}

type NotifierConfig struct {
	Lark     map[string]string `mapstructure:"lark"`
	Dingding map[string]string `mapstructure:"dingding"`
}

type WebConfig struct {
	Host string `mapstructure:"host"`
	port string `mapstructure:"port"`
}

func have(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetConfigFile(options *types.Options) string {
	if options.ConfigFile != "" {
		return options.ConfigFile
	}
	return getDefaultConfigFile()
}

func getDefaultConfigFile() string {
	pwd, err := os.Getwd()
	if err != nil {
		// TODO log error info
	}
	return strings.Join([]string{pwd, "v1/pkg/config", "config.yaml"}, "/")
}

func loadConfigFromFile(path string, conf *Config) {
	var err error
	if have(path) {
		fileViper := viper.New()
		fileViper.SetConfigFile(path)
		if err = fileViper.ReadInConfig(); err == nil {
			if err = fileViper.Unmarshal(conf); err == nil {
				log.Printf("Load config from %s success\n", path)
				return
			}
		}
	}

	if err != nil {
		log.Fatalf("Load config from %s failed: %s\n", path, err)
	}
}

func main() {
	config := Config{}
	cfgPath := getDefaultConfigFile()
	loadConfigFromFile(cfgPath, &config)

}
