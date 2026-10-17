// Package middleware 提供全局中间件: 统一响应 / 跨域 / JWT 鉴权。
package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域中间件 (开发环境前后端分离部署时放行跨域请求)。
// 通过 config.yaml 的 cors.enabled 控制开关, 生产环境前后端同源部署时可关闭。
func CORS(r *ghttp.Request) {
	r.Response.CORS(ghttp.CORSOptions{
		AllowOrigin:  "*", // 放行所有来源, 也可在配置中改为指定域名
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		// accessToken 为前端携带令牌的请求头 (与 setting.config.js 的 tokenName 一致)
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With,accessToken",
		MaxAge:       86400,
	})
	if r.Method == "OPTIONS" {
		r.Response.WriteHeader(204)
		r.ExitAll()
		return
	}
	r.Middleware.Next()
}
