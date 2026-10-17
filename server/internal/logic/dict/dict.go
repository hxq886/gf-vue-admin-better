// Package dict 字典业务实现: 字典类型 / 字典数据 CRUD。
package dict

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/dao"
	"gf-vue-admin-better/internal/model/do"
	"gf-vue-admin-better/internal/model/entity"
	"gf-vue-admin-better/internal/service"
)

// sDict 字典业务实现
type sDict struct{}

func init() {
	service.RegisterDict(&sDict{})
}

// TypeList 字典类型列表。
func (s *sDict) TypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error) {
	res = &v1.DictTypeListRes{List: make([]*entity.DictType, 0)}
	err = dao.DictType.Ctx(ctx).
		OrderAsc(dao.DictType.Columns().Id).
		Scan(&res.List)
	if err != nil {
		return nil, gerror.Wrap(err, "查询字典类型失败")
	}
	return res, nil
}

// TypeCreate 新增字典类型 (类型键唯一)。
func (s *sDict) TypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error) {
	if err = checkTypeUnique(ctx, 0, req.DictType); err != nil {
		return nil, err
	}
	if _, err = dao.DictType.Ctx(ctx).Data(do.DictType{
		Name:     req.Name,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "新增字典类型失败")
	}
	return &v1.DictTypeCreateRes{}, nil
}

// TypeUpdate 修改字典类型 (类型键变更时同步更新字典数据)。
func (s *sDict) TypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error) {
	record, err := dao.DictType.Ctx(ctx).WherePri(req.Id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询字典类型失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("字典类型不存在")
	}
	if err = checkTypeUnique(ctx, req.Id, req.DictType); err != nil {
		return nil, err
	}
	oldType := record[dao.DictType.Columns().DictType].String()

	err = dao.DictType.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.DictType.Ctx(ctx).WherePri(req.Id).Data(do.DictType{
			Name:     req.Name,
			DictType: req.DictType,
			Status:   req.Status,
			Remark:   req.Remark,
		}).Update(); err != nil {
			return gerror.Wrap(err, "修改字典类型失败")
		}
		// 类型键变更时, 同步更新其下字典数据的归属
		if oldType != req.DictType {
			if _, err = dao.DictData.Ctx(ctx).
				Where(dao.DictData.Columns().DictType, oldType).
				Data(do.DictData{DictType: req.DictType}).
				Update(); err != nil {
				return gerror.Wrap(err, "同步字典数据失败")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.DictTypeUpdateRes{}, nil
}

// TypeDelete 批量删除字典类型 (连同其下字典数据)。
func (s *sDict) TypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error) {
	records, err := dao.DictType.Ctx(ctx).WhereIn(dao.DictType.Columns().Id, req.Ids).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询字典类型失败")
	}
	types := records.Array(dao.DictType.Columns().DictType).Strings()

	err = dao.DictType.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.DictType.Ctx(ctx).
			WhereIn(dao.DictType.Columns().Id, req.Ids).
			Delete(); err != nil {
			return gerror.Wrap(err, "删除字典类型失败")
		}
		if len(types) > 0 {
			if _, err = dao.DictData.Ctx(ctx).
				WhereIn(dao.DictData.Columns().DictType, types).
				Delete(); err != nil {
				return gerror.Wrap(err, "删除字典数据失败")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.DictTypeDeleteRes{}, nil
}

// DataList 按类型键查询字典数据。
func (s *sDict) DataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error) {
	res = &v1.DictDataListRes{List: make([]*entity.DictData, 0)}
	err = dao.DictData.Ctx(ctx).
		Where(dao.DictData.Columns().DictType, req.DictType).
		OrderAsc(dao.DictData.Columns().Sort).
		OrderAsc(dao.DictData.Columns().Id).
		Scan(&res.List)
	if err != nil {
		return nil, gerror.Wrap(err, "查询字典数据失败")
	}
	return res, nil
}

// DataCreate 新增字典数据。
func (s *sDict) DataCreate(ctx context.Context, req *v1.DictDataCreateReq) (res *v1.DictDataCreateRes, err error) {
	if _, err = dao.DictData.Ctx(ctx).Data(do.DictData{
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
		TagClass: req.TagClass,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "新增字典数据失败")
	}
	return &v1.DictDataCreateRes{}, nil
}

// DataUpdate 修改字典数据。
func (s *sDict) DataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error) {
	if _, err = dao.DictData.Ctx(ctx).WherePri(req.Id).Data(do.DictData{
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
		TagClass: req.TagClass,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "修改字典数据失败")
	}
	return &v1.DictDataUpdateRes{}, nil
}

// DataDelete 批量删除字典数据。
func (s *sDict) DataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error) {
	if _, err = dao.DictData.Ctx(ctx).
		WhereIn(dao.DictData.Columns().Id, req.Ids).
		Delete(); err != nil {
		return nil, gerror.Wrap(err, "删除字典数据失败")
	}
	return &v1.DictDataDeleteRes{}, nil
}

// checkTypeUnique 校验字典类型键唯一 (excludeId>0 时排除自身)
func checkTypeUnique(ctx context.Context, excludeId uint64, dictType string) error {
	m := dao.DictType.Ctx(ctx).Where(dao.DictType.Columns().DictType, dictType)
	if excludeId > 0 {
		m = m.WhereNot(dao.DictType.Columns().Id, excludeId)
	}
	count, err := m.Count()
	if err != nil {
		return gerror.Wrap(err, "查询字典类型失败")
	}
	if count > 0 {
		return gerror.Newf("字典类型 [%s] 已存在", dictType)
	}
	return nil
}

// 确保 sDict 实现接口 (编译期校验)
var _ service.IDict = (*sDict)(nil)
