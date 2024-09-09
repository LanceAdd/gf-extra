package ops

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
)

func doHandleFieldOpsPrefixColumnsToMap(prefix string, columns any) *map[string]string {
	m := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &m
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeFieldOpsColumnsMap(cacheKey, &m, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.PkgPath(), of.Name())
		doMergeFieldOpsColumnsMap(cacheKey, &m, prefix, columns)
	}
	return &m
}

func doMergeFieldOpsColumnsMap(cacheKey string, data *map[string]string, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		for _, value := range m {
			(*data)[doHandlePrefix(prefix, value)] = doHandleReal(value)
		}
	} else {
		columnsMap := gconv.MapStrStr(columns)
		cache.Store(cacheKey, columnsMap)
		for _, value := range columnsMap {
			(*data)[doHandlePrefix(prefix, value)] = doHandleReal(value)
		}
	}
}

func doHandleFieldOpsObjectToMap(prefix string, columns any) *map[string]string {
	m := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &m
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeFieldOpsObjectMap(cacheKey, &m, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.PkgPath(), of.Name())
		doMergeFieldOpsObjectMap(cacheKey, &m, prefix, columns)
	}
	return &m
}

func doMergeFieldOpsObjectMap(cacheKey string, data *map[string]string, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		for key := range m {
			(*data)[doHandlePrefix(prefix, key)] = doHandleReal(key)
		}
	} else {
		m, _ := gstructs.TagMapName(columns, []string{"orm", "json"})
		cache.Store(cacheKey, m)
		for key := range m {
			(*data)[doHandlePrefix(prefix, key)] = doHandleReal(key)
		}
	}
}

func doHandleFieldOpsPrefixSliceToMap(prefix string, columns any) *map[string]string {
	m := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &m
	}
	switch field := columns.(type) {
	case []string:
		for _, value := range field {
			m[fmt.Sprintf("%s.%s", prefix, value)] = value
		}
	case *[]string:
		for _, value := range *field {
			m[fmt.Sprintf("%s.%s", prefix, value)] = value
		}
	case gset.StrSet:
		field.Iterator(func(v string) bool {
			m[fmt.Sprintf("%s.%s", prefix, v)] = v
			return true
		})
	case *gset.StrSet:
		field.Iterator(func(v string) bool {
			m[fmt.Sprintf("%s.%s", prefix, v)] = v
			return true
		})
	}
	return &m
}

func doHandleFieldOpsPrefixMapToMap(prefix string, columns any, aliasIsKey ...bool) *map[string]string {
	m := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &m
	}
	var reverse bool
	if len(aliasIsKey) > 0 {
		reverse = !aliasIsKey[0]
	} else {
		reverse = false
	}
	tmpMap := make(map[string]string)
	switch field := columns.(type) {
	case map[string]string:
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&field))
		} else {
			doMergeMap(&tmpMap, &field)
		}

	case *map[string]string:
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(field))
		} else {
			doMergeMap(&tmpMap, field)
		}
	case gmap.StrStrMap:
		m2 := field.Map()
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&m2))
		} else {
			doMergeMap(&tmpMap, &m2)
		}
	case *gmap.StrStrMap:
		m2 := field.Map()
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&m2))
		} else {
			doMergeMap(&tmpMap, &m2)
		}
	}
	for key, value := range tmpMap {
		m[fmt.Sprintf("%s.%s", prefix, value)] = key
	}
	return &m
}

func doHandleFieldOpsMapToMap(columns any, aliasIsKey ...bool) *map[string]string {
	m := make(map[string]string)
	if columns == nil {
		return &m
	}
	var reverse bool
	if len(aliasIsKey) > 0 {
		reverse = !aliasIsKey[0]
	} else {
		reverse = false
	}
	tmpMap := make(map[string]string)
	switch field := columns.(type) {
	case map[string]string:
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&field))
		} else {
			doMergeMap(&tmpMap, &field)
		}

	case *map[string]string:
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(field))
		} else {
			doMergeMap(&tmpMap, field)
		}
	case gmap.StrStrMap:
		m2 := field.Map()
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&m2))
		} else {
			doMergeMap(&tmpMap, &m2)
		}
	case *gmap.StrStrMap:
		m2 := field.Map()
		if reverse {
			doMergeMap(&tmpMap, doSwapMapKeyValue(&m2))
		} else {
			doMergeMap(&tmpMap, &m2)
		}
	}
	if len(tmpMap) == 0 {
		return &m
	}
	return &m
}
