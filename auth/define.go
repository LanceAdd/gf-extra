package auth

import (
	"sync"
	"time"
)

type DefaultHandlerResponse struct {
	Code     int         `json:"code"    dc:"Error code"`
	Message  string      `json:"message" dc:"Error message"`
	Platform string      `json:"platform" dc:"Error Platform" `
	Data     interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

type RuleConfig struct {
	sync.RWMutex
	RedisName string `json:"redisName" dc:"redis名称"`
	Issuer    string `json:"issuer" v:"required" dc:"Token签发者"`
	CacheMode string `json:"cacheMode" v:"required|in:redis,memory,none" d:"no" dc:"token是否使用缓存以及缓存的方式"`
	ExpireDt  int64  `json:"expireDt" v:"required-unless:cacheMode" dc:"缓存有效时间"`
	Rule      []Rule `json:"rule" dc:"权限路由集合"`
}

type Rule struct {
	sync.RWMutex
	Path       string           `json:"path" v:"required" dc:"路由"`
	Permission []RulePermission `json:"permission" dc:"权限内容"`
}
type RulePermission struct {
	Group     []string `json:"group" dc:"角色限制"`
	GroupMode string   `json:"groupMode" v:"required-with:group|in:or,and,exclude,OR,AND,EXCLUDE" d:"or" dc:"角色限制模式" dc:"权限以OR还是AND的方式校验"`
	Code      []string `json:"code"`
	CodeMode  string   `json:"codeMode" v:"required-with:code|in:or,and,exclude,OR,AND,EXCLUDE" d:"or" dc:"权限以OR还是AND的方式校验"`
}

type CurrentUser struct {
	UserId int64
}

// SimpleTokenContent jwt token原始内容
type SimpleTokenContent struct {
	UserId   int64     `json:"userId" v:"required" dc:"用户ID"`
	ExpireDt int64     `json:"expireDt"  v:"required" dc:"有效期"`
	IssuedDt time.Time `json:"issuedDt"  v:"required" dc:"签发时间"`
	Issuer   string    `json:"issuer"  v:"required" dc:"签发者"`
}

func (j *SimpleTokenContent) IsValidate() bool {
	deadlineDt := j.IssuedDt.Add(time.Duration(j.ExpireDt) * time.Second)
	return deadlineDt.After(time.Now())
}

// NewSimpleTokenContent 快速创建token内容
func NewSimpleTokenContent(userId int64) *SimpleTokenContent {
	return &SimpleTokenContent{
		UserId:   userId,
		ExpireDt: GetCacheExpireDt(),
		IssuedDt: time.Now(),
		Issuer:   GetIssuer(),
	}
}
