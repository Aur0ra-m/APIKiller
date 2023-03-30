package main

import (
	"APIKiller/internal/runner"
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
	"flag"
	"os"
	"os/signal"
)

var (
	options = &types.Options{}
)

func main() {
	_ = readConfig()

	runner.ParseOptions(options)
	cfg := config.GetConf()
	apiRunner, err := runner.New(options, cfg)
	if err != nil {
		logger.Fatalf("Could not create runner: %s\n", err)
	}
	if apiRunner == nil {
		return
	}

	// Setup graceful exits
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Info("CTRL+C pressed: Exiting\n")
			apiRunner.Close()
			os.Exit(1)
		}
	}()

	if err := apiRunner.Run(); err != nil {
		logger.Fatalf("Could not run ApiKiller: %s\n", err)
	}
	apiRunner.Close()
}

func readConfig() *flag.FlagSet {
	flagSet := runner.NewFlagSet()

	flagSet.BoolVar(&options.EnableWeb, "web", false, "enable web platform")
	flagSet.IntVar(&options.Thread, "thread", 100, "maximum number of threads to be executed in parallel")
	flagSet.StringVar(&options.BurpFile, "burp-file", "", "load http requests from burp file")
	flagSet.StringVar(&options.ConfigFile, "config", "", "path to the configuration file")

	// parse go flags
	_ = runner.Parse(flagSet)
	// load config file
	config.Setup(options)

	return flagSet
}
