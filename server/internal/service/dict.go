// Package service 定义业务层接口, 由 internal/logic 提供实现并注册。
package service

import (
	"context"

	"gf-vue-admin-better/api/system/v1"
)

// IDict 字典服务接口。
type IDict interface {
	// TypeList 字典类型列表
	TypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error)
	// TypeCreate 新增字典类型
	TypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error)
	// TypeUpdate 修改字典类型
	TypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error)
	// TypeDelete 批量删除字典类型 (连同其下字典数据)
	TypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error)
	// DataList 按类型键查询字典数据
	DataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error)
	// DataCreate 新增字典数据
	DataCreate(ctx context.Context, req *v1.DictDataCreateReq) (res *v1.DictDataCreateRes, err error)
	// DataUpdate 修改字典数据
	DataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error)
	// DataDelete 批量删除字典数据
	DataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error)
}

var localIDict IDict

// Dict 获取字典服务实例
func Dict() IDict {
	if localIDict == nil {
		panic("service IDict not registered")
	}
	return localIDict
}

// RegisterDict 注册字典服务实现
func RegisterDict(i IDict) {
	localIDict = i
}
