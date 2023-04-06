package module

import (
	"APIKiller/core/data"
)

const (
	AsyncDetectVulnTypeSeperator = "^"
)

var (
	Modules []Detecter
)

type Detecter interface {
	Detect(item *data.DataItem) (result *data.DataItem)
}

func RegisterModule(d Detecter) {
	if d == nil {
		return
	}

	Modules = append(Modules, d)
}
