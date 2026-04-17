package models

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

// Config defines the structure of Config data
type Config struct {
	Target ConfigTarget  `json:"target" yaml:"target"`
	Tables []ConfigTable `json:"tables" yaml:"tables"`
}

type ConfigTarget struct {
	Type      string `json:"type" yaml:"type"`
	DbType    string `json:"db_type" yaml:"db_type"`
	DbConnStr string `json:"db_conn_str" yaml:"db_conn_str"`
	ToPath    string `json:"to_path" yaml:"to_path"`
}

// ConfigTable defines the structure of Table specific overrides
type ConfigTable struct {
	Name       string              `json:"name" yaml:"name"`
	Mode       string              `json:"mode" yaml:"mode"`
	RowCount   int                 `json:"row_count" yaml:"row_count"`
	RawColumns []map[string]string `json:"columns" yaml:"columns"`
	Columns    []*Column           `json:"-" yaml:"-"`
}

type Column struct {
	Name string `json:"name" yaml:"name"`

	FnName string `json:"fn_name" yaml:"fn_name"`

	Value interface{} `json:"value" yaml:"value"`

	// params
	Len       *int    `json:"len" yaml:"len"`
	Min       *int    `json:"min" yaml:"min"`
	Max       *int    `json:"max" yaml:"max"`
	Format    *string `json:"format" yaml:"format"`
	Case      *string `json:"case" yaml:"case"`
	NumPairs  *int    `json:"num_pairs" yaml:"num_pairs"`
	Separator *string `json:"separator" yaml:"separator"`
}
