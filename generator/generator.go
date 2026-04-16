package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"
	"kumarvv.com/mockdata/constants/targettypes"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

func Generate(ctx context.Context, config *models.Config) error {
	path := config.Target.ToPath
	targetType := config.Target.Type
	utils.Log("[%s] processing target path: %s", targetType, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create directory path
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create directory path: %s", path)
		}
	}

	utils.Log("[%s] processing tables. count=%d", targetType, len(config.Tables))
	for i, table := range config.Tables {
		tablePath := filepath.Join(path, fmt.Sprintf("%s.%s", table.Name, targetType))
		utils.Log("[%s] processing table(%d of %d): name=%s, tablePath=%s", targetType, i+1, len(config.Tables), table.Name, tablePath)

		if table.RowCount == 0 {
			table.RowCount = 1
		}
		utils.Log("generating rows: %d", table.RowCount)
		logMarker := 0
		rows := make([]map[string]interface{}, 0)
		for r := 0; r < table.RowCount; r++ {
			row := make(map[string]interface{})
			gender := utils.RandomOneOf(randomdata.Male, randomdata.Female)
			for _, column := range table.Columns {
				if value, err := generateValue(ctx, &table, &column, gender, r); err != nil {
					return errors.Wrapf(err, "failed to generate value for table(%d of %d): name=%s, column=%s", i+1,
						len(config.Tables), table.Name, column.Name)
				} else {
					row[column.Name] = value
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
		utils.Log("generating %s data", targetType)
		var b []byte
		var err error
		if targetType == targettypes.JSON {
			if b, err = generateJSON(ctx, rows); err != nil {
				return errors.Wrapf(err, "failed to generate json data: %s", tablePath)
			}
		} else if targetType == targettypes.SQL {
			if b, err = generateSQL(ctx, &table, rows); err != nil {
				return errors.Wrapf(err, "failed to generate sql data: %s", tablePath)
			}
		}

		utils.Log("creating %s file at tablePath=%s", targetType, tablePath)
		if err = os.WriteFile(tablePath, b, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to write file: %s", tablePath)
		}
		utils.Log("[%s] data file created at %s", tablePath, targetType)
	}

	return nil
}
