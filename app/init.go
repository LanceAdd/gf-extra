package app

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	appCfg *ApplicationConfig
)

func InitAppInfo() {
	ctx := gctx.GetInitCtx()
	appCfg = &ApplicationConfig{}
	name, err := g.Cfg().Get(ctx, "app.name")
	if err != nil {
		panic(err)
	}
	if name.IsEmpty() {
		panic("app.name is empty")
	}
	appCfg.AppName = name.String()
}
