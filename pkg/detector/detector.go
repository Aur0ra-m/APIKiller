package detector

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/detector/authorize"
	"APIKiller/pkg/detector/bypass"
	"APIKiller/pkg/detector/csrf"
	"APIKiller/pkg/detector/openredirect"
	"APIKiller/pkg/types"
)

type Detector interface {
	Detect(item *types.DataItem)
}

func NewDetectors(cfg *config.DetectorConfig) []Detector {
	var detectors []Detector

	detectors = append(detectors, authorize.NewUnauthorizedDetector(cfg))
	detectors = append(detectors, bypass.New40xBypassDetector(cfg))
	detectors = append(detectors, csrf.NewCsrfDetector(cfg))
	detectors = append(detectors, openredirect.NewOpenRedirectDetector(cfg))

	return detectors
}
