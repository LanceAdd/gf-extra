package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gutil"
)

func Login(ctx context.Context, userId int64) (string, error) {
	if gutil.IsEmpty(userId) {
		return "", errors.New("userId empty")
	}
	mode := GetCacheMode()
	switch mode {
	case CacheModeRedis:
		return GenerateTokenToRedis(ctx, userId)
	case CacheModeMemory:
		return GenerateTokenToMemory(ctx, userId)
	default:
		return "", errors.New("invalid cache mode")
	}
}

func Logout(ctx context.Context, userId int64) (string, error) {
	mode := GetCacheMode()
	switch mode {
	case CacheModeRedis:
		return GenerateTokenToRedis(ctx, userId)
	case CacheModeMemory:
		return GenerateTokenToMemory(ctx, userId)
	default:
		return "", nil
	}
}

func GenerateTokenToRedis(ctx context.Context, userId int64) (string, error) {
	token := uuid.New().String()
	cacheKey := CachePrefixUserToken + ":" + token
	err := redisOps().SetEX(ctx, cacheKey, userId, GetCacheExpireDt())
	if err != nil {
		return "", err
	}
	arrayCacheKey := fmt.Sprintf("%s:%d", CachePrefixUserTokenArray, userId)
	_, err = redisOps().SAdd(ctx, arrayCacheKey, token)
	if err != nil {
		return "", err
	}
	_, err = redisOps().Expire(ctx, arrayCacheKey, GetCacheExpireDt())
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateTokenToMemory(ctx context.Context, userId int64) (string, error) {
	token := uuid.New().String()
	cacheKey := CachePrefixUserToken + ":" + token
	err := gcache.Set(ctx, cacheKey, userId, GetCacheExpireDtDuration())
	if err != nil {
		return "", err
	}
	arrayCacheKey := fmt.Sprintf("%s:%d", CachePrefixUserTokenArray, userId)
	tokenArrayCache, err := gcache.Get(ctx, arrayCacheKey)
	if err != nil {
		return "", err
	}
	if tokenArrayCache.IsEmpty() {
		err := gcache.Set(ctx, arrayCacheKey, []string{token}, GetCacheExpireDtDuration())
		if err != nil {
			return "", err
		}
	} else {
		tokenArray := tokenArrayCache.Strings()
		tokenArray = append(tokenArray, token)
		err := gcache.Set(ctx, arrayCacheKey, tokenArray, GetCacheExpireDtDuration())
		if err != nil {
			return "", err
		}
	}
	return token, nil
}
func ClearUserTokenInRedis(ctx context.Context, userId int64) error {
	arrayCacheKey := fmt.Sprintf("%s:%d", CachePrefixUserTokenArray, userId)

	tokenArrayCache, err := redisOps().Get(ctx, arrayCacheKey)
	if err != nil {
		return err
	}
	tokenArray := tokenArrayCache.Strings()
	cacheKeys := make([]string, len(tokenArray))
	for _, token := range tokenArray {
		cacheKey := CachePrefixUserToken + ":" + token
		cacheKeys = append(cacheKeys, cacheKey)
	}
	_, err = redisOps().Del(ctx, cacheKeys...)
	if err != nil {
		return err
	}
	return nil
}

func ClearUserTokenInMemory(ctx context.Context, userId int64) error {
	arrayCacheKey := fmt.Sprintf("%s:%d", CachePrefixUserTokenArray, userId)

	tokenArrayCache, err := gcache.Get(ctx, arrayCacheKey)
	if err != nil {
		return err
	}
	tokenArray := tokenArrayCache.Strings()
	cacheKeys := g.Slice{}
	for _, token := range tokenArray {
		cacheKey := CachePrefixUserToken + ":" + token
		cacheKeys = append(cacheKeys, cacheKey)
	}
	err = gcache.Removes(ctx, cacheKeys)
	if err != nil {
		return err
	}
	return nil
}
