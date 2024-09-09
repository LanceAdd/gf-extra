package auth

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

var (
	IllegalPermissionError = apiError(40000, "非法权限")
	IllegalTokensError     = apiError(40001, "非法token")
	UrlPrefixError         = apiError(40002, "路由前缀错误")
	PlatformError          = apiError(40003, "平台名称错误")
)

func apiError(errorCode int, errorMessage string) gcode.Code {
	return gcode.New(errorCode, errorMessage, nil)
}
