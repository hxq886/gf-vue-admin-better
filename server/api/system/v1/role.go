// Package v1 定义系统角色相关接口契约。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleListReq 角色列表查询 (角色数据量小, 不分页)
type RoleListReq struct {
	g.Meta `path:"/role/list" method:"get" tags:"角色管理" summary:"角色列表"`
}

// RoleItem 角色列表项 (含已授权菜单, 供角色编辑回显)
type RoleItem struct {
	Id        uint64      `json:"id"        dc:"角色ID"`
	Name      string      `json:"name"      dc:"角色名称"`
	RoleKey   string      `json:"roleKey"   dc:"角色权限字符"`
	Sort      int         `json:"sort"      dc:"显示顺序"`
	Status    int         `json:"status"    dc:"状态(1启用 0停用)"`
	Remark    string      `json:"remark"    dc:"备注"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"创建时间"`
	MenuIds   []uint64    `json:"menuIds"   dc:"已授权菜单ID数组"`
}

// RoleListRes 角色列表响应
type RoleListRes struct {
	List []*RoleItem `json:"list" dc:"角色列表"`
}

// RoleCreateReq 新增角色
type RoleCreateReq struct {
	g.Meta  `path:"/role/create" method:"post" tags:"角色管理" summary:"新增角色"`
	Name    string   `json:"name"    v:"required#请输入角色名称" dc:"角色名称"`
	RoleKey string   `json:"roleKey" v:"required#请输入角色权限字符" dc:"角色权限字符"`
	Sort    *int     `json:"sort"    dc:"显示顺序"`
	Status  *int     `json:"status"  dc:"状态(1启用 0停用)" d:"1"`
	Remark  string   `json:"remark"  dc:"备注"`
	MenuIds []uint64 `json:"menuIds" dc:"授权菜单ID数组"`
}

// RoleCreateRes 新增角色响应 (无业务数据)
type RoleCreateRes struct{}

// RoleUpdateReq 修改角色
type RoleUpdateReq struct {
	g.Meta  `path:"/role/update" method:"put" tags:"角色管理" summary:"修改角色"`
	Id      uint64   `json:"id"      v:"required#缺少角色ID" dc:"角色ID"`
	Name    string   `json:"name"    v:"required#请输入角色名称" dc:"角色名称"`
	RoleKey string   `json:"roleKey" v:"required#请输入角色权限字符" dc:"角色权限字符"`
	Sort    *int     `json:"sort"    dc:"显示顺序"`
	Status  *int     `json:"status"  dc:"状态(1启用 0停用)"`
	Remark  string   `json:"remark"  dc:"备注"`
	MenuIds []uint64 `json:"menuIds" dc:"授权菜单ID数组"`
}

// RoleUpdateRes 修改角色响应 (无业务数据)
type RoleUpdateRes struct{}

// RoleDeleteReq 删除角色 (批量)
type RoleDeleteReq struct {
	g.Meta `path:"/role/delete" method:"delete" tags:"角色管理" summary:"删除角色"`
	Ids    []uint64 `json:"ids" v:"required#请选择要删除的角色" dc:"角色ID数组"`
}

// RoleDeleteRes 删除角色响应 (无业务数据)
type RoleDeleteRes struct{}
