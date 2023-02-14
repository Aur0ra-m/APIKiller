package notify

import (
	"APIKiller/core/data"
	"context"
)

type Notify interface {
	Notify(ctx context.Context, item *data.DataItem)
}
