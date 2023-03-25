package main

import (
	"APIKiller/cmd"
	"APIKiller/core"
	"APIKiller/core/ahttp"
	"APIKiller/core/ahttp/hook"
	"APIKiller/core/aio"
	"APIKiller/core/data"
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
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"plugin"
	"runtime"
	"strings"
)

const (
	VERSION = "1.0.0"
)

func main() {
	// print Logo
	printLogo()

	// init cmd line
	cmd := cmd.CmdInit()

	// Context initial
	ctx := context.TODO()

	// load database\modules\filters\notifier and so on
	ctx = loadConfig(ctx, cmd.ConfigPath)
	ctx = loadDatabase(ctx)
	ctx = loadModules(ctx)
	ctx = loadFilter(ctx)
	ctx = loadNotifer(ctx)
	ctx = loadHooks(ctx)

	// start web server and so on
	if cmd.Web {
		logger.Infoln("load web server")
		go backend.NewAPIServer(ctx)
	}

	// create a httpItem channel
	httpItemQueue := make(chan *origin.TransferItem)

	// load request from different origins
	go func() {
		if cmd.FileInput != "" {
			//inputOrigin := fileInputOrigin.NewFileInputOrigin("C:\\Users\\Lenovo\\Desktop\\src.txt")
			inputOrigin := fileInputOrigin.NewFileInputOrigin(cmd.FileInput)
			inputOrigin.LoadOriginRequest(ctx, httpItemQueue)
		} else {
			inputOrigin := realTimeOrigin.NewRealTimeOrigin()
			inputOrigin.LoadOriginRequest(ctx, httpItemQueue)
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
		filters := ctx.Value("filters").([]filter.Filter)

		flag := true // true -pass false -block
		for _, f := range filters {
			if f.Filter(ctx, httpItem.Req) == filter.FilterBlocked {
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

			core.NewHandler(ctx, httpItem)

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

func loadNotifer(ctx context.Context) context.Context {
	logger.Infoln("loading notifier")

	var notifer notify.Notify

	if viper.GetString("app.notifier.Lark.webhookUrl") != "" {
		notifer = notify.NewLarkNotifier(ctx)
	} else if viper.GetString("app.notifier.Dingding.webhookUrl") != "" {
		notifer = notify.NewDingdingNotifer(ctx)
	} else {
		return ctx
	}

	// init notify queue
	notifer.SetNotifyQueue(make(chan *data.DataItem, 30))

	// message queue
	go func() {
		var item *data.DataItem
		for {
			item = <-notifer.NotifyQueue()
			notifer.Notify(item)
		}
	}()

	return context.WithValue(ctx, "notifier", notifer)
}

func loadDatabase(ctx context.Context) context.Context {
	logger.Infoln("loading database")

	db := database.NewMysqlClient(ctx)

	// init queue
	db.SetItemAddQueue(make(chan *data.DataItem, 100))

	// message queue
	go func() {
		var item *data.DataItem
		for {
			item = <-db.ItemAddQueue()
			db.AddInfo(item)
		}
	}()

	return context.WithValue(ctx, "db", db)
}

func loadModules(ctx context.Context) context.Context {
	logger.Infoln("loading modules")

	var modules []module.Detecter

	modules = append(modules, authorizedDetector.NewAuthorizedDetector(ctx))
	modules = append(modules, A40xBypasserModule.NewA40xBypassModule(ctx))
	modules = append(modules, CSRFDetector.NewCSRFDetector(ctx))
	modules = append(modules, openRedirectDetector.NewOpenRedirectDetector(ctx))
	modules = append(modules, DoSDetector.NewDoSDetector(ctx))

	return context.WithValue(ctx, "modules", modules)
}

func loadFilter(ctx context.Context) context.Context {
	logger.Infoln("loading filters")
	var filters []filter.Filter

	filters = append(filters, filter.NewHttpFilter())
	filters = append(filters, filter.NewStaticFileFilter(ctx))
	filters = append(filters, filter.NewDuplicateFilter())

	return context.WithValue(ctx, "filters", filters)
}

func loadConfig(ctx context.Context, configPath string) context.Context {
	logger.Infoln("loading config")

	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return ctx
}

func loadHooks(ctx context.Context) context.Context {
	// except windows os
	if runtime.GOOS == "windows" {
		logger.Infoln("not support windows operation system")
		return ctx
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

		return ctx
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

		var greeter hook.RequestHook
		greeter, ok := Hook.(hook.RequestHook)
		if !ok {
			logger.Errorln(fmt.Sprintf("load hook %s error: unexpected type from module symbol", soName))
			panic(err)
		}

		ahttp.RegisterHooks(greeter)
	}

	return ctx
}
