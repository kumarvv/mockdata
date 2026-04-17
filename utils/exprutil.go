package utils

import (
	"fmt"
	"strings"

	"kumarvv.com/mockdata/models"
)

func ParseValueExpr(expr string) (*models.Column, error) {

	return nil, nil
}

func getExprTokens(expr string) (string, map[string]string, error) {
	expr = strings.TrimSpace(expr)

	// fn start
	fs := strings.Index(expr, "(")
	if fs == -1 {
		return "", nil, fmt.Errorf("function expression required: %s", expr)
	}
	// fn end
	fe := strings.Index(expr, ")")
	if fe == -1 {
		return "", nil, fmt.Errorf("valid function expression required - missing ')' : %s", expr)
	}

	// fn name
	fnName := expr[:fs]
	if IsBlank(fnName) {
		return "", nil, fmt.Errorf("function name required: %s", expr)
	}

	// params
	paramsExpr := expr[fs+1 : fe]
	paramsKVs := strings.Split(paramsExpr, ",")
	params := make(map[string]string)
	for _, paramKV := range paramsKVs {
		items := strings.Split(paramKV, "=")
		kv := strings.TrimSpace(items[0])
		if len(items) > 1 {
			params[kv] = strings.TrimSpace(items[1])
		}
	}

	return fnName, params, nil
}
