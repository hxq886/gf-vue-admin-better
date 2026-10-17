// Package service 定义业务层接口, 由 internal/logic 提供实现并注册。
package service

import (
	"context"

	"gf-vue-admin-better/api/system/v1"
)

// IUser 用户服务接口。
type IUser interface {
	// GetInfo 获取当前登录用户信息 (昵称/角色/按钮权限标识)
	GetInfo(ctx context.Context) (res *v1.GetInfoRes, err error)
	// List 分页查询用户列表
	List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)
	// Create 新增用户
	Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error)
	// Update 修改用户基本信息
	Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error)
	// Password 重置用户密码
	Password(ctx context.Context, req *v1.UserPasswordReq) (res *v1.UserPasswordRes, err error)
	// Delete 批量删除用户
	Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error)
}

var localIUser IUser

// User 获取用户服务实例
func User() IUser {
	if localIUser == nil {
		panic("service IUser not registered")
	}
	return localIUser
}

// RegisterUser 注册用户服务实现
func RegisterUser(i IUser) {
	localIUser = i
}
