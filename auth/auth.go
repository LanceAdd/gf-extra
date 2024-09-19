package auth

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gutil"
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
	currentUser := doVerifyToken(r.GetCtx(), token, cacheMode)
	if currentUser == nil {
		return gerror.NewCode(IllegalTokensError)
	}
	r.SetCtxVar(CtxUserId, currentUser.UserId)
	return nil
}

func doVerifyToken(ctx context.Context, token string, mode string) *CurrentUser {
	if strings.TrimSpace(token) == "" {
		return nil
	}
	var (
		user *CurrentUser
		err  error
	)
	switch mode {
	case CacheModeRedis:
		user, err = doExistsTokenFromRedis(ctx, token)
	case CacheModeMemory:
		user, err = doExistsTokenFromMemory(ctx, token)
	case CacheModeNone:
		user, err = doVerifyTokenByUnSignVerify(ctx, token)
	default:
		g.Log().Errorf(ctx, "illegal cache mode: %s", mode)
		return nil
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return nil
	}
	if user != nil {
		return user
	}
	return nil
}

func doExistsTokenFromRedis(ctx context.Context, token string) (*CurrentUser, error) {
	key := CachePrefixUserToken + ":" + token
	tmp, err := redisOps().Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if tmp.IsNil() {
		return nil, nil
	}
	content := &SimpleTokenContent{}
	err = tmp.Scan(content)
	if err != nil {
		return nil, err
	}
	if gutil.IsEmpty(content.UserId) {
		return nil, nil
	}
	if !content.IsValidate() {
		return nil, nil
	}
	return &CurrentUser{UserId: content.UserId}, nil
}

func doExistsTokenFromMemory(ctx context.Context, token string) (*CurrentUser, error) {
	key := CachePrefixUserToken + ":" + token
	tmp, err := gcache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	content := &SimpleTokenContent{}
	err = tmp.Scan(content)
	if err != nil {
		return nil, err
	}
	if gutil.IsEmpty(content.UserId) {
		return nil, nil
	}
	if !content.IsValidate() {
		return nil, nil
	}
	return &CurrentUser{UserId: content.UserId}, nil
}
func doVerifyTokenByUnSignVerify(ctx context.Context, token string) (*CurrentUser, error) {
	if unSignToken == nil {
		g.Log().Error(ctx, "func unSignToken do not init")
		return nil, errors.New("func signToken do not init")
	}
	content, err := unSignToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if gutil.IsEmpty(content.UserId) {
		return nil, nil
	}
	if !content.IsValidate() {
		return nil, nil
	}
	return &CurrentUser{UserId: content.UserId}, nil
}
