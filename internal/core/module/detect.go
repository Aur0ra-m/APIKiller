package module

import (
	"APIKiller/internal/core/data"
)

const (
	AsyncDetectVulnTypeSeperator = "^"
)

var (
	modules []Detecter
)

type Detecter interface {
	//
	// Detect
	//  @Description: detect the target api and return the result
	//  @param item
	//  @return result
	//
	Detect(item *data.DataItem) (result *data.DataItem)
}

func RegisterModule(d Detecter) {
	if d == nil {
		return
	}

	modules = append(modules, d)
}

func GetModules() []Detecter {
	if modules != nil {
		return modules
	}
	return nil
}
