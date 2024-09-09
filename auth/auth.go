package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gutil"
)

const defaultPlatform = "Default"

func RequiredAuth(r *ghttp.Request) {
	err := doAuthRequired(r)
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

func doAuthRequired(r *ghttp.Request) error {
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
	path := fmt.Sprintf("%s:%s", r.Method, r.Router.Uri)
	permissions := GetRule(path)
	cacheMode := GetCacheMode()
	currentUser := doVerifyToken(r.GetCtx(), token, cacheMode)
	if currentUser == nil {
		return gerror.NewCode(IllegalTokensError)
	}
	r.SetCtxVar(CtxUserId, currentUser.UserId)
	if len(permissions) == 0 {
		return nil
	}
	pass := doVerifyAuth(r.GetCtx(), currentUser.UserId, permissions, cacheMode)
	if pass {
		return nil
	}
	return gerror.NewCode(IllegalPermissionError)
}

func doVerifyAuth(ctx context.Context, userId int64, permissions []RulePermission, mode string) bool {
	if getGroupCodes == nil {
		g.Log().Error(ctx, "func getGroupCodes do not init")
		return false
	}
	if getPermissionCodes == nil {
		g.Log().Error(ctx, "func getPermissionCodes do not init")
		return false
	}
	switch mode {
	case CacheModeRedis:
		return doAuthFromRedis(ctx, userId, permissions)
	case CacheModeMemory:
		return doAuthFromMemory(ctx, userId, permissions)
	case CacheModeNone:
		return doAuthFromDb(ctx, userId, permissions)
	}
	return false
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

func doAuth(permissions []RulePermission, groupCodes *gset.StrSet, permissionCodes *gset.StrSet) bool {
	for _, permission := range permissions {
		if permission.Group != nil && len(permission.Group) > 0 {
			if strings.TrimSpace(permission.GroupMode) == "" || strings.EqualFold(permission.GroupMode, OR) {
				set := groupCodes.Intersect(doUpperStrSet(permission.Group))
				if set.Size() > 0 {
					return true
				}
			}
			if strings.EqualFold(permission.GroupMode, AND) {
				set := groupCodes.Intersect(doUpperStrSet(permission.Group))
				if set.Size() == groupCodes.Size() {
					return true
				}
			}
			if strings.EqualFold(permission.GroupMode, EXCLUDE) {
				set := groupCodes.Intersect(doUpperStrSet(permission.Group))
				if set.Size() == 0 {
					return true
				}
			}
		}
		if permission.Code != nil && len(permission.Code) > 0 {
			if strings.TrimSpace(permission.CodeMode) == "" || strings.EqualFold(permission.CodeMode, OR) {
				set := permissionCodes.Intersect(doUpperStrSet(permission.Code))
				if set.Size() > 0 {
					return true
				}
			}
			if strings.EqualFold(permission.CodeMode, AND) {
				set := permissionCodes.Intersect(doUpperStrSet(permission.Code))
				if set.Size() == permissionCodes.Size() {
					return true
				}
			}
			if strings.EqualFold(permission.CodeMode, EXCLUDE) {
				set := permissionCodes.Intersect(doUpperStrSet(permission.Code))
				if set.Size() == 0 {
					return true
				}
			}
		}
	}
	return false
}

func doAuthFromRedis(ctx context.Context, userId int64, permissions []RulePermission) bool {
	cacheKey := fmt.Sprintf("%s:%d", CachePrefixUserPermission, userId)
	data, err := redisOps().HMGet(ctx, cacheKey, "GroupCodes", "PermissionCodes")
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if data[0].IsEmpty() || data[1].IsEmpty() {
		groupCodes, err := getGroupCodes(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		permissionCodes, err := getPermissionCodes(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		err = redisOps().HMSet(ctx, cacheKey, map[string]interface{}{"GroupCodes": groupCodes, "PermissionCodes": permissionCodes})
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		expire, err := redisOps().Expire(ctx, cacheKey, GetCacheExpireDt())
		if err != nil {
			g.Log().Errorf(ctx, "expireResult: %d, err: %v", expire, err)
			return false
		}
		return doAuth(permissions, gset.NewStrSetFrom(groupCodes), gset.NewStrSetFrom(permissionCodes))
	}
	groupCodes := gset.NewStrSetFrom(data[0].Strings())
	permissionCodes := gset.NewStrSetFrom(data[1].Strings())
	return doAuth(permissions, groupCodes, permissionCodes)
}

func doAuthFromMemory(ctx context.Context, userId int64, permissions []RulePermission) bool {
	cacheKey := fmt.Sprintf("%s:%d", CachePrefixUserPermission, userId)
	data, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if data.IsNil() {
		groupCodes, err := getGroupCodes(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		permissionCodes, err := getPermissionCodes(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		err = gcache.Set(ctx, cacheKey, map[string][]string{"GroupCodes": groupCodes, "PermissionCodes": permissionCodes}, GetCacheExpireDtDuration())
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		return doAuth(permissions, gset.NewStrSetFrom(groupCodes), gset.NewStrSetFrom(permissionCodes))

	}
	varMap := data.MapStrVar()
	groupCodes := varMap["GroupCodes"].Strings()
	permissionCodes := varMap["PermissionCodes"].Strings()
	return doAuth(permissions, gset.NewStrSetFrom(groupCodes), gset.NewStrSetFrom(permissionCodes))
}

func doAuthFromDb(ctx context.Context, userId int64, permissions []RulePermission) bool {
	groupCodes, err := getGroupCodes(ctx, userId)
	if err != nil {
		return false
	}
	permissionCodes, err := getPermissionCodes(ctx, userId)
	if err != nil {
		return false
	}
	return doAuth(permissions, gset.NewStrSetFrom(groupCodes), gset.NewStrSetFrom(permissionCodes))
}

func doUpperStrSet(source []string) *gset.StrSet {
	set := gset.NewStrSet()
	for _, item := range source {
		set.Add(strings.ToUpper(item))
	}
	return set
}
