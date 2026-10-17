// Package middleware 提供全局中间件: 统一响应 / 跨域 / JWT 鉴权。
package middleware

import (
	"context"

	gfjwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"

	"gf-vue-admin-better/internal/consts"
)

// TokenHeader 前端携带令牌的请求头名称 (与 setting.config.js 的 tokenName 一致)
const TokenHeader = "accessToken"

var authInstance *gfjwt.GfJWTMiddleware

// Auth 返回 JWT 鉴权实例 (登录生成令牌 / 登出拉黑令牌均通过此实例)。
func Auth() *gfjwt.GfJWTMiddleware {
	return authInstance
}

// InitAuth 根据配置文件初始化 JWT 实例 (config.yaml 的 jwt 节点), 服务启动时调用。
func InitAuth() {
	var (
		ctx     = gctx.New()
		key     = g.Cfg().MustGet(ctx, "jwt.key", "gf-vue-admin-better-secret-key").String()
		timeout = g.Cfg().MustGet(ctx, "jwt.timeout", "24h").Duration()
		maxRefr = g.Cfg().MustGet(ctx, "jwt.maxRefresh", "24h").Duration()
	)

	authInstance = gfjwt.New(&gfjwt.GfJWTMiddleware{
		Realm:       "gf-vue-admin-better",
		Key:         []byte(key),
		Timeout:     timeout,
		MaxRefresh:  maxRefr,
		TokenLookup: "header: " + TokenHeader,
		// 注意: 载荷键名不能用 "id"。
		// gf-jwt 中间件会执行 r.SetParam(IdentityKey, identity) 将当前用户ID写入路由参数,
		// 而 GF 参数解析优先级为 路由参数 > 请求体, 若键名为 "id" 会污染所有请求结构体的 id 字段
		// (如修改/重置密码接口会错传为登录用户自身ID)。使用 "uid" 避免冲突。
		IdentityKey: "uid",
		// TokenGenerator 传入的数据映射为 JWT 载荷
		PayloadFunc: func(data interface{}) gfjwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				claims := gfjwt.MapClaims{}
				for key, value := range v {
					claims[key] = value
				}
				return claims
			}
			return gfjwt.MapClaims{}
		},
		// 鉴权失败统一输出 {code:401/403, msg:...}
		Unauthorized: unauthorized,
	})
}

// unauthorized JWT 校验失败响应: 按统一结构输出 401(未登录) / 403(权限不足)
func unauthorized(ctx context.Context, code int, message string) {
	r := g.RequestFromCtx(ctx)
	if code == 403 {
		WriteJSON(r, consts.CodeNoPermission, "权限不足", nil)
		return
	}
	WriteJSON(r, consts.CodeNoLogin, "登录状态已过期, 请重新登录", nil)
}

// AdaptTokenHeader 适配前端令牌头格式。
// 前端按 tokenName 在请求头携带原始令牌(无 Bearer 前缀),
// gf-jwt 解析要求 "Bearer <token>" 格式, 此处统一补齐, 避免修改前端核心源码。
func AdaptTokenHeader(r *ghttp.Request) {
	if v := r.Header.Get(TokenHeader); v != "" && !gstr.HasPrefix(v, "Bearer ") {
		r.Header.Set(TokenHeader, "Bearer "+v)
	}
}

// MiddlewareAuth JWT 鉴权中间件: 校验令牌有效性并检查登出黑名单, 通过后放行。
func MiddlewareAuth(r *ghttp.Request) {
	AdaptTokenHeader(r)
	authInstance.MiddlewareFunc()(r)
	r.Middleware.Next()
}
