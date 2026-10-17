// Package role 角色控制器: 角色CRUD。
package role

import (
	"context"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/service"
)

// Controller 角色控制器
type Controller struct{}

// New 创建角色控制器实例
func New() *Controller {
	return &Controller{}
}

// List 角色列表
func (c *Controller) List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	return service.Role().List(ctx, req)
}

// Create 新增角色
func (c *Controller) Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	return service.Role().Create(ctx, req)
}

// Update 修改角色
func (c *Controller) Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	return service.Role().Update(ctx, req)
}

// Delete 删除角色
func (c *Controller) Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	return service.Role().Delete(ctx, req)
}
