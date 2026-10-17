// Package service 定义业务层接口, 由 internal/logic 提供实现并注册。
package service

import (
	"context"

	"gf-vue-admin-better/api/system/v1"
)

// IRole 角色服务接口。
type IRole interface {
	// List 查询角色列表 (含已授权菜单ID)
	List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error)
	// Create 新增角色并授权菜单
	Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error)
	// Update 修改角色并重授菜单
	Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error)
	// Delete 批量删除角色
	Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error)
}

var localIRole IRole

// Role 获取角色服务实例
func Role() IRole {
	if localIRole == nil {
		panic("service IRole not registered")
	}
	return localIRole
}

// RegisterRole 注册角色服务实现
func RegisterRole(i IRole) {
	localIRole = i
}
