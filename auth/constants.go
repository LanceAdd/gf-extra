package auth

const (
	CachePrefixUserToken      = "USER_TOKEN"
	CachePrefixUserPermission = "USER_PERMISSION"
	CachePrefixUserTokenArray = "USER_TOKEN_ARRAY"
)

const (
	CacheModeRedis  = "REDIS"
	CacheModeMemory = "MEMORY"
	CacheModeNone   = "NONE"
)

const (
	OR      = "OR"
	AND     = "AND"
	EXCLUDE = "EXCLUDE"
)

const CtxUserId = "CtxUserId"

const defaultPlatform = "Default"
