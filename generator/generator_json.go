package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	orderedmap "github.com/wk8/go-ordered-map"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

func generateJSON(ctx context.Context, config *models.Config) error {
	path := config.Target.ToPath
	utils.Log("[JSON] processing target path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create directory path
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create directory path: %s", path)
		}
	}

	utils.Log("[JSON] processing tables. count=%d", len(config.Tables))
	for i, table := range config.Tables {
		tablePath := filepath.Join(path, fmt.Sprintf("%s.json", table.Name))
		utils.Log("[JSON] processing table(%d of %d): name=%s, filePath=%s", i+1, len(config.Tables), table.Name, tablePath)

		utils.Log("generating rows: %d", table.RowCount)
		logMarker := 0
		rows := make([]*orderedmap.OrderedMap, 0)
		for r := 0; r < table.RowCount; r++ {
			row := orderedmap.New[string, interface{}]()
			for _, col := range table.Columns {
				for column, valueExpr := range col {
					if value, err := generateValue(ctx, valueExpr); err != nil {
						return errors.Wrapf(err, "failed to generate value for table(%d of %d): name=%s, column=%s", i+1, len(config.Tables), table.Name, column)
					} else {
						row.Set(column, value)
					}
				}
			}
			rows = append(rows, row)
			logMarker++
			if logMarker%100 == 0 {
				utils.Log("generated rows %d out of %d", logMarker, table.RowCount)
			}
		}
		utils.Log("total rows generated: %d", logMarker)

		// write file
		utils.Log("creating json file at path=%s", tablePath)
		if b, err := json.MarshalIndent(rows, "", "  "); err != nil {
			if err = os.WriteFile(tablePath, b, os.ModePerm); err != nil {
				return errors.Wrapf(err, "failed to write file: %s", tablePath)
			}
		}
		utils.Log("[JSON] data file created at %s", tablePath)
	}

	return nil
}
