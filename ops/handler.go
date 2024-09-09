package ops

import (
	"fmt"
	"sort"
	"strings"
)

func doHandleMapValuePrefix(prefix string, m map[string]string) *map[string]string {
	data := make(map[string]string)
	for key, value := range m {
		data[key] = doHandlePrefix(prefix, value)
	}
	return &data
}

func doHandlePrefix(prefix string, value string) string {
	return fmt.Sprintf("`%s`.`%s`", prefix, value)
}

func doHandleReal(value string) string {
	return fmt.Sprintf("`%s`", value)
}

func doSwapMapKeyValue(m *map[string]string) *map[string]string {
	data := make(map[string]string)
	for key, value := range *m {
		data[value] = key
	}
	return &data
}

func doMergeMap(source *map[string]string, target *map[string]string) {
	for k, v := range *target {
		(*source)[k] = v
	}
}

func IsBlank(prefixOrKey string) bool {
	return strings.TrimSpace(prefixOrKey) == ""
}

func doCopySlice(source []string) []string {
	copyFields := make([]string, len(source))
	copy(copyFields, source)
	return copyFields
}

func doSortFields(source []string) {
	sort.SliceStable(source, func(i, j int) bool {
		// 先忽略大小写排序
		s1, s2 := strings.ToLower(source[i]), strings.ToLower(source[j])
		if len(s1) != len(s2) {
			return len(s1) < len(s2) // 根据长度排序
		}
		return s1 < s2
	})
}

func doBoolDefaultTrue(b ...bool) bool {
	var res bool
	if len(b) > 0 {
		res = b[0]
	} else {
		res = true
	}
	return res
}

func doBoolDefaultFalse(b ...bool) bool {
	var res bool
	if len(b) > 0 {
		res = b[0]
	} else {
		res = false
	}
	return res
}
