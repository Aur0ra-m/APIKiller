package dos

import "APIKiller/pkg/types"

type rateLimitDetector struct {
}

func (r *rateLimitDetector) Detect(item *types.DataItem) {
	//
}

func newRateLimitDetector() *rateLimitDetector {
	return &rateLimitDetector{}
}
