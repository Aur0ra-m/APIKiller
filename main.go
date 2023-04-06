package main

import (
	"APIKiller/cmd"
	"APIKiller/core"
	"APIKiller/core/ahttp/hook"
	"APIKiller/core/aio"
	"APIKiller/core/database"
	"APIKiller/core/filter"
	"APIKiller/core/module"
	"APIKiller/core/module/A40xBypasserModule"
	"APIKiller/core/module/CSRFDetector"
	"APIKiller/core/module/DoSDetector"
	"APIKiller/core/module/authorizedDetector"
	"APIKiller/core/module/openRedirectDetector"
	"APIKiller/core/notify"
	"APIKiller/core/origin"
	"APIKiller/core/origin/fileInputOrigin"
	"APIKiller/core/origin/realTimeOrigin"
	logger "APIKiller/logger"
	"APIKiller/web/backend"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"plugin"
	"runtime"
	"strings"
)

const (
	VERSION     = "1.0.0"
	LoggerLevel = logrus.InfoLevel
)

func main() {
	// print Logo
	printLogo()

	// init cmd line
	cmd := cmd.CmdInit()

	// load database\modules\filters\notifier and so on
	loadLogger()
	loadConfig(cmd.ConfigPath)
	loadDatabase()
	loadModules()
	loadFilter()
	loadNotifer()
	loadHooks()

	// start web server and so on
	if cmd.Web {
		logger.Infoln("load web server")
		go backend.NewAPIServer()
	}

	// create a httpItem channel
	httpItemQueue := make(chan *origin.TransferItem)

	// load request from different origins
	go func() {
		if cmd.FileInput != "" {
			//inputOrigin := fileInputOrigin.NewFileInputOrigin("C:\\Users\\Lenovo\\Desktop\\src.txt")
			inputOrigin := fileInputOrigin.NewFileInputOrigin(cmd.FileInput)
			inputOrigin.LoadOriginRequest(httpItemQueue)
		} else {
			inputOrigin := realTimeOrigin.NewRealTimeOrigin()
			inputOrigin.LoadOriginRequest(httpItemQueue)
		}
	}()

	// goroutine control
	limit := make(chan int, cmd.Thread)

	for {
		httpItem := <-httpItemQueue

		// transform io.Reader
		httpItem.Req.Body = aio.TransformReadCloser(httpItem.Req.Body)
		httpItem.Resp.Body = aio.TransformReadCloser(httpItem.Resp.Body)

		// filter requests
		flag := true // true -pass false -block
		for _, f := range filter.Filters {
			if f.Filter(httpItem.Req) == filter.FilterBlocked {
				flag = false

				logger.Infoln(fmt.Sprintf("filter %v, %v", httpItem.Req.Host, httpItem.Req.URL.Path))
				break
			}
		}
		if !flag {
			continue
		}

		go func() {
			limit <- 1

			core.NewHandler(httpItem)

			<-limit
		}()
	}

}

// print logo in form of ascii
func printLogo() {
	fmt.Printf(`
 █████╗ ██████╗ ██╗██╗  ██╗██╗██╗     ██╗     ███████╗██████╗ 
██╔══██╗██╔══██╗██║██║ ██╔╝██║██║     ██║     ██╔════╝██╔══██╗
███████║██████╔╝██║█████╔╝ ██║██║     ██║     █████╗  ██████╔╝
██╔══██║██╔═══╝ ██║██╔═██╗ ██║██║     ██║     ██╔══╝  ██╔══██╗
██║  ██║██║     ██║██║  ██╗██║███████╗███████╗███████╗██║  ██║
╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝
Version: %s`+"\n",
		VERSION)
}

func loadLogger() {
	logger.Initial(LoggerLevel, ".")
}

func loadNotifer() {
	logger.Infoln("loading notifier")

	if viper.GetString("app.notifier.Lark.webhookUrl") != "" {
		notify.BindNotifier(notify.NewLarkNotifier())
	} else if viper.GetString("app.notifier.Dingding.webhookUrl") != "" {
		notify.BindNotifier(notify.NewDingdingNotifer())
	} else {
	}
}

func loadDatabase() {
	logger.Infoln("loading database")

	// bind global database
	database.BindDatabase(database.NewMysqlClient())
}

func loadModules() {
	logger.Infoln("loading modules")

	module.RegisterModule(authorizedDetector.NewAuthorizedDetector())
	module.RegisterModule(A40xBypasserModule.NewA40xBypassModule())
	module.RegisterModule(CSRFDetector.NewCSRFDetector())
	module.RegisterModule(openRedirectDetector.NewOpenRedirectDetector())
	module.RegisterModule(DoSDetector.NewDoSDetector())

}

func loadFilter() {
	logger.Infoln("loading filters")

	filter.RegisterFilter(filter.NewHttpFilter())
	filter.RegisterFilter(filter.NewStaticFileFilter())
	filter.RegisterFilter(filter.NewDuplicateFilter())

}

func loadConfig(configPath string) {
	logger.Infoln("loading config")

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

		var Hookk hook.RequestHook
		Hookk, ok := Hook.(hook.RequestHook)
		if !ok {
			logger.Errorln(fmt.Sprintf("load hook %s error: unexpected type from module symbol", soName))
			panic(err)
		}

		hook.RegisterHooks(Hookk)
	}
}
