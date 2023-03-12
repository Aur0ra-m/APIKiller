package config

import (
	"APIKiller/v1/pkg/types"
	"fmt"
	"os"
	"strings"
)

const VERSION = "0.0.4"

type Config struct {
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
	return strings.Join([]string{pwd, "config", "config.yaml"}, "/")
}

func main() {
	fmt.Println(getDefaultConfigFile())
}
