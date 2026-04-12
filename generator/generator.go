package generator

import (
	"context"

	"kumarvv.com/mockdata/constants/targettypes"
	"kumarvv.com/mockdata/models"
)

func Generate(ctx context.Context, config *models.Config) error {
	if config.Target.Type == targettypes.JSON {
		return generateJSON(ctx, config)
	}

	return nil
}
