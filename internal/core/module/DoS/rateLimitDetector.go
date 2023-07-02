package DoS

import (
	"APIKiller/internal/core/data"
)

type rateLimitDetector struct {
}

func (r *rateLimitDetector) Detect(item *data.DataItem) {
	//
}

func newRateLimitDetector() *rateLimitDetector {
	return &rateLimitDetector{}
}
