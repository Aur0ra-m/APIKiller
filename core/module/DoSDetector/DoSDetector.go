package DoSDetector

import (
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/logger"
	"github.com/spf13/viper"
)

type DosDetector struct {
	typeFlag string
	d1       *rateLimitDetector
	d2       *resourceSizeDetector
}

func (d DosDetector) Detect(item *data.DataItem) (result *data.DataItem) {
	logger.Debugln("[Detect] DoS detect")

	// rate limit
	//d.d1.Detect( item)

	// the size of resource lack of control
	return d.d2.Detect(item)
}

func NewDoSDetector() module.Detecter {
	if viper.GetInt("app.module.DoSDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] DoS detect module")

	return &DosDetector{
		typeFlag: viper.GetString("app.module.DoSDetector.typeFlag"),
		d1:       newRateLimitDetector(),
		d2:       newResourceSizeDetector(),
	}
}
