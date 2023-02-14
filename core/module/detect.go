package module

import (
	"APIKiller/core/data"
	"context"
)

type Detecter interface {
	Detect(ctx context.Context, item *data.DataItem)
}
