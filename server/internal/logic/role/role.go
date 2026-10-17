// Package role 角色业务实现: 角色CRUD / 菜单授权。
package role

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/consts"
	"gf-vue-admin-better/internal/dao"
	"gf-vue-admin-better/internal/model/do"
	"gf-vue-admin-better/internal/service"
)

// sRole 角色业务实现
type sRole struct{}

func init() {
	service.RegisterRole(&sRole{})
}

// List 查询全部角色 (含已授权菜单ID, 供角色编辑回显)。
func (s *sRole) List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	res = &v1.RoleListRes{List: make([]*v1.RoleItem, 0)}

	records, err := dao.Role.Ctx(ctx).
		OrderAsc(dao.Role.Columns().Sort).
		OrderAsc(dao.Role.Columns().Id).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色列表失败")
	}
	if records.IsEmpty() {
		return res, nil
	}

	if err = records.Structs(&res.List); err != nil {
		return nil, gerror.Wrap(err, "角色数据转换失败")
	}

	// 一次性查询角色-菜单关联
	menuMap := make(map[uint64][]uint64)
	roleIds := records.Array(dao.Role.Columns().Id).Uint64s()
	rmRecords, err := dao.RoleMenu.Ctx(ctx).
		WhereIn(dao.RoleMenu.Columns().RoleId, roleIds).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色菜单关联失败")
	}
	for _, record := range rmRecords {
		roleId := record[dao.RoleMenu.Columns().RoleId].Uint64()
		menuMap[roleId] = append(menuMap[roleId], record[dao.RoleMenu.Columns().MenuId].Uint64())
	}
	for _, item := range res.List {
		if menuMap[item.Id] == nil {
			item.MenuIds = []uint64{}
		} else {
			item.MenuIds = menuMap[item.Id]
		}
	}
	return res, nil
}

// Create 新增角色并授权菜单 (角色权限字符唯一)。
func (s *sRole) Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	if err = checkRoleKeyUnique(ctx, 0, req.RoleKey); err != nil {
		return nil, err
	}

	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := dao.Role.Ctx(ctx).Data(do.Role{
			Name:    req.Name,
			RoleKey: req.RoleKey,
			Sort:    req.Sort,
			Status:  req.Status,
			Remark:  req.Remark,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "新增角色失败")
		}
		roleId, err := result.LastInsertId()
		if err != nil {
			return gerror.Wrap(err, "获取角色ID失败")
		}
		return grantRoleMenus(ctx, uint64(roleId), req.MenuIds)
	})
	if err != nil {
		return nil, err
	}
	return &v1.RoleCreateRes{}, nil
}

// Update 修改角色并重新授权菜单 (内置超管角色受保护)。
func (s *sRole) Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	if req.Id == consts.RoleSuper {
		return nil, gerror.New("内置超级管理员角色不允许修改")
	}
	record, err := dao.Role.Ctx(ctx).WherePri(req.Id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("角色不存在")
	}
	if err = checkRoleKeyUnique(ctx, req.Id, req.RoleKey); err != nil {
		return nil, err
	}

	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.Role.Ctx(ctx).WherePri(req.Id).Data(do.Role{
			Name:    req.Name,
			RoleKey: req.RoleKey,
			Sort:    req.Sort,
			Status:  req.Status,
			Remark:  req.Remark,
		}).Update(); err != nil {
			return gerror.Wrap(err, "修改角色失败")
		}
		// 菜单重新授权: 先清空旧关联
		if _, err = dao.RoleMenu.Ctx(ctx).
			Where(dao.RoleMenu.Columns().RoleId, req.Id).
			Delete(); err != nil {
			return gerror.Wrap(err, "清除角色菜单失败")
		}
		return grantRoleMenus(ctx, req.Id, req.MenuIds)
	})
	if err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}

// Delete 批量删除角色 (内置超管受保护), 同步清理角色菜单/用户角色关联。
func (s *sRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	for _, id := range req.Ids {
		if id == consts.RoleSuper {
			return nil, gerror.New("内置超级管理员角色不允许删除")
		}
	}
	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.Role.Ctx(ctx).WhereIn(dao.Role.Columns().Id, req.Ids).Delete(); err != nil {
			return gerror.Wrap(err, "删除角色失败")
		}
		if _, err = dao.RoleMenu.Ctx(ctx).
			WhereIn(dao.RoleMenu.Columns().RoleId, req.Ids).
			Delete(); err != nil {
			return gerror.Wrap(err, "删除角色菜单关联失败")
		}
		if _, err = dao.UserRole.Ctx(ctx).
			WhereIn(dao.UserRole.Columns().RoleId, req.Ids).
			Delete(); err != nil {
			return gerror.Wrap(err, "删除用户角色关联失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}

// checkRoleKeyUnique 校验角色权限字符唯一 (excludeId>0 时排除自身)
func checkRoleKeyUnique(ctx context.Context, excludeId uint64, roleKey string) error {
	m := dao.Role.Ctx(ctx).Where(dao.Role.Columns().RoleKey, roleKey)
	if excludeId > 0 {
		m = m.WhereNot(dao.Role.Columns().Id, excludeId)
	}
	count, err := m.Count()
	if err != nil {
		return gerror.Wrap(err, "查询角色权限字符失败")
	}
	if count > 0 {
		return gerror.Newf("角色权限字符 [%s] 已存在", roleKey)
	}
	return nil
}

// grantRoleMenus 写入角色-菜单关联 (空数组跳过)
func grantRoleMenus(ctx context.Context, roleId uint64, menuIds []uint64) error {
	if len(menuIds) == 0 {
		return nil
	}
	list := g.List{}
	for _, menuId := range menuIds {
		list = append(list, g.Map{
			dao.RoleMenu.Columns().RoleId: roleId,
			dao.RoleMenu.Columns().MenuId: menuId,
		})
	}
	if _, err := dao.RoleMenu.Ctx(ctx).Data(list).Insert(); err != nil {
		return gerror.Wrap(err, "授权角色菜单失败")
	}
	return nil
}

// 确保 sRole 实现接口 (编译期校验)
var _ service.IRole = (*sRole)(nil)
