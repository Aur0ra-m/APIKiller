package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/database"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
)

type Runner struct {
	options *types.Options
	config  *config.Config
	db      *database.MysqlConn
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

	runner := &Runner{
		options: options,
		config:  cfg,
		db:      mysqlConn,
	}

	return runner, nil
}

func (r *Runner) Run() error {
	return nil
}

func (r *Runner) Close() {

}
