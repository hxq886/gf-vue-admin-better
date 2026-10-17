// Package v1 定义字典相关接口契约 (字典类型 + 字典数据)。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-vue-admin-better/internal/model/entity"
)

// DictTypeListReq 字典类型列表
type DictTypeListReq struct {
	g.Meta `path:"/dict/type/list" method:"get" tags:"字典管理" summary:"字典类型列表"`
}

// DictTypeListRes 字典类型列表响应
type DictTypeListRes struct {
	List []*entity.DictType `json:"list" dc:"字典类型列表"`
}

// DictTypeCreateReq 新增字典类型
type DictTypeCreateReq struct {
	g.Meta    `path:"/dict/type/create" method:"post" tags:"字典管理" summary:"新增字典类型"`
	Name      string `json:"name"     v:"required#请输入字典名称" dc:"字典名称"`
	DictType  string `json:"dictType" v:"required#请输入类型键" dc:"类型键(唯一)"`
	Status    *int   `json:"status"   dc:"状态(1启用 0停用)" d:"1"`
	Remark    string `json:"remark"   dc:"备注"`
}

// DictTypeCreateRes 新增字典类型响应 (无业务数据)
type DictTypeCreateRes struct{}

// DictTypeUpdateReq 修改字典类型 (同步更新字典数据的类型键)
type DictTypeUpdateReq struct {
	g.Meta    `path:"/dict/type/update" method:"put" tags:"字典管理" summary:"修改字典类型"`
	Id        uint64 `json:"id"       v:"required#缺少字典类型ID" dc:"字典类型ID"`
	Name      string `json:"name"     v:"required#请输入字典名称" dc:"字典名称"`
	DictType  string `json:"dictType" v:"required#请输入类型键" dc:"类型键(唯一)"`
	Status    *int   `json:"status"   dc:"状态(1启用 0停用)"`
	Remark    string `json:"remark"   dc:"备注"`
}

// DictTypeUpdateRes 修改字典类型响应 (无业务数据)
type DictTypeUpdateRes struct{}

// DictTypeDeleteReq 删除字典类型 (批量, 同时删除其下字典数据)
type DictTypeDeleteReq struct {
	g.Meta `path:"/dict/type/delete" method:"delete" tags:"字典管理" summary:"删除字典类型"`
	Ids    []uint64 `json:"ids" v:"required#请选择要删除的字典类型" dc:"字典类型ID数组"`
}

// DictTypeDeleteRes 删除字典类型响应 (无业务数据)
type DictTypeDeleteRes struct{}

// DictDataListReq 字典数据列表 (按类型键查询)
type DictDataListReq struct {
	g.Meta    `path:"/dict/data/list" method:"get" tags:"字典管理" summary:"字典数据列表"`
	DictType  string `json:"dictType" v:"required#缺少字典类型" dc:"字典类型键"`
}

// DictDataListRes 字典数据列表响应
type DictDataListRes struct {
	List []*entity.DictData `json:"list" dc:"字典数据列表"`
}

// DictDataCreateReq 新增字典数据
type DictDataCreateReq struct {
	g.Meta    `path:"/dict/data/create" method:"post" tags:"字典管理" summary:"新增字典数据"`
	DictType  string `json:"dictType" v:"required#请选择字典类型" dc:"字典类型键"`
	Label     string `json:"label"    v:"required#请输入标签" dc:"标签(展示用)"`
	Value     string `json:"value"    v:"required#请输入键值" dc:"键值(存储用)"`
	TagClass  string `json:"tagClass" dc:"标签样式(primary/success/warning/danger)"`
	Sort      *int   `json:"sort"     dc:"显示顺序"`
	Status    *int   `json:"status"   dc:"状态(1启用 0停用)" d:"1"`
	Remark    string `json:"remark"   dc:"备注"`
}

// DictDataCreateRes 新增字典数据响应 (无业务数据)
type DictDataCreateRes struct{}

// DictDataUpdateReq 修改字典数据
type DictDataUpdateReq struct {
	g.Meta    `path:"/dict/data/update" method:"put" tags:"字典管理" summary:"修改字典数据"`
	Id        uint64 `json:"id"       v:"required#缺少字典数据ID" dc:"字典数据ID"`
	DictType  string `json:"dictType" v:"required#请选择字典类型" dc:"字典类型键"`
	Label     string `json:"label"    v:"required#请输入标签" dc:"标签(展示用)"`
	Value     string `json:"value"    v:"required#请输入键值" dc:"键值(存储用)"`
	TagClass  string `json:"tagClass" dc:"标签样式(primary/success/warning/danger)"`
	Sort      *int   `json:"sort"     dc:"显示顺序"`
	Status    *int   `json:"status"   dc:"状态(1启用 0停用)"`
	Remark    string `json:"remark"   dc:"备注"`
}

// DictDataUpdateRes 修改字典数据响应 (无业务数据)
type DictDataUpdateRes struct{}

// DictDataDeleteReq 删除字典数据 (批量)
type DictDataDeleteReq struct {
	g.Meta `path:"/dict/data/delete" method:"delete" tags:"字典管理" summary:"删除字典数据"`
	Ids    []uint64 `json:"ids" v:"required#请选择要删除的字典数据" dc:"字典数据ID数组"`
}

// DictDataDeleteRes 删除字典数据响应 (无业务数据)
type DictDataDeleteRes struct{}
