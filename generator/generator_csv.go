package generator

import (
	"bytes"
	"context"
	"encoding/csv"

	"github.com/pkg/errors"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

func generateCSV(ctx context.Context, table *models.ConfigTable, rows []map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// header row
	headers := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		headers[i] = col.Name
	}
	if err := w.Write(headers); err != nil {
		return nil, errors.Wrap(err, "failed to write csv header")
	}

	// data rows
	for _, row := range rows {
		record := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			record[i] = utils.ToString(row[col.Name])
		}
		if err := w.Write(record); err != nil {
			return nil, errors.Wrap(err, "failed to write csv row")
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, errors.Wrap(err, "csv flush error")
	}

	return buf.Bytes(), nil
}
