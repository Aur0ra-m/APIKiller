package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/types"
)

type Runner struct {
	options *types.Options
	config  *config.Config
}

func New(options *types.Options, config *config.Config) (*Runner, error) {

	return &Runner{
		options: options,
		config:  config,
	}, nil
}
