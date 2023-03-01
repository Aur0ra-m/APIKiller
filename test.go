package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func test() {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("app"))
}
