package encoding

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	jwtKeys *JwtKeys
	jwtCfg  *JwtConfig
)

var (
	rsaKeys *RsaKeys
	rsaCfg  *RsaConfig
)

func InitJwtEncodingKeys() {
	jwtCfg = &JwtConfig{}
	jwtKeys = &JwtKeys{}
	ctx := gctx.GetInitCtx()
	doLoadJwtConfig(ctx, jwtCfg)
	doLoadJwt(jwtCfg, jwtKeys)
	g.Log().Info(ctx, "[SUCCESS] Load jwt config")
}

func InitRsaEncodingKeys() {
	rsaKeys = &RsaKeys{}
	rsaCfg = &RsaConfig{}
	ctx := gctx.GetInitCtx()
	doLoadRsaConfig(ctx, rsaCfg)
	doLoadRsaKeys(rsaCfg, rsaKeys)
	g.Log().Info(ctx, "[SUCCESS] Load rsa config")
}
func ReInitJwtEncodingKeys() {
	jwtCfg.Lock()
	defer jwtCfg.Unlock()
	jwtKeys.Lock()
	defer jwtKeys.Unlock()
	newJwtCfg := &JwtConfig{}
	ctx := gctx.GetInitCtx()
	doLoadJwtConfig(ctx, newJwtCfg)
	jwtCfgEqual := doJwtConfigCompare(newJwtCfg, jwtCfg)
	if !jwtCfgEqual {
		jwtCfg = newJwtCfg
		newKeys := &JwtKeys{}
		doLoadJwt(newJwtCfg, newKeys)
		g.Log().Infof(ctx, "[Success] ReLoad Security Config")
	}

	rsaCfg.Lock()
	defer rsaCfg.Unlock()
	rsaKeys.Lock()
	defer rsaKeys.Unlock()
	newRsaCfg := &RsaConfig{}
	doLoadRsaConfig(ctx, newRsaCfg)
	rsaKeys = &RsaKeys{}
	doLoadRsaKeys(rsaCfg, rsaKeys)
	rsaCfgEqual := doRsaCfgCompare(newRsaCfg, rsaCfg)
	if !rsaCfgEqual {
		rsaCfg = newRsaCfg
		newRsaKeys := &RsaKeys{}
		doLoadRsaKeys(rsaCfg, newRsaKeys)
	}
}

func ReInitRsaEncodingKeys() {
	ctx := gctx.GetInitCtx()
	rsaCfg.Lock()
	defer rsaCfg.Unlock()
	rsaKeys.Lock()
	defer rsaKeys.Unlock()
	newRsaCfg := &RsaConfig{}
	doLoadRsaConfig(ctx, newRsaCfg)
	rsaKeys = &RsaKeys{}
	doLoadRsaKeys(rsaCfg, rsaKeys)
	rsaCfgEqual := doRsaCfgCompare(newRsaCfg, rsaCfg)
	if !rsaCfgEqual {
		rsaCfg = newRsaCfg
		newRsaKeys := &RsaKeys{}
		doLoadRsaKeys(rsaCfg, newRsaKeys)
	}
}
