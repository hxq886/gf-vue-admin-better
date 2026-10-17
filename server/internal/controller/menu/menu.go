// Package menu 菜单控制器: 动态路由 / 菜单CRUD。
package menu

import (
	"context"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/service"
)

// Controller 菜单控制器
type Controller struct{}

// New 创建菜单控制器实例
func New() *Controller {
	return &Controller{}
}

// GetRouters 获取当前登录用户的动态路由
func (c *Controller) GetRouters(ctx context.Context, req *v1.GetRoutersReq) (res *v1.GetRoutersRes, err error) {
	return service.Menu().GetRouters(ctx, req)
}

// List 菜单列表
func (c *Controller) List(ctx context.Context, req *v1.MenuListReq) (res *v1.MenuListRes, err error) {
	return service.Menu().List(ctx, req)
}

// Create 新增菜单
func (c *Controller) Create(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error) {
	return service.Menu().Create(ctx, req)
}

// Update 修改菜单
func (c *Controller) Update(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	return service.Menu().Update(ctx, req)
}

// Delete 删除菜单
func (c *Controller) Delete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error) {
	return service.Menu().Delete(ctx, req)
}
