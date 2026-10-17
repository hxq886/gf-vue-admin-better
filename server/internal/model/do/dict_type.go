// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DictType is the golang structure of table sys_dict_type for DAO operations like Where/Data.
type DictType struct {
	g.Meta    `orm:"table:sys_dict_type, do:true"`
	Id        any         // 字典类型ID
	Name      any         // 字典名称
	DictType  any         // 类型键(唯一)
	Status    any         // 状态(1启用 0停用)
	Remark    any         // 备注
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
