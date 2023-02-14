package util

import (
	"context"
	"github.com/tidwall/gjson"
	"strings"
)

var (
	ConfigSeperator = ","
)

//
// GetConfig
//  @Description: Get config in format of string. Every config endpoint must use string, cannot use other
//  @param ctx
//  @param key
//  @return string
//
func GetConfig(ctx context.Context, key string) string {
	conf := ctx.Value("configJsonString").(string)

	result := gjson.Get(conf, key).String()

	// trim the last comma at the end of string
	for result != "" && result[len(result)-1] == ',' {
		result = result[:len(result)-1]
	}

	return result
}

func SplitConfigString(configString string) []string {
	if configString == "" {
		return []string{}
	}

	return strings.Split(configString, ConfigSeperator)
}
