package ops

import (
	"github.com/gogf/gf/v2/container/gset"
)

type GroupOpsImpl struct {
	items *gset.StrSet
}

func GroupOps() *GroupOpsImpl {
	return &GroupOpsImpl{items: gset.NewStrSet(false)}
}

func (g *GroupOpsImpl) IncludePrefixColumns(prefix string, columns any) *GroupOpsImpl {
	slice := doHandleGroupOpsPrefixColumnsToSet(prefix, columns)
	g.items.Merge(slice)
	return g
}

func (g *GroupOpsImpl) IncludePrefixObject(prefix string, columns any) *GroupOpsImpl {
	slice := doHandleGroupOpsObjectToSet(prefix, columns)
	g.items.Merge(slice)
	return g
}

func (g *GroupOpsImpl) IncludePrefixSlice(prefixOrKey string, columns any) *GroupOpsImpl {
	slice := doHandleGroupOpsPrefixSliceToSet(prefixOrKey, columns)
	g.items.Merge(slice)
	return g
}

func (g *GroupOpsImpl) Exclude(keys ...string) *GroupOpsImpl {
	if len(keys) == 0 {
		return g
	}
	for _, key := range keys {
		g.items.Remove(key)
	}
	return g
}

func (g *GroupOpsImpl) GetGroupFields() []string {
	return g.items.Slice()
}
