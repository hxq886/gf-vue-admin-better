// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id        uint64      `json:"id"        orm:"id"         description:"角色ID"`          // 角色ID
	Name      string      `json:"name"      orm:"name"       description:"角色名称"`          // 角色名称
	RoleKey   string      `json:"roleKey"   orm:"role_key"   description:"角色权限字符(接口鉴权用)"` // 角色权限字符(接口鉴权用)
	Sort      int         `json:"sort"      orm:"sort"       description:"显示顺序"`          // 显示顺序
	Status    int         `json:"status"    orm:"status"     description:"状态(1启用 0停用)"`   // 状态(1启用 0停用)
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`            // 备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`          // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`          // 更新时间
}
