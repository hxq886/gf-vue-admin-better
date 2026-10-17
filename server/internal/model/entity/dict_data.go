// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DictData is the golang structure for table dict_data.
type DictData struct {
	Id        uint64      `json:"id"        orm:"id"         description:"字典数据ID"`                               // 字典数据ID
	DictType  string      `json:"dictType"  orm:"dict_type"  description:"所属类型键"`                                // 所属类型键
	Label     string      `json:"label"     orm:"label"      description:"标签(展示用)"`                              // 标签(展示用)
	Value     string      `json:"value"     orm:"value"      description:"键值(存储用)"`                              // 键值(存储用)
	TagClass  string      `json:"tagClass"  orm:"tag_class"  description:"标签样式(primary/success/warning/danger)"` // 标签样式(primary/success/warning/danger)
	Sort      int         `json:"sort"      orm:"sort"       description:"显示顺序"`                                 // 显示顺序
	Status    int         `json:"status"    orm:"status"     description:"状态(1启用 0停用)"`                          // 状态(1启用 0停用)
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`                                   // 备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                                 // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                                 // 更新时间
}
