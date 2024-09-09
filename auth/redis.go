package auth

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
)

func redisOps() *gredis.Redis {
	name := GetRedisName()
	if gutil.IsEmpty(name) {
		return g.Redis()
	}
	return g.Redis(name)
}
