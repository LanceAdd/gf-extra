package auth

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"strings"
)

func RequiredToken(r *ghttp.Request) {
	err := doTokenRequired(r)
	platform := defaultPlatform
	if getPlatform != nil {
		platform = getPlatform(r.GetCtx())
	}
	if err != nil {
		code := gerror.Code(err)
		r.Response.WriteJson(DefaultHandlerResponse{
			Code:     code.Code(),
			Message:  code.Message(),
			Data:     nil,
			Platform: platform,
		})
		return
	}
	r.Middleware.Next()
}

func doTokenRequired(r *ghttp.Request) error {
	token := r.GetHeader("token")
	if strings.TrimSpace(token) == "" {
		return gerror.NewCode(IllegalTokensError)
	}
	if getPlatform == nil {
		return gerror.NewCode(PlatformError)
	}
	platform := getPlatform(r.GetCtx())
	if !strings.HasPrefix(r.Router.Uri, "/"+platform) {
		return gerror.NewCode(UrlPrefixError)
	}
	cacheMode := GetCacheMode()
	currentUserId := doVerifyToken(r.GetCtx(), token, cacheMode)
	if currentUserId == 0 {
		return gerror.NewCode(IllegalTokensError)
	}
	r.SetCtxVar(CtxUserId, currentUserId)
	return nil
}

func doVerifyToken(ctx context.Context, token string, mode string) int64 {
	if strings.TrimSpace(token) == "" {
		return 0
	}
	var (
		userId int64
		err    error
	)
	switch mode {
	case CacheModeRedis:
		userId, err = doExistsTokenFromRedis(ctx, token)
	case CacheModeMemory:
		userId, err = doExistsTokenFromMemory(ctx, token)
	default:
		g.Log().Errorf(ctx, "illegal cache mode: %s", mode)
		return 0
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return 0
	}
	if userId > 0 {
		return userId
	}
	return 0
}

func doExistsTokenFromRedis(ctx context.Context, token string) (int64, error) {
	key := CachePrefixUserToken + ":" + token
	value, err := redisOps().Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if value.IsNil() {
		return 0, nil
	}
	return value.Int64(), nil
}

func doExistsTokenFromMemory(ctx context.Context, token string) (int64, error) {
	key := CachePrefixUserToken + ":" + token
	value, err := gcache.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if value.IsNil() {
		return 0, nil
	}
	return value.Int64(), nil
}
