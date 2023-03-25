package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/database"
	"APIKiller/pkg/detector"
	"APIKiller/pkg/filter"
	"APIKiller/pkg/hooks"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/notifier"
	"APIKiller/pkg/types"
	"APIKiller/web/backend"
)

type Runner struct {
	options   *types.Options
	config    *config.Config
	db        *database.MysqlConn
	detectors []detector.Detector
	filters   []filter.Filter
	notify    notifier.Notify
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

	detectors := detector.NewDetectors(&cfg.Detector)
	filters := filter.NewFilter(&cfg.Filter)
	notify := notifier.NewNotify(&cfg.Notifier)
	hooks.NewHook(cfg)

	if cfg.Web.Enable {
		logger.Info("load web server\n")
		go backend.NewAPIServer(&cfg.Web)
	}

	// create a httpItem channel
	go func() {
		if options.ConfigFile != "" {

		}
	}()

	runner := &Runner{
		options:   options,
		config:    cfg,
		db:        mysqlConn,
		detectors: detectors,
		filters:   filters,
		notify:    notify,
	}

	return runner, nil
}

func (r *Runner) Run() error {
	return nil
}

func (r *Runner) Close() {

}
