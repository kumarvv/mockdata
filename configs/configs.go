package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"kumarvv.com/mockdata/constants/dbtypes"
	"kumarvv.com/mockdata/constants/functiontypes"
	"kumarvv.com/mockdata/constants/tablemodes"
	"kumarvv.com/mockdata/constants/targettypes"
	"kumarvv.com/mockdata/utils"

	"kumarvv.com/mockdata/models"
)

//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// 2022 (c) Vijay Vijayaram
//

func Load(path string) (*models.Config, []error) {
	data, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to read file %s", path)
		return nil, []error{err}
	}

	var config models.Config
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal file %s", path)
		return nil, []error{err}
	}

	fillEnvVars(&config)

	if errs := validate(&config); len(errs) > 0 {
		return nil, errs
	}

	return &config, nil
}

func fillEnvVars(config *models.Config) {
	if utils.IsBlank(config.Target.DbConnStr) {
		return
	}

	str := config.Target.DbConnStr
	tokens := make([]string, 0)
	for {
		startx := strings.Index(str, "%%")
		if startx <= 0 {
			break
		}
		startx += 2

		rstr := str[startx:]
		endx := strings.Index(rstr, "%%")
		if endx <= 0 {
			break
		}

		tokens = append(tokens, str[startx:startx+endx])
		str = str[startx+endx+2:]
	}

	for _, token := range tokens {
		value := os.Getenv(token)
		config.Target.DbConnStr = strings.ReplaceAll(config.Target.DbConnStr, "%%"+token+"%%", value)
	}
}

func validate(config *models.Config) []error {
	if config == nil {
		return []error{errors.New("config is nil")}
	}

	errs := make([]error, 0)

	// target type
	if utils.IsBlank(config.Target.Type) {
		config.Target.Type = targettypes.SQL
	}
	if !utils.Includes(targettypes.List(), config.Target.Type) {
		errs = append(errs, errors.Errorf("invalid target type %s", config.Target.Type))
	}
	// target type: db
	if config.Target.Type == targettypes.DB {
		if utils.IsBlank(config.Target.DbType) {
			errs = append(errs, errors.New("db_type is required"))
		} else if !utils.Includes(dbtypes.List(), config.Target.DbType) {
			errs = append(errs, errors.Errorf("invalid db_type %s", config.Target.DbType))
		}
		if utils.IsBlank(config.Target.DbConnStr) {
			errs = append(errs, errors.New("db_conn_str is required"))
		}
	} else {
		if utils.IsBlank(config.Target.ToPath) {
			errs = append(errs, errors.New("to_path is required"))
		}
	}

	// tables
	if len(config.Tables) == 0 {
		errs = append(errs, errors.New("at least one table in `tables` key is required"))
	}
	for _, table := range config.Tables {
		if utils.IsBlank(table.Name) {
			errs = append(errs, errors.Errorf("table name is required for table %s", table.Name))
		}
		if utils.IsBlank(table.Mode) {
			table.Mode = tablemodes.Append
		}
		if !utils.Includes(tablemodes.List(), table.Mode) {
			errs = append(errs, errors.Errorf("invalid table mode %s for table %s", table.Mode, table.Name))
		}

		table.Columns = make([]*models.Column, 0)
		for _, columnMap := range table.RawColumns {
			for columnName, valueExpr := range columnMap {
				if column, err := parseValueExpr(columnName, valueExpr); err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to parse value expression for table.column %s.%s",
						table.Name, columnName))
				} else {
					column.Name = columnName
					table.Columns = append(table.Columns, column)
				}
			}
		}
	}

	return errs
}

func parseValueExpr(columnName, expr string) (*models.Column, error) {
	expr = strings.TrimSpace(expr)

	// fn start
	fs := strings.Index(expr, "(")
	if fs == -1 {
		return nil, fmt.Errorf("function expression required: %s", expr)
	}
	// fn end
	fe := strings.Index(expr, ")")
	if fe == -1 {
		return nil, fmt.Errorf("valid function expression required - missing ')' : %s", expr)
	}

	// fn name
	fnName := expr[:fs]
	if utils.IsBlank(fnName) {
		return nil, fmt.Errorf("function name required")
	}
	if !utils.Includes(functiontypes.List(), fnName) {
		return nil, fmt.Errorf("invalid function name [%s]", fnName)
	}

	// params
	params := make(map[string]string)
	paramsExpr := expr[fs+1 : fe]
	if !utils.IsBlank(paramsExpr) {
		paramsKVs := strings.Split(paramsExpr, ",")
		for _, paramKV := range paramsKVs {
			items := strings.Split(paramKV, "=")
			if len(items) > 1 {
				paramKey := strings.TrimSpace(items[0])
				if !utils.Includes(functiontypes.GetParams(fnName), paramKey) {
					return nil, fmt.Errorf("invalid param key [%s] for function [%s]", paramKey, fnName)
				}
				params[paramKey] = strings.TrimSpace(items[1])
			} else {
				// simple value becomes key "value"
				params["value"] = strings.TrimSpace(items[0])
			}
		}
	}

	column := buildColumn(columnName, fnName, params)

	// value param required
	if column.Value == nil && functiontypes.IsRequiredValueExpr(fnName) {
		return nil, fmt.Errorf("value param required: %s", expr)
	}

	return &column, nil
}

func buildColumn(columnName, fnName string, params map[string]string) models.Column {
	// column
	column := models.Column{
		Name:   columnName,
		FnName: fnName,
	}

	for k, v := range params {
		if k == "value" {
			column.Value = v
		} else if k == "len" {
			column.Len = utils.IntPtr(utils.ToInt(v))
		} else if k == "min" {
			column.Min = utils.IntPtr(utils.ToInt(v))
		} else if k == "max" {
			column.Max = utils.IntPtr(utils.ToInt(v))
		} else if k == "format" {
			column.Format = utils.StrPtr(utils.ToString(v))
		} else if k == "case" {
			column.Case = utils.StrPtr(utils.ToString(v))
		} else if k == "numpairs" {
			column.NumPairs = utils.IntPtr(utils.ToInt(v))
		} else if k == "separator" {
			column.Separator = utils.StrPtr(utils.ToString(v))
		}
	}

	return column
}
