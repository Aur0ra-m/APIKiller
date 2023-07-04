package runner

import (
	"APIKiller/internal/core"
	hook2 "APIKiller/internal/core/ahttp/hook"
	"APIKiller/internal/core/aio"
	"APIKiller/internal/core/async"
	database2 "APIKiller/internal/core/database"
	filter2 "APIKiller/internal/core/filter"
	"APIKiller/internal/core/module"
	"APIKiller/internal/core/module/CSRF"
	"APIKiller/internal/core/module/DoS"
	"APIKiller/internal/core/module/OpenRedirect"
	"APIKiller/internal/core/module/SSRF"
	"APIKiller/internal/core/module/authorize"
	notify2 "APIKiller/internal/core/notify"
	"APIKiller/internal/core/origin"
	"APIKiller/internal/core/origin/fileInputOrigin"
	"APIKiller/internal/core/origin/realTimeOrigin"
	"APIKiller/internal/web/backend"
	"APIKiller/pkg/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"plugin"
	"runtime"
	"strings"
)

const (
	VERSION     = "1.2.0"
	LoggerLevel = logrus.InfoLevel
)

type Runner struct {
}

//
// Start
//  @Description: start the core application
//  @receiver r
//
func (r *Runner) Start(cmdOptions *CommandOptions) {

	// load database\modules\filters\notifier and so on
	loadLogger()
	loadConfig(cmdOptions.ConfigPath)
	loadDatabase()
	loadModules()
	loadAsyncCheckEngine()
	loadFilter()
	loadNotifer()
	loadHooks()

	// load ops platform
	if cmdOptions.Web {
		loadOPSPlatform()
	}

	// load request data from different origins
	loadRequestDatafromOrigin(cmdOptions.FileInput)

	// load data from channel and start to handle
	startHandle(cmdOptions.Thread)

	//
}

//
// Stop
//  @Description: stop the core application
//  @receiver r
//
func (r *Runner) Stop() {

}

//
// NewRunner
//  @Description:
//
func NewRunner() {
	R := &Runner{}

	// show banner
	showBanner()

	// parse command option
	cmdOptions := ParseCommandOptions()

	// start runner
	R.Start(cmdOptions)

	// and so on
}

func startHandle(MaxThreadNum int) {
	logger.Infoln("start handle")

	// goroutine control
	limit := make(chan int, MaxThreadNum)

	for {
		transferItem := origin.GetOriginRequest()

		// filter requests
		flag := true // true -pass false -block
		for _, f := range filter2.GetFilters() {
			if f.Filter(transferItem.Req) == filter2.FilterBlocked {
				flag = false

				logger.Infoln(fmt.Sprintf("filter %v, %v", transferItem.Req.Host, transferItem.Req.URL.Path))
				break
			}
		}
		if !flag {
			continue
		}

		// transform io.Reader
		transferItem.Req.Body = aio.TransformReadCloser(transferItem.Req.Body)
		transferItem.Resp.Body = aio.TransformReadCloser(transferItem.Resp.Body)

		go func() {
			limit <- 1

			core.NewHandler(transferItem)

			<-limit
		}()
	}
}

func loadRequestDatafromOrigin(filePath string) {
	// load request from different origins
	go func() {
		if filePath != "" {
			inputOrigin := fileInputOrigin.NewFileInputOrigin(filePath)
			inputOrigin.LoadOriginRequest()
		} else {
			inputOrigin := realTimeOrigin.NewRealTimeOrigin()
			inputOrigin.LoadOriginRequest()
		}
	}()
}

func loadOPSPlatform() {
	logger.Infoln("loading OPS platform")

	go backend.NewAPIServer()
}

func loadLogger() {
	logger.Initial(LoggerLevel, ".")
}

func loadNotifer() {
	logger.Infoln("loading notifier")

	if viper.GetString("app.notifier.Lark.webhookUrl") != "" {
		notify2.BindNotifier(notify2.NewLarkNotifier())
	} else if viper.GetString("app.notifier.Dingding.webhookUrl") != "" {
		notify2.BindNotifier(notify2.NewDingdingNotifer())
	} else {
	}
}

func loadDatabase() {
	logger.Infoln("loading database")

	// bind global database
	database2.BindDatabase(database2.NewMysqlClient())
}

func loadModules() {
	logger.Infoln("loading modules")

	module.RegisterModule(authorize.NewAuthorizedDetector())
	module.RegisterModule(CSRF.NewCSRFDetector())
	module.RegisterModule(OpenRedirect.NewOpenRedirectDetector())
	module.RegisterModule(DoS.NewDoSDetector())
	module.RegisterModule(SSRF.NewSSRFDetector())

}

func loadFilter() {
	logger.Infoln("loading filters")

	filter2.RegisterFilter(filter2.NewHttpFilter())
	filter2.RegisterFilter(filter2.NewStaticFileFilter())
	filter2.RegisterFilter(filter2.NewDuplicateFilter())

}

func loadConfig(configPath string) {
	logger.Infoln("loading config")

	// use the specified configuration file when configPath option is not blank
	if configPath == "" {
		// using environment variable
		env := os.Getenv("APIKiller_env")
		if env == "dev" || true {
			configPath = "./config/config.dev.yaml"
		} else {
			configPath = "./config/config.release.yaml"
		}
	}

	logger.Debugln(fmt.Sprintf("current config: %s", configPath))

	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func loadHooks() {
	// except windows os
	if runtime.GOOS == "windows" {
		logger.Infoln("not support windows operation system")
	}

	logger.Infoln("loading hooks")

	// ./hooks directory does not exist
	_, err2 := os.Stat("./hooks")
	if os.IsNotExist(err2) {
		logger.Errorln("target directory does not exist")

		// make directory
		err := os.Mkdir("./hooks", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// list directory
	entries, err := os.ReadDir("./hooks")
	if err != nil {
		logger.Errorln(fmt.Sprintf("loading hooks error: %v", err))
		panic(entries)
	}

	for _, entry := range entries {
		soName := entry.Name()

		// filter directory and none so file
		if entry.IsDir() == true || strings.Index(soName, ".so") == -1 {
			continue
		}

		// load plugins and register them via RegisterHooks
		logger.Infoln(fmt.Sprintf("[Load Hook] load hook %s", strings.Replace(soName, ".so", "", 1)))
		open, err := plugin.Open("./hooks/" + soName)
		if err != nil {
			logger.Errorln(fmt.Sprintf("load hook %s error: %v", soName, err))
			panic(err)
		}

		Hook, err := open.Lookup("Hook")
		if err != nil {
			logger.Errorln(fmt.Sprintf("load hook %s error: %v", soName, err))
			panic(err)
		}

		var Hookk hook2.RequestHook
		Hookk, ok := Hook.(hook2.RequestHook)
		if !ok {
			logger.Errorln(fmt.Sprintf("load hook %s error: unexpected type from module symbol", soName))
			panic(err)
		}

		hook2.RegisterHooks(Hookk)
	}
}

func loadAsyncCheckEngine() {
	logger.Infoln("loading asynchronous check engine")

	// start asynchronous check engine
	go async.NewAsyncCheckEngine().Start()
}
