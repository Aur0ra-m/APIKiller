package main

import (
	"APIKiller/cmd"
	"APIKiller/core"
	"APIKiller/core/aio"
	"APIKiller/core/database"
	"APIKiller/core/filter"
	"APIKiller/core/module"
	"APIKiller/core/module/authorizedDetector"
	"APIKiller/core/module/csrfDetector"
	"APIKiller/core/notify"
	"APIKiller/core/origin"
	"APIKiller/core/origin/fileInputOrigin"
	"APIKiller/core/origin/realTimeOrigin"
	logger "APIKiller/log"
	"APIKiller/util"
	"APIKiller/web/backend"
	"context"
	"fmt"
	"os"
)

const (
	VERSION = "0.0.1"
)

func main() {
	// print Logo
	printLogo()

	// init cmd line
	cmd := cmd.CmdInit()

	// Context initial
	ctx := context.TODO()

	// load database\modules\filters\notifier and so on
	ctx = loadConfigJsonStr(ctx)
	ctx = loadDatabase(ctx)
	ctx = loadModules(ctx)
	ctx = loadFilter(ctx)
	ctx = loadNotifer(ctx)

	// start web server and so on
	cmd.Web = false
	if cmd.Web {
		logger.Infoln("load web server")
		go backend.NewAPIServer(ctx)
	}

	// create a httpItem channel
	httpItemQueue := make(chan *origin.HttpItem)

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

				logger.Debugln("filter %v, %v", httpItem.Req.Host, httpItem.Req.URL.RawPath)
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

	webhookUrl := util.GetConfig(ctx, "app.notifier.Lark.webhookUrl")
	secret := util.GetConfig(ctx, "app.notifier.Lark.secret")

	notifer := notify.NewLarkNotifier(webhookUrl, secret)

	return context.WithValue(ctx, "notifier", notifer)
}

func loadDatabase(ctx context.Context) context.Context {
	logger.Infoln("loading database")

	client := database.NewMysqlClient(ctx)

	return context.WithValue(ctx, "db", client)
}

func loadModules(ctx context.Context) context.Context {
	logger.Infoln("loading modules")

	var modules []module.Detecter

	modules = append(modules, authorizedDetector.NewAuthorizedDetector(ctx))
	modules = append(modules, csrfDetector.NewCsrfDetector(ctx))

	return context.WithValue(ctx, "modules", modules)
}

func loadFilter(ctx context.Context) context.Context { //only support single filter currently
	logger.Infoln("loading filters")
	var filters []filter.Filter

	filters = append(filters, filter.NewDuplicateFilter())
	filters = append(filters, filter.NewStaticFileFilter(ctx))
	filters = append(filters, filter.NewHttpFilter())

	return context.WithValue(ctx, "filters", filters)
}

func loadConfigJsonStr(ctx context.Context) context.Context {
	logger.Infoln("loading config into config string")

	bytes, err := os.ReadFile("config.json")

	if err != nil {
		logger.Errorln("read config file error ", err)
		panic(err)
	}

	return context.WithValue(ctx, "configJsonString", string(bytes))
}
