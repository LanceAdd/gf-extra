package ops

import (
	"fmt"
)

type FieldOpsImpl struct {
	itemMap map[string]string
	items   []string
}

func FieldOps() *FieldOpsImpl {
	return &FieldOpsImpl{itemMap: make(map[string]string)}
}
func (f *FieldOpsImpl) Include(key string, alias string) *FieldOpsImpl {
	if IsBlank(key) || IsBlank(alias) {
		return f
	}
	f.itemMap[key] = alias
	return f
}

func (f *FieldOpsImpl) IncludeRawSql(rawSql string) *FieldOpsImpl {
	if IsBlank(rawSql) {
		return f
	}
	f.items = append(f.items, rawSql)
	return f
}

func (f *FieldOpsImpl) IncludePrefixColumns(prefix string, columns any) *FieldOpsImpl {
	m := doHandleFieldOpsPrefixColumnsToMap(prefix, columns)
	doMergeMap(&f.itemMap, m)
	return f
}

func (f *FieldOpsImpl) IncludePrefixObject(prefix string, columns any) *FieldOpsImpl {
	m := doHandleFieldOpsObjectToMap(prefix, columns)
	doMergeMap(&f.itemMap, m)
	return f
}

func (f *FieldOpsImpl) IncludePrefixSlice(prefixOrKey string, columns any) *FieldOpsImpl {
	m := doHandleFieldOpsPrefixSliceToMap(prefixOrKey, columns)
	doMergeMap(&f.itemMap, m)
	return f
}

func (f *FieldOpsImpl) IncludeMap(columns any) *FieldOpsImpl {
	m := doHandleFieldOpsMapToMap(columns)
	doMergeMap(&f.itemMap, m)
	return f
}

func (f *FieldOpsImpl) IncludePrefixMap(prefix string, columns any) *FieldOpsImpl {
	m := doHandleFieldOpsPrefixMapToMap(prefix, columns)
	doMergeMap(&f.itemMap, m)
	return f
}

func (f *FieldOpsImpl) Omitempty() *FieldOpsImpl {
	delete(f.itemMap, "")
	var keys []string
	for key, value := range f.itemMap {
		if IsBlank(value) {
			keys = append(keys, key)
		}
		if IsBlank(key) {
			keys = append(keys, value)
		}
	}
	for _, key := range keys {
		delete(f.itemMap, key)
	}

	var items []string
	for _, item := range f.items {
		if !IsBlank(item) {
			items = append(items, item)
		}
	}
	f.items = items
	return f
}

func (f *FieldOpsImpl) ExcludeByColumns(keys ...string) *FieldOpsImpl {
	if len(keys) == 0 {
		return f
	}
	for _, key := range keys {
		delete(f.itemMap, key)
	}
	return f
}

func (f *FieldOpsImpl) ExcludeByAlias(keys ...string) *FieldOpsImpl {
	if len(keys) == 0 {
		return f
	}
	var deleteKeys []string
	for column, alias := range f.itemMap {
		for _, key := range keys {
			if key == alias {
				deleteKeys = append(deleteKeys, column)
			}
		}
	}
	for _, key := range deleteKeys {
		delete(f.itemMap, key)
	}
	return f
}

func (f *FieldOpsImpl) ExcludeRawSql(keys ...string) *FieldOpsImpl {
	f.items = []string{}
	return f
}

func (f *FieldOpsImpl) Items() []string {
	var items []string
	for key, value := range f.itemMap {
		items = append(items, fmt.Sprintf("%s as %s", key, value))
	}
	items = append(items, f.items...)
	return items
}

func (f *FieldOpsImpl) Fields() []string {
	return f.Omitempty().Items()
}
