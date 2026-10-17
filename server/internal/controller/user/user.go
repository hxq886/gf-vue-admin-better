// Package user 用户控制器: 用户信息 / 用户CRUD。
package user

import (
	"context"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/service"
)

// Controller 用户控制器
type Controller struct{}

// New 创建用户控制器实例
func New() *Controller {
	return &Controller{}
}

// GetInfo 获取当前登录用户信息
func (c *Controller) GetInfo(ctx context.Context, req *v1.GetInfoReq) (res *v1.GetInfoRes, err error) {
	return service.User().GetInfo(ctx)
}

// List 用户列表
func (c *Controller) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	return service.User().List(ctx, req)
}

// Create 新增用户
func (c *Controller) Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	return service.User().Create(ctx, req)
}

// Update 修改用户
func (c *Controller) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	return service.User().Update(ctx, req)
}

// Password 重置密码
func (c *Controller) Password(ctx context.Context, req *v1.UserPasswordReq) (res *v1.UserPasswordRes, err error) {
	return service.User().Password(ctx, req)
}

// Delete 删除用户
func (c *Controller) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	return service.User().Delete(ctx, req)
}
