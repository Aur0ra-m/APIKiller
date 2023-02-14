package authorizedDetector

import (
	"APIKiller/core/data"
	"APIKiller/core/module"
	logger "APIKiller/log"
	"APIKiller/util"
	"context"
)

type AuthorizedDetector struct {
	multiRolesDetector   module.Detecter
	singleRoleDetector   module.Detecter
	unauthorizedDetector module.Detecter
}

func (d *AuthorizedDetector) Detect(ctx context.Context, item *data.DataItem) {
	logger.Debugln("[Detect] main  authorized detect")

	d.unauthorizedDetector.Detect(ctx, item)
	for _, t := range item.VulnType {
		if t == "unauthorized" {
			return
		}
	}

	//if d.singleRoleDetector != nil {
	//	d.singleRoleDetector.Detect(ctx, item)
	//}

	if d.multiRolesDetector != nil {
		d.multiRolesDetector.Detect(ctx, item)
	}

}

func NewAuthorizedDetector(ctx context.Context) module.Detecter {
	if util.GetConfig(ctx, "app.detectors.authorizedDetector.option") != "1" {
		return nil
	}

	logger.Infoln("[Load Module] authorized module")

	var detector *AuthorizedDetector

	if util.GetConfig(ctx, "app.detectors.authorizedDetector.multiRolesDetector.role") != "" {
		detector = &AuthorizedDetector{
			multiRolesDetector:   newMultiRolesDetector(ctx),
			singleRoleDetector:   nil,
			unauthorizedDetector: newUnauthorizedDetector(ctx),
		}
	} else {
		detector = &AuthorizedDetector{
			multiRolesDetector:   nil,
			singleRoleDetector:   newSingleRoleDetector(ctx),
			unauthorizedDetector: newUnauthorizedDetector(ctx),
		}
	}

	return detector
}
