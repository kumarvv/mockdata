package core

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
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

func LoadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to read file %s", "./dbaudit.yaml")
		return nil, err
	}

	var config models.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal file %s", "./dbaudit.yaml")
		return nil, err
	}

	fillEnvVars(&config)

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
