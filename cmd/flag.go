package cmd

import (
	"flag"
)

type Cmd struct {
	Web       bool
	Thread    int
	FileInput string
}

func CmdInit() *Cmd {
	c := &Cmd{}
	// bind data
	flag.BoolVar(&c.Web, "web", false, "web operations platform option")
	flag.IntVar(&c.Thread, "thread", 100, "go routine concurrency control")
	flag.StringVar(&c.FileInput, "f", "", "load requests from target brup file")

	// parse cmd line
	flag.Parse()
	//fmt.Println(c)
	return c
}
