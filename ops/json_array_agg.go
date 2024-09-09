package ops

import (
	"fmt"
	"strings"
)

type JsonArrayAggOpsImpl struct {
	item string
}

func JsonArrayAggOps() *JsonArrayAggOpsImpl {
	return &JsonArrayAggOpsImpl{}
}

func (j *JsonArrayAggOpsImpl) Include(field string) *JsonArrayAggOpsImpl {
	j.item = strings.TrimSpace(field)
	return j
}

func (j *JsonArrayAggOpsImpl) Items() string {
	return strings.TrimSpace(j.item)
}

func (j *JsonArrayAggOpsImpl) As(alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_ARRAYAGG(%s) AS `%s`", j.Items(), alias[0])
	}
	return fmt.Sprintf("JSON_ARRAYAGG(%s)", j.Items())
}
