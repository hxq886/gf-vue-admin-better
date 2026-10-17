// Package service 定义业务层接口, 由 internal/logic 提供实现并注册。
// 遵循 GoFrame 分层规范: controller -> service -> logic -> dao。
package service

import (
	"context"

	"gf-vue-admin-better/api/auth/v1"
)

// IAuth 认证服务接口。
type IAuth interface {
	// Login 账号密码登录, 返回 JWT 令牌
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	// Logout 登出, 将令牌加入黑名单
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
}

var localIAuth IAuth

// Auth 获取认证服务实例 (未注册时 panic, 由启动流程保证注册)
func Auth() IAuth {
	if localIAuth == nil {
		panic("service IAuth not registered")
	}
	return localIAuth
}

// RegisterAuth 注册认证服务实现
func RegisterAuth(i IAuth) {
	localIAuth = i
}
