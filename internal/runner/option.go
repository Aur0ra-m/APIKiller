package runner

import "flag"

type CommandOptions struct {
	ConfigPath string
	Web        bool
	Thread     int
	FileInput  string
}

//
// ParseCommandOptions
//  @Description: parse command options through flag package
//  @return *CommandOptions
//
func ParseCommandOptions() *CommandOptions {
	c := &CommandOptions{}
	// bind data
	flag.StringVar(&c.ConfigPath, "conf", "", "project config path")
	flag.BoolVar(&c.Web, "web", false, "web operations platform option")
	flag.IntVar(&c.Thread, "thread", 100, "go routine concurrency control")
	flag.StringVar(&c.FileInput, "f", "", "load requests from target brup file")

	// parse cmd line
	flag.Parse()
	//fmt.Println(c)
	return c
}
