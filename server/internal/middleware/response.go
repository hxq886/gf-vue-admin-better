// Package middleware 提供全局中间件: 统一响应 / 跨域 / JWT 鉴权。
package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gf-vue-admin-better/internal/consts"
)

// Response 统一响应中间件。
// 将 handler 返回值/错误统一包装为 {code:number, msg:string, data:any} 结构,
// 与前端 net.config.js 约定对齐: 200成功 / 401未登录 / 403无权限 / 500服务异常。
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	// 已有自定义输出(如 JWT 未授权响应), 直接返回
	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
	)

	if err == nil {
		WriteJSON(r, consts.CodeSuccess, "success", res)
		return
	}

	// 保留约定的业务码(401/403), 其余错误统一归为 500
	code := gerror.Code(err).Code()
	if code != consts.CodeNoLogin && code != consts.CodeNoPermission {
		code = consts.CodeInternalError
	}
	WriteJSON(r, code, err.Error(), nil)
}

// WriteJSON 按统一格式输出 JSON 响应 (HTTP 状态码恒为 200, 业务码由 body.code 承载)
func WriteJSON(r *ghttp.Request, code int, msg string, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	r.Response.WriteJson(g.Map{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// CodeToGcode 将约定业务码转换为 gcode, 供业务层 gerror.NewCode 使用
func CodeToGcode(code int) gcode.Code {
	return gcode.New(code, "", nil)
}
