package runner

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/database"
	"APIKiller/pkg/detector"
	"APIKiller/pkg/filter"
	"APIKiller/pkg/hooks"
	"APIKiller/pkg/http"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/notifier"
	"APIKiller/pkg/origin"
	"APIKiller/pkg/origin/fileinput"
	"APIKiller/pkg/origin/realtime"
	"APIKiller/pkg/types"
	"APIKiller/web/backend"
	"fmt"
	"math/rand"
	"time"
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

	return &Runner{
		options:   options,
		config:    cfg,
		db:        mysqlConn,
		detectors: detectors,
		filters:   filters,
		notify:    notify,
	}, nil
}

func (r *Runner) Run() error {
	if r.config.Web.Enable {
		logger.Info("load web server\n")
		go backend.NewAPIServer(&r.config.Web)
	}

	// create a httpItem channel
	httpItemQueue := make(chan *origin.TransferItem)
	// load request from different origins
	go func() {
		if r.options.BurpFile != "" {
			inputOrigin := fileinput.NewFileInputOrigin(r.options.BurpFile)
			inputOrigin.LoadOriginRequest(&r.config.Origin, httpItemQueue)
		} else {
			inputOrigin := realtime.NewRealTimeOrigin()
			inputOrigin.LoadOriginRequest(&r.config.Origin, httpItemQueue)
		}
	}()

	// goroutine control
	limit := make(chan int, r.options.Thread)
	filtered := make(chan bool, 1)
	passed := make(chan bool, 1)
	for {
		httpItem := <-httpItemQueue
		// transform io.Reader
		httpItem.Req.Body = http.TransformReadCloser(httpItem.Req.Body)
		httpItem.Resp.Body = http.TransformReadCloser(httpItem.Resp.Body)

		go func() {
			for _, f := range r.filters {
				if f.Filter(httpItem.Req) == filter.FilterBlocked {
					filtered <- true
					return
				}
			}
			passed <- true
		}()

		select {
		case <-filtered:
			continue
		case <-passed:
			limit <- 1
			r.handle(httpItem)
			<-limit
		default:
			logger.Error("something wrong!\n")
		}
	}
}

func (r *Runner) handle(httpItem *origin.TransferItem) {
	req := httpItem.Req

	// assembly DataItem
	item := &types.DataItem{
		Id:             fmt.Sprintf("%v%v", time.Now().Unix(), rand.Int()),
		Domain:         req.Host,
		Url:            req.URL.Path,
		Https:          req.URL.Scheme == "https",
		Method:         req.Method,
		SourceRequest:  req,
		SourceResponse: httpItem.Resp,
		VulnType:       []string{},
		VulnRequest:    nil,
		VulnResponse:   nil,
		ReportTime:     fmt.Sprintf("%v", time.Now().Unix()),
		CheckState:     false,
	}

	for _, d := range r.detectors {
		if d != nil {
			d.Detect(item)
		}
	}

	if r.notify != nil {
		r.notify.NotifyQueue() <- item
	}
	logger.Infof("%v %v checkout: %v\n", item.Domain, item.Url, item.VulnType)
	r.db.ItemAddQueue() <- item
}

func (r *Runner) Close() {

}
