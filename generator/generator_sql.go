package generator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"kumarvv.com/mockdata/models"
)

const insertTemplate = `INSER INTO %s 
	(%s) 
VALUES 
	(%s); `

func generateSQL(ctx context.Context, table *models.ConfigTable, rows []map[string]interface{}) ([]byte, error) {
	contents := ""

	columnNames := make([]string, 0)
	for _, column := range table.Columns {
		columnNames = append(columnNames, strings.ToLower(column.Name))
	}
	columnNamesStr := strings.Join(columnNames, ", ")

	for _, row := range rows {
		values := make([]string, 0)
		for _, columnName := range columnNames {
			value := row[columnName]
			valueStr := ""
			if value == nil {
				value = "NULL"
			} else if reflect.TypeOf(value).Kind() == reflect.String {
				valueStr = fmt.Sprintf("'%s'", value)
			} else {
				valueStr = fmt.Sprintf("%v", value)
			}
			values = append(values, valueStr)
		}
		allValues := strings.Join(values, ", ")

		sql := fmt.Sprintf(insertTemplate, strings.ToLower(table.Name), columnNamesStr, allValues)
		contents += fmt.Sprintf("%s\n", sql)
	}

	return []byte(contents), nil
}
