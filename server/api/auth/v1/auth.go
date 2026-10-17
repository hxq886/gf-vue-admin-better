// Package v1 定义认证相关接口契约 (登录 / 登出)。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LoginReq 登录请求
type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"认证管理" summary:"账号密码登录"`
	Username string `json:"username" v:"required#请输入用户名" dc:"登录账号"`
	Password string `json:"password" v:"required#请输入密码" dc:"登录密码"`
}

// LoginRes 登录响应 (返回 JWT 令牌, 字段名与前端 tokenName 配置一致)
type LoginRes struct {
	AccessToken string `json:"accessToken" dc:"访问令牌"`
}

// LogoutReq 登出请求 (携带令牌, 服务端将令牌加入黑名单)
type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" tags:"认证管理" summary:"退出登录"`
}

// LogoutRes 登出响应 (无业务数据)
type LogoutRes struct{}
