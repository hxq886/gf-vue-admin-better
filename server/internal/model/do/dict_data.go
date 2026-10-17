// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DictData is the golang structure of table sys_dict_data for DAO operations like Where/Data.
type DictData struct {
	g.Meta    `orm:"table:sys_dict_data, do:true"`
	Id        any         // 字典数据ID
	DictType  any         // 所属类型键
	Label     any         // 标签(展示用)
	Value     any         // 键值(存储用)
	TagClass  any         // 标签样式(primary/success/warning/danger)
	Sort      any         // 显示顺序
	Status    any         // 状态(1启用 0停用)
	Remark    any         // 备注
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
