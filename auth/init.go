package auth

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	ruleCfg            *RuleConfig
	ruleMap            *sync.Map
	getGroupCodes      func(ctx context.Context, userId int64) ([]string, error)
	getPermissionCodes func(ctx context.Context, userId int64) ([]string, error)
	signToken          func(ctx context.Context, s *SimpleTokenContent) (string, error)
	unSignToken        func(ctx context.Context, content string) (*SimpleTokenContent, error)
	getPlatform        func(ctx context.Context) string
)

func InitAuthRule() {
	ctx := gctx.GetInitCtx()
	ruleCfg = &RuleConfig{}
	ruleMap = &sync.Map{}
	doLoadAuthConfig(ctx, ruleCfg)
	doInitAuthRule(ruleCfg, ruleMap)
	g.Log().Infof(ctx, "[SUCCESS] Load %d auth rule", len(ruleCfg.Rule))
}

func SetGetGroupCodesFunc(f func(ctx context.Context, userId int64) ([]string, error)) {
	getGroupCodes = f
}

func SetGetPermissionCodes(f func(ctx context.Context, userId int64) ([]string, error)) {
	getPermissionCodes = f
}

func SetSignTokenFunc(f func(ctx context.Context, s *SimpleTokenContent) (string, error)) {
	signToken = f
}

func SetUnSignTokenFunc(f func(ctx context.Context, content string) (*SimpleTokenContent, error)) {
	unSignToken = f
}

func SetGetPlatformFunc(f func(ctx context.Context) string) {
	getPlatform = f
}

func GetRule(path string) []RulePermission {
	value, ok := ruleMap.Load(path)
	if ok {
		permissions := value.([]RulePermission)
		slice := make([]RulePermission, len(permissions))
		copy(slice, permissions)
		return slice
	}
	return nil
}

func GetCacheMode() string {
	ruleCfg.RLock()
	defer ruleCfg.RUnlock()
	return strings.ToUpper(ruleCfg.CacheMode)
}

func GetCacheExpireDt() int64 {
	ruleCfg.RLock()
	defer ruleCfg.RUnlock()
	return ruleCfg.ExpireDt
}

func GetRedisName() string {
	ruleCfg.RLock()
	defer ruleCfg.RUnlock()
	return ruleCfg.RedisName
}

func GetCacheExpireDtDuration() time.Duration {
	ruleCfg.RLock()
	defer ruleCfg.RUnlock()
	return time.Duration(ruleCfg.ExpireDt) * time.Second
}

func GetIssuer() string {
	ruleCfg.RLock()
	defer ruleCfg.RUnlock()
	return ruleCfg.Issuer
}

func ReInitAuthRule() {
	ruleCfg.Lock()
	defer ruleCfg.Unlock()
	ctx := gctx.GetInitCtx()
	newAuthCfg := &RuleConfig{}
	doLoadAuthConfig(ctx, newAuthCfg)
	equal := reflect.DeepEqual(newAuthCfg, ruleCfg)
	if !equal {
		ruleCfg = newAuthCfg
		newAuthRule := &sync.Map{}
		doInitAuthRule(ruleCfg, newAuthRule)
		ruleMap = newAuthRule
		g.Log().Infof(ctx, "[SUCCESS] ReLoad %d auth rule", len(ruleCfg.Rule))
	}
}

func doLoadAuthConfig(ctx context.Context, authCfg *RuleConfig) {
	config, err := g.Cfg().Get(ctx, "auth")
	if err != nil {
		panic(err)
	}
	err = gconv.Scan(config, authCfg)
	if err != nil {
		panic(err)
	}
	err = g.Validator().Data(authCfg).Run(ctx)
	if err != nil {
		panic(err)
	}
}

func doInitAuthRule(config *RuleConfig, authRule *sync.Map) {
	for index := range config.Rule {
		authRule.Store(config.Rule[index].Path, config.Rule[index].Permission)
	}
}
