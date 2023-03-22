package main

import (
	"APIKiller/pkg/config"
	gohttp "APIKiller/pkg/http"
	"APIKiller/pkg/http/hook"
	"APIKiller/pkg/logger"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

const Separator = string(filepath.Separator)

func NewHook(cfg *config.Config) {
	if runtime.GOOS == "windows" {
		logger.Error("not support windows operation system\n")
		return
	}

	pwd, _ := os.Getwd()
	dir := pwd + "/pkg/hooks"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		logger.Error("target directory does not exist\n")
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		return
	}

	entries, err := os.ReadDir(dir)
	for _, entry := range entries {
		soName := entry.Name()
		// filter directory and none so file
		if entry.IsDir() == true || strings.Index(soName, ".so") == -1 {
			continue
		}

		// load plugins and register them via RegisterHooks
		logger.Infof("[Load Hook] load hook %s\n", strings.Replace(soName, ".so", "", 1))
		open, err := plugin.Open(dir + Separator + soName)
		if err != nil {
			logger.Errorf("load hook %s error: %v\n", soName, err)
			panic(err)
		}

		Hook, err := open.Lookup("Hook")
		if err != nil {
			logger.Errorf("load hook %s error: %v\n", soName, err)
			panic(err)
		}

		var greeter hook.RequestHook
		greeter, ok := Hook.(hook.RequestHook)
		if !ok {
			logger.Errorf("load hook %s error: unexpected type from module symbol\n", soName)
			panic(err)
		}

		gohttp.RegisterHooks(greeter)
	}
}

func main() {

}
