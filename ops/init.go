package ops

import (
	"sync"
)

var (
	cache       *sync.Map // 用于存储结构体或者指针指向的结构体的orm数据库映射map
	fieldsCache *sync.Map // 用于存储查询结果的结构体的Fields列表避免重复构建
)

func init() {
	cache = &sync.Map{}
	fieldsCache = &sync.Map{}
}
