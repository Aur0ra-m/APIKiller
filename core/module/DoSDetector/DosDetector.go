package DoSDetector

import (
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/log"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type DosDetector struct {
	d1 *rateLimitDetector
	d2 *resourceSizeDetector
}

func (d DosDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] DoS detect")

	// rate limit
	//d.d1.Detect(ctx, item)

	// the size of resource lack of control
	d.d2.Detect(ctx, item)
	
}

func NewDosDetector(ctx context.Context) module.Detecter {
	if viper.GetInt("app.module.DoSDetector.option") == 0 {
		return nil
	}

	logger.Infoln("[Load Module] DoS detect module")

	return &DosDetector{
		d1: newRateLimitDetector(),
		d2: newResourceSizeDetector(),
	}
}
