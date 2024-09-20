package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"strings"
)

func RequiredPermission(r *ghttp.Request) {
	err := doPermissionRequired(r)
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

func doPermissionRequired(r *ghttp.Request) error {
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
	value := r.GetCtxVar(CtxUserId)
	var userId int64
	if value.IsNil() {
		currentUserId := doVerifyToken(r.GetCtx(), token, cacheMode)
		if currentUserId == 0 {
			return gerror.NewCode(IllegalTokensError)
		}
		r.SetCtxVar(CtxUserId, currentUserId)
	} else {
		userId = value.Int64()
	}
	path := fmt.Sprintf("%s:%s", r.Method, r.Router.Uri)
	permissions := GetRule(path)
	if len(permissions) == 0 {
		return nil
	}
	pass := doVerifyAuth(r.GetCtx(), userId, permissions, cacheMode)
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
	}
	return false
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
