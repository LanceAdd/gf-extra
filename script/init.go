package script

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	scriptsMap     *sync.Map
	scriptsSha1Map *sync.Map
)

func LoadLuaScripts(m map[string]string) {
	scriptsMap = &sync.Map{}
	scriptsSha1Map = &sync.Map{}
	ctx := gctx.GetInitCtx()
	total := 0
	for k, v := range m {
		if strings.TrimSpace(k) != "" && strings.TrimSpace(v) != "" {
			scriptsMap.Store(k, v)
			sign, err := g.Redis().ScriptLoad(ctx, v)
			if err != nil {
				panic(gerror.Wrap(err, "加载lua脚本失败: "+k))
			}
			scriptsSha1Map.Store(k, sign)
			total++
		}
	}

	builder := strings.Builder{}
	scriptsSha1Map.Range(func(key, value any) bool {
		builder.WriteString(fmt.Sprintf("[SUCCESS] Load lua script: %s => %s\n", key, value))
		return true
	})
	g.Log().Infof(ctx, "[SUCCESS] Load [%d] lua script\n%s", total, builder.String())
}

func GetLuaSha1(key string) string {
	if value, ok := scriptsSha1Map.Load(key); ok {
		return value.(string)
	}
	return ""
}
