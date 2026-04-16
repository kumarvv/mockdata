package generator

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func generateJSON(ctx context.Context, rows []map[string]interface{}) ([]byte, error) {
	if b, err := json.MarshalIndent(rows, "", "  "); err != nil {
		return nil, errors.Wrapf(err, "failed to marshal data")
	} else {
		return b, nil
	}
}
