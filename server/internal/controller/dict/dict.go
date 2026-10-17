// Package dict 字典控制器: 字典类型 / 字典数据 CRUD。
package dict

import (
	"context"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/service"
)

// Controller 字典控制器
type Controller struct{}

// New 创建字典控制器实例
func New() *Controller {
	return &Controller{}
}

// TypeList 字典类型列表
func (c *Controller) TypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error) {
	return service.Dict().TypeList(ctx, req)
}

// TypeCreate 新增字典类型
func (c *Controller) TypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error) {
	return service.Dict().TypeCreate(ctx, req)
}

// TypeUpdate 修改字典类型
func (c *Controller) TypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error) {
	return service.Dict().TypeUpdate(ctx, req)
}

// TypeDelete 删除字典类型
func (c *Controller) TypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error) {
	return service.Dict().TypeDelete(ctx, req)
}

// DataList 字典数据列表
func (c *Controller) DataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error) {
	return service.Dict().DataList(ctx, req)
}

// DataCreate 新增字典数据
func (c *Controller) DataCreate(ctx context.Context, req *v1.DictDataCreateReq) (res *v1.DictDataCreateRes, err error) {
	return service.Dict().DataCreate(ctx, req)
}

// DataUpdate 修改字典数据
func (c *Controller) DataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error) {
	return service.Dict().DataUpdate(ctx, req)
}

// DataDelete 删除字典数据
func (c *Controller) DataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error) {
	return service.Dict().DataDelete(ctx, req)
}
