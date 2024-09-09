package ops

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
)

type JsonArrayOpsImpl struct {
	items []string
}

func JsonArrayOps() *JsonArrayOpsImpl {
	return &JsonArrayOpsImpl{}
}
func (j *JsonArrayOpsImpl) Include(data any) *JsonArrayOpsImpl {
	encode, _ := gjson.EncodeString(data)
	j.items = append(j.items, fmt.Sprintf("%s", encode))
	return j
}
func (j *JsonArrayOpsImpl) IncludeColumn(column string) *JsonArrayOpsImpl {
	j.items = append(j.items, column)
	return j
}
func (j *JsonArrayOpsImpl) Items() string {
	return strings.Join(j.items, ", ")
}

func (j *JsonArrayOpsImpl) As(alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_ARRAY(%s) AS `%s`", j.Items(), alias[0])
	}
	return fmt.Sprintf("JSON_ARRAY(%s)", j.Items())
}
