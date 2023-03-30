package dos

import (
	"APIKiller/pkg/logger"
	"APIKiller/pkg/types"
)

type DosDetector struct {
	d1 *rateLimitDetector
	d2 *resourceSizeDetector
}

func (d DosDetector) Detect(item *types.DataItem) {
	logger.Debug("[Detect] DoS detect\n")

	// rate limit
	//d.d1.Detect(ctx, item)

	// the size of resource lack of control
	d.d2.Detect(item)

}
