package ops

import (
	"fmt"
	"strings"
)

type JsonObjectOpsImpl struct {
	itemMap map[string]string
}

func JsonObjectOps() *JsonObjectOpsImpl {
	return &JsonObjectOpsImpl{itemMap: make(map[string]string)}
}

func (j *JsonObjectOpsImpl) IncludePrefixColumns(prefixOrKey string, columns any) *JsonObjectOpsImpl {
	m := doHandleJsonObjectOpsPrefixColumnsToMap(prefixOrKey, columns)
	doMergeMap(&j.itemMap, m)
	return j
}

func (j *JsonObjectOpsImpl) IncludePrefixObject(prefixOrKey string, columns any) *JsonObjectOpsImpl {
	m := doHandleJsonObjectOpsObjectToMap(prefixOrKey, columns)
	doMergeMap(&j.itemMap, m)
	return j
}

func (j *JsonObjectOpsImpl) IncludePrefixMap(prefixOrKey string, columns any) *JsonObjectOpsImpl {
	m := doHandleJsonObjectOpsPrefixMapToMap(prefixOrKey, columns)
	doMergeMap(&j.itemMap, m)
	return j
}

func (j *JsonObjectOpsImpl) IncludeMap(columns any) *JsonObjectOpsImpl {
	m := doHandleJsonObjectOpsMapToMap(columns)
	doMergeMap(&j.itemMap, m)
	return j
}

func (j *JsonObjectOpsImpl) IncludePrefixSlice(prefixOrKey string, columns any) *JsonObjectOpsImpl {
	m := doHandleJsonObjectOpsPrefixSliceToMap(prefixOrKey, columns)
	doMergeMap(&j.itemMap, m)
	return j
}

func (j *JsonObjectOpsImpl) IncludeItem(key string, value string) *JsonObjectOpsImpl {
	if IsBlank(key) || IsBlank(value) {
		return j
	}
	j.itemMap[key] = value
	return j
}
func (j *JsonObjectOpsImpl) Exclude(keys ...string) *JsonObjectOpsImpl {
	if len(keys) == 0 {
		return j
	}
	for _, key := range keys {
		delete(j.itemMap, key)
	}
	return j
}

func (j *JsonObjectOpsImpl) Items() []string {
	var data []string
	for key, value := range j.itemMap {
		data = append(data, fmt.Sprintf("'%s', %s", key, value))
	}
	return data
}

func (j *JsonObjectOpsImpl) Omitempty() *JsonObjectOpsImpl {
	delete(j.itemMap, "")
	var keys []string
	for key, value := range j.itemMap {
		if IsBlank(key) || IsBlank(value) {
			keys = append(keys, key)
		}
	}
	for _, key := range keys {
		delete(j.itemMap, key)
	}
	return j
}

func (j *JsonObjectOpsImpl) As(alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_OBJECT(%s) AS `%s`", strings.Join(j.Omitempty().Items(), ", "), alias[0])

	}
	return fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(j.Omitempty().Items(), ", "))

}
