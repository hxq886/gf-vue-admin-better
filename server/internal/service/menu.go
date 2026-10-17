// Package service 定义业务层接口, 由 internal/logic 提供实现并注册。
package service

import (
	"context"

	"gf-vue-admin-better/api/system/v1"
)

// IMenu 菜单服务接口。
type IMenu interface {
	// GetRouters 获取当前登录用户的动态路由树 (严格适配 vue3-admin-better)
	GetRouters(ctx context.Context, req *v1.GetRoutersReq) (res *v1.GetRoutersRes, err error)
	// List 查询全部菜单 (管理页用)
	List(ctx context.Context, req *v1.MenuListReq) (res *v1.MenuListRes, err error)
	// Create 新增菜单
	Create(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error)
	// Update 修改菜单
	Update(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error)
	// Delete 批量删除菜单 (存在子菜单时拒绝)
	Delete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error)
}

var localIMenu IMenu

// Menu 获取菜单服务实例
func Menu() IMenu {
	if localIMenu == nil {
		panic("service IMenu not registered")
	}
	return localIMenu
}

// RegisterMenu 注册菜单服务实现
func RegisterMenu(i IMenu) {
	localIMenu = i
}
