package watcher

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	snapShot *SnapShot
)

func init() {
	snapShot = &SnapShot{}
	ctx := gctx.GetInitCtx()
	data, err := g.Cfg().Data(ctx)
	if err != nil {
		panic(err)
	}
	dataMap := doToPropertiesMap(&data)
	snapShot.CfgMap = dataMap
}

func ReInitCfgWatcher() {
	snapShot.Lock()
	defer snapShot.Unlock()
	ctx := gctx.GetInitCtx()
	data, err := g.Cfg().Data(ctx)
	if err != nil {
		panic(err)
	}
	dataMap := doToPropertiesMap(&data)
	result := Compare(dataMap, snapShot.CfgMap)
	json, err := gjson.DecodeToJson(result)
	snapShot.CfgMap = dataMap
	g.Log().Infof(ctx, "[Config Change] %s", json.String())
}
