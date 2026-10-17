// Package consts 定义全局常量。
package consts

// 接口统一响应码约定 (与前端 net.config.js 严格对齐)
const (
	CodeSuccess        = 200 // 成功
	CodeNoLogin        = 401 // 未登录 / Token 失效
	CodeNoPermission   = 403 // 权限不足
	CodeInternalError  = 500 // 服务异常 / 业务错误
)

// 菜单类型
const (
	MenuTypeDirectory = "M" // 目录
	MenuTypeMenu      = "C" // 菜单
	MenuTypeButton    = "F" // 按钮
)

// 内置角色与账号约定
const (
	RoleKeyAdmin  = "admin"    // 超级管理员角色权限字符
	UserIdAdmin   = uint64(1)  // 内置管理员用户ID
	RoleSuper     = uint64(1)  // 内置超级管理员角色ID
)
