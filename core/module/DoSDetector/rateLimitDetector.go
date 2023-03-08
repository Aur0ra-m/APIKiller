package DoSDetector

import (
	"APIKiller/core/data"
	"context"
)

type rateLimitDetector struct {
}

func (r *rateLimitDetector) Detect(ctx context.Context, item *data.DataItem) {
	//
}

func newRateLimitDetector() *rateLimitDetector {
	return &rateLimitDetector{}
}
