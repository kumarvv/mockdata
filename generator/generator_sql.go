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
	columnNames := generateSQLColumns(table)

	for _, row := range rows {
		if insertSql, err := generateSQLInsert(table, row, columnNames); err != nil {
			return nil, err
		} else {
			contents += fmt.Sprintf("%s\n", insertSql)
		}
	}

	return []byte(contents), nil
}

func generateSQLColumns(table *models.ConfigTable) string {
	columnNames := make([]string, 0)
	for _, column := range table.Columns {
		columnNames = append(columnNames, strings.ToLower(column.Name))
	}
	return strings.Join(columnNames, ", ")
}

func generateSQLInsert(table *models.ConfigTable, row map[string]interface{}, columnNames string) (string, error) {
	values := make([]string, 0)
	for _, column := range table.Columns {
		value := row[column.Name]
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

	sql := fmt.Sprintf(insertTemplate, strings.ToLower(table.Name), columnNames, allValues)
	return sql, nil
}
