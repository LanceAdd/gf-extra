package app

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type RespBody struct {
	Code     int         `json:"code"    dc:"Error code"`
	Message  string      `json:"message" dc:"Error message"`
	Platform string      `json:"platform" dc:"Error Platform" `
	Data     interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

func RespWriter(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.BufferLength() > 0 {
		return
	}
	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			msg = http.StatusText(r.Response.Status)
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			err = gerror.NewCode(code, msg)
			r.SetError(err)
		} else {
			code = gcode.CodeOK
		}
	}
	r.Response.WriteJson(RespBody{
		Code:     code.Code(),
		Message:  msg,
		Data:     res,
		Platform: GetAppName(),
	})
}
