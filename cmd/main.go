package main

import (
	runner2 "APIKiller/internal/runner"
	"APIKiller/pkg/config"
	"APIKiller/pkg/types"
	"flag"
)

var (
	options = &types.Options{}
)

func main() {
	_ = readConfig()

	runner2.ParseOptions(options)
}

func readConfig() *flag.FlagSet {
	flagSet := runner2.NewFlagSet()

	flagSet.BoolVar(&options.EnableWeb, "web", false, "enable web platform")
	flagSet.IntVar(&options.Thread, "thread", 100, "maximum number of threads to be executed in parallel")
	flagSet.StringVar(&options.BurpFile, "burp-file", "", "load http requests from burp file")
	flagSet.StringVar(&options.ConfigFile, "config", "", "path to the configuration file")

	// parse go flags
	_ = runner2.Parse(flagSet)
	// load config file
	config.Setup(options)

	return flagSet
}
