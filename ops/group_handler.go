package ops

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
)

func doHandleGroupOpsPrefixColumnsToSet(prefix string, columns any) *gset.StrSet {
	set := gset.NewStrSet()
	if IsBlank(prefix) || columns == nil {
		return set
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeGroupOpsColumnsSet(cacheKey, set, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyColumnsField, of.PkgPath(), of.Name())
		doMergeGroupOpsColumnsSet(cacheKey, set, prefix, columns)
	}
	return set
}

func doMergeGroupOpsColumnsSet(cacheKey string, set *gset.StrSet, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		for _, value := range m {
			set.Add(fmt.Sprintf("%s.%s", prefix, value))
		}
	} else {
		columnsMap := gconv.MapStrStr(columns)
		cache.Store(cacheKey, columnsMap)
		for _, value := range columnsMap {
			set.Add(fmt.Sprintf("%s.%s", prefix, value))
		}
	}
}

func doHandleGroupOpsObjectToSet(prefix string, columns any) *gset.StrSet {
	set := gset.NewStrSet()
	if IsBlank(prefix) || columns == nil {
		return set
	}
	of := reflect.TypeOf(columns)
	if of.Kind() == reflect.Ptr {
		if of.Elem().Kind() == reflect.Struct {
			cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.Elem().PkgPath(), of.Elem().Name())
			doMergeGroupOpsObjectMap(cacheKey, set, prefix, columns)
		}
	}
	if of.Kind() == reflect.Struct {
		cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyObjectField, of.PkgPath(), of.Name())
		doMergeGroupOpsObjectMap(cacheKey, set, prefix, columns)
	}
	return set
}

func doMergeGroupOpsObjectMap(cacheKey string, set *gset.StrSet, prefix string, columns any) {
	if cached, ok := cache.Load(cacheKey); ok {
		m := cached.(map[string]string)
		for key := range m {
			set.Add(fmt.Sprintf("%s.%s", prefix, key))
		}
	} else {
		m, _ := gstructs.TagMapName(columns, []string{"orm", "json"})
		cache.Store(cacheKey, m)
		for key := range m {
			set.Add(fmt.Sprintf("%s.%s", prefix, key))
		}
	}
}

func doHandleGroupOpsPrefixSliceToSet(prefix string, columns any) *gset.StrSet {
	set := gset.NewStrSet()
	if IsBlank(prefix) || columns == nil {
		return set
	}
	switch field := columns.(type) {
	case []string:
		for _, value := range field {
			set.Add(fmt.Sprintf("%s.%s", prefix, value))
		}
	case *[]string:
		for _, value := range *field {
			set.Add(fmt.Sprintf("%s.%s", prefix, value))
		}
	case gset.StrSet:
		field.Iterator(func(v string) bool {
			set.Add(fmt.Sprintf("%s.%s", prefix, v))
			return true
		})
	case *gset.StrSet:
		field.Iterator(func(v string) bool {
			set.Add(fmt.Sprintf("%s.%s", prefix, v))
			return true
		})
	}
	return set
}
