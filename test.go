package main

import (
	"fmt"
	"regexp"
)

func test() {
	r := regexp.MustCompile(`"` + `key` + `"\s*?:\s*?"?(.*?)?"?,?\s`)
	strings := r.FindStringSubmatch("{\n    \"key\":\"1\"\n}")
	fmt.Println(strings)
	fmt.Println(r.SubexpNames())
}
