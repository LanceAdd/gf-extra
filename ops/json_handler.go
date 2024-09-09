package ops

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func doHandleJsonObjectOpsPrefixColumnsToMap(prefix string, columns any) *map[string]string {
	dataMap := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &dataMap
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeJsonObjectOpsColumnsMap(cacheKey, &dataMap, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.PkgPath(), of.Name())
		doMergeJsonObjectOpsColumnsMap(cacheKey, &dataMap, prefix, columns)
	}
	return &dataMap
}

func doMergeJsonObjectOpsColumnsMap(cacheKey string, dataMap *map[string]string, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		prefixMap := doHandleMapValuePrefix(prefix, m)
		doMergeMap(dataMap, prefixMap)
	} else {
		columnsMap := gconv.MapStrStr(columns)
		cache.Store(cacheKey, columnsMap)
		prefixMap := doHandleMapValuePrefix(prefix, columnsMap)
		doMergeMap(dataMap, prefixMap)
	}
}

func doHandleJsonObjectOpsPrefixMapToMap(prefix string, columns any) *map[string]string {
	dataMap := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &dataMap
	}
	switch field := columns.(type) {
	case map[string]string:
		for key, value := range field {
			field[key] = doHandlePrefix(prefix, value)
		}
		doMergeMap(&dataMap, &field)
	case *map[string]string:
		for key, value := range *field {
			(*field)[key] = doHandlePrefix(prefix, value)
		}
		doMergeMap(&dataMap, field)
	case gmap.StrStrMap:
		field.Iterator(func(k string, v string) bool {
			dataMap[k] = doHandlePrefix(prefix, v)
			return true
		})
	case *gmap.StrStrMap:
		field.Iterator(func(k string, v string) bool {
			dataMap[k] = doHandlePrefix(prefix, v)
			return true
		})
	}
	return &dataMap
}

func doHandleJsonObjectOpsMapToMap(columns any) *map[string]string {
	dataMap := make(map[string]string)
	if columns == nil {
		return &dataMap
	}
	switch field := columns.(type) {
	case map[string]string:
		doMergeMap(&dataMap, &field)
	case *map[string]string:
		doMergeMap(&dataMap, field)
	case gmap.StrStrMap:
		field.Iterator(func(k string, v string) bool {
			dataMap[k] = v
			return true
		})
	case *gmap.StrStrMap:
		field.Iterator(func(k string, v string) bool {
			dataMap[k] = v
			return true
		})
	}
	return &dataMap
}

func doHandleJsonObjectOpsPrefixSliceToMap(prefix string, columns any) *map[string]string {
	dataMap := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &dataMap
	}
	switch field := columns.(type) {
	case []string:
		for _, value := range field {
			dataMap[gstr.CaseCamelLower(value)] = doHandlePrefix(prefix, value)
		}
	case *[]string:
		for _, value := range *field {
			dataMap[gstr.CaseCamelLower(value)] = doHandlePrefix(prefix, value)
		}
	case gset.StrSet:
		field.Iterator(func(v string) bool {
			dataMap[gstr.CaseCamelLower(v)] = doHandlePrefix(prefix, v)
			return true
		})
	case *gset.StrSet:
		field.Iterator(func(v string) bool {
			dataMap[gstr.CaseCamelLower(v)] = doHandlePrefix(prefix, v)
			return true
		})
	}
	return &dataMap
}
func doHandleJsonObjectOpsObjectToMap(prefix string, columns any) *map[string]string {
	dataMap := make(map[string]string)
	if IsBlank(prefix) || columns == nil {
		return &dataMap
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeJsonObjectOpsObjectMap(cacheKey, &dataMap, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.PkgPath(), of.Name())
		doMergeJsonObjectOpsObjectMap(cacheKey, &dataMap, prefix, columns)
	}
	return &dataMap
}

func doMergeJsonObjectOpsObjectMap(cacheKey string, dataMap *map[string]string, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		prefixMap := doHandleMapValuePrefix(prefix, m)
		doMergeMap(dataMap, prefixMap)
	} else {
		tagMap, _ := gstructs.TagMapName(columns, []string{"orm", "json"})
		m := doSwapMapKeyValue(&tagMap)
		cache.Store(cacheKey, *m)
		prefixMap := doHandleMapValuePrefix(prefix, *m)
		doMergeMap(dataMap, prefixMap)
	}
}
