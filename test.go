package main

import (
	"fmt"
	"regexp"
)

func test() {
	r := regexp.MustCompile("name=\"" + "description" + "\".*?" + "\n\n" + `(.*)`)

	s := "---------------------------974767299852498929531610575\nContent-Disposition: form-data; name=\"description\"\n\nsome text"

	strings := r.FindStringSubmatch(s)
	fmt.Println(strings)
	fmt.Println(r.SubexpNames())
}
