// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table sys_role for DAO operations like Where/Data.
type Role struct {
	g.Meta    `orm:"table:sys_role, do:true"`
	Id        any         // 角色ID
	Name      any         // 角色名称
	RoleKey   any         // 角色权限字符(接口鉴权用)
	Sort      any         // 显示顺序
	Status    any         // 状态(1启用 0停用)
	Remark    any         // 备注
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
