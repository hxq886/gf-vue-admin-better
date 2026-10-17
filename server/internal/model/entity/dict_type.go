// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DictType is the golang structure for table dict_type.
type DictType struct {
	Id        uint64      `json:"id"        orm:"id"         description:"字典类型ID"`      // 字典类型ID
	Name      string      `json:"name"      orm:"name"       description:"字典名称"`        // 字典名称
	DictType  string      `json:"dictType"  orm:"dict_type"  description:"类型键(唯一)"`     // 类型键(唯一)
	Status    int         `json:"status"    orm:"status"     description:"状态(1启用 0停用)"` // 状态(1启用 0停用)
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`          // 备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`        // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`        // 更新时间
}
