package configs

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"kumarvv.com/mockdata/constants/dbtypes"
	"kumarvv.com/mockdata/constants/tablemodes"
	"kumarvv.com/mockdata/constants/targettypes"
	"kumarvv.com/mockdata/constants/valuetypes"
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
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal file %s", path)
		return nil, []error{err}
	}

	fillEnvVars(&config)

	if errs := Validate(&config); len(errs) > 0 {
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

func Validate(config *models.Config) []error {
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

		if len(table.Columns) == 0 {
			errs = append(errs, errors.Errorf("at least one column is required for table %s", table.Name))
		} else {
			for _, colMap := range table.Columns {
				for column, value := range colMap {
					tokens := strings.Split(value, "|")
					valueType := tokens[0]
					valueExpr := ""
					if len(tokens) > 1 {
						valueExpr = tokens[1]
					}
					if !utils.Includes(valuetypes.List(), valueType) {
						errs = append(errs, errors.Errorf("invalid value type %s for table.column %s.%s",
							table.Name, column, valueType))
					} else if valuetypes.IsRequiredValueExpr(valueType) && utils.IsBlank(valueExpr) {
						errs = append(errs, errors.Errorf("value expression is required for table.column %s.%s",
							table.Name, column))
					}
					break
				}
			}
		}
	}

	return errs
}
