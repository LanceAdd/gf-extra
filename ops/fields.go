package ops

import (
	"fmt"
	"reflect"
)

type FieldsOpsImpl[T FieldsInterface] struct {
	item *T
}

func FieldsOps[T FieldsInterface]() *FieldsOpsImpl[T] {
	return &FieldsOpsImpl[T]{item: new(T)}
}

func (f *FieldsOpsImpl[T]) GetQueryFields(cached ...bool) []string {
	fromCache := doBoolDefaultTrue(cached...)
	if fromCache {
		of := reflect.TypeOf(f.item)
		if of.Kind() == reflect.Ptr {
			if of.Elem().Kind() == reflect.Struct {
				cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyQueryFields, of.Elem().PkgPath(), of.Elem().Name())
				if value, ok := fieldsCache.Load(cacheKey); ok {
					return doCopySlice(value.([]string))
				} else {
					fields := (*f.item).GetQueryFields()
					doSortFields(fields)
					fieldsCache.Store(cacheKey, doCopySlice(fields))
					return fields
				}
			}
		}
	}
	return (*f.item).GetQueryFields()
}

func (f *FieldsOpsImpl[T]) GetGroupFields(cached ...bool) []string {
	fromCache := doBoolDefaultTrue(cached...)
	if fromCache {
		of := reflect.TypeOf(f.item)
		if of.Kind() == reflect.Ptr {
			if of.Elem().Kind() == reflect.Struct {
				cacheKey := fmt.Sprintf("%s:%s/%s", CacheKeyGroupFields, of.Elem().PkgPath(), of.Elem().Name())
				if value, ok := fieldsCache.Load(cacheKey); ok {
					return doCopySlice(value.([]string))
				} else {
					fields := (*f.item).GetGroupFields()
					doSortFields(fields)
					fieldsCache.Store(cacheKey, doCopySlice(fields))
					return fields
				}
			}
		}
	}
	return (*f.item).GetGroupFields()
}
