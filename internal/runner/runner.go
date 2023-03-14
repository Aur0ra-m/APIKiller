package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/types"
)

type Runner struct {
	options *types.Options
	config  *config.Config
}

func New(options *types.Options, cfg *config.Config) (*Runner, error) {
	runner := &Runner{
		options: options,
		config:  cfg,
	}

	return runner, nil
}

func (r *Runner) Run() error {
	return nil
}

func (r *Runner) Close() {

}
