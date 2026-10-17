// Package auth 认证控制器: 登录 / 登出。
package auth

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gf-vue-admin-better/api/auth/v1"
	"gf-vue-admin-better/internal/middleware"
	"gf-vue-admin-better/internal/service"
)

// Controller 认证控制器
type Controller struct{}

// New 创建认证控制器实例
func New() *Controller {
	return &Controller{}
}

// Login 账号密码登录
func (c *Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	return service.Auth().Login(ctx, req)
}

// Logout 退出登录 (令牌加入黑名单)
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	// 前端请求头携带原始令牌, 适配 gf-jwt 解析格式后再执行登出
	middleware.AdaptTokenHeader(g.RequestFromCtx(ctx))
	return service.Auth().Logout(ctx, req)
}

// 确保 Controller 实现接口 (编译期校验)
var _ service.IAuth = (*Controller)(nil)
