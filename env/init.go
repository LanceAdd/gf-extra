package env

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
)

const defaultEnv = "default"

func InitEnv() {
	var Env string
	ctx := gctx.GetInitCtx()
	env := genv.GetWithCmd("ENV")
	if env.IsNil() {
		Env = defaultEnv
	} else {
		Env = env.String()
	}
	if Env == defaultEnv {
		g.Log().Infof(ctx, "[Success] App start with ENV=%s, config=%s", Env, "config.yaml")
		return
	}
	config := fmt.Sprintf("config_%s.yaml", Env)
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(config)
	g.Log().Infof(ctx, "[Success] App start with ENV=%s, config=%s", Env, config)
}
