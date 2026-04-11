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
	Type       string `json:"type" yaml:"type"`
	DbConnStr  string `json:"db_conn_str" yaml:"db_conn_str"`
	DbUsername string `json:"db_username" yaml:"db_username"`
	DbPassword string `json:"db_password" yaml:"db_password"`
}

// ConfigTable defines the structure of Table specific overrides
type ConfigTable struct {
	Name    string              `json:"name" yaml:"name"`
	Method  string              `json:"method" yaml:"method"`
	Columns []map[string]string `json:"columns" yaml:"columns"`
}
