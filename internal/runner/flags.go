package runner

import (
	"flag"
	"os"
)

// NewFlagSet creates a flag structure for the application
func NewFlagSet() *flag.FlagSet {
	flag.CommandLine.ErrorHandling()
	return flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// Parse parses the flags provided to the application
func Parse(flagSet *flag.FlagSet) error {
	flagSet.SetOutput(os.Stdout)
	_ = flagSet.Parse(os.Args[1:])

	return nil
}
