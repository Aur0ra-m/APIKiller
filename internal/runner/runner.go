package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/database"
	"APIKiller/pkg/detector"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
)

type Runner struct {
	options   *types.Options
	config    *config.Config
	db        *database.MysqlConn
	detectors []detector.Detector
}

func New(options *types.Options, cfg *config.Config) (*Runner, error) {
	var (
		err       error
		mysqlConn *database.MysqlConn
	)

	if mysqlConn, err = database.NewMysqlConnection(cfg); err != nil {
		logger.Errorf("connect mysql server error: %s", err)
		return nil, err
	}

	detectors := detector.NewDetectors(cfg)

	runner := &Runner{
		options:   options,
		config:    cfg,
		db:        mysqlConn,
		detectors: detectors,
	}

	return runner, nil
}

func (r *Runner) Run() error {
	return nil
}

func (r *Runner) Close() {

}
