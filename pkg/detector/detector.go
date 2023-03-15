package detector

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/detector/authorize"
	"APIKiller/pkg/detector/bypass"
	"APIKiller/pkg/types"
)

type Detector interface {
	Detect(item *types.DataItem)
}

func NewDetectors(cfg *config.Config) []Detector {
	var detectors []Detector

	detectors = append(detectors, authorize.NewUnauthorizedDetector(cfg))
	detectors = append(detectors, bypass.New40xBypassDetector(cfg))

	return detectors
}
