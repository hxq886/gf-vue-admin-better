// Package menu 菜单业务实现: 动态路由构建 / 菜单CRUD。
package menu

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/consts"
	"gf-vue-admin-better/internal/dao"
	"gf-vue-admin-better/internal/logic/user"
	"gf-vue-admin-better/internal/model/do"
	"gf-vue-admin-better/internal/model/entity"
	"gf-vue-admin-better/internal/service"
)

// sMenu 菜单业务实现
type sMenu struct{}

func init() {
	service.RegisterMenu(&sMenu{})
}

// GetRouters 获取当前登录用户的动态路由树。
// 超管返回全部启用目录/菜单, 其余按 用户->角色->菜单 关联查询。
func (s *sMenu) GetRouters(ctx context.Context, req *v1.GetRoutersReq) (res *v1.GetRoutersRes, err error) {
	userId, err := user.GetUserId(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := user.GetUserRoles(ctx, userId)
	if err != nil {
		return nil, err
	}
	isAdmin := containsString(roles, consts.RoleKeyAdmin)

	// 查询可见的目录/菜单 (按钮不参与路由)
	var records gdb.Result
	menuTypes := []string{consts.MenuTypeDirectory, consts.MenuTypeMenu}
	if isAdmin {
		records, err = dao.Menu.Ctx(ctx).
			WhereIn(dao.Menu.Columns().MenuType, menuTypes).
			Where(dao.Menu.Columns().Status, 1).
			OrderAsc(dao.Menu.Columns().Sort).
			OrderAsc(dao.Menu.Columns().Id).
			All()
	} else {
		records, err = dao.Menu.Ctx(ctx).As("m").
			InnerJoin(dao.RoleMenu.Table()+" rm", "m.id=rm.menu_id").
			InnerJoin(dao.UserRole.Table()+" ur", "rm.role_id=ur.role_id").
			Fields("m.*").
			Where("ur.user_id", userId).
			WhereIn("m.menu_type", menuTypes).
			Where("m.status", 1).
			Group("m.id").
			OrderAsc("m.sort").
			OrderAsc("m.id").
			All()
	}
	if err != nil {
		return nil, gerror.Wrap(err, "查询菜单失败")
	}

	// 转为实体并构建树 (按 parent_id 分组)
	var menus []*entity.Menu
	if err = records.Structs(&menus); err != nil {
		return nil, gerror.Wrap(err, "菜单数据转换失败")
	}
	childrenMap := make(map[uint64][]*entity.Menu)
	var roots []*entity.Menu
	for _, m := range menus {
		if m.ParentId == 0 {
			roots = append(roots, m)
		} else {
			childrenMap[m.ParentId] = append(childrenMap[m.ParentId], m)
		}
	}

	// 递归组装动态路由
	routers := make(v1.GetRoutersRes, 0, len(roots))
	for _, m := range roots {
		routers = append(routers, toRouter(m, childrenMap))
	}
	return &routers, nil
}

// toRouter 实体菜单 -> 前端路由节点 (字段结构严格适配 vue3-admin-better)
func toRouter(m *entity.Menu, childrenMap map[uint64][]*entity.Menu) *v1.Router {
	router := &v1.Router{
		Path:       m.Path,
		Component:  m.Component,
		Redirect:   m.Redirect,
		Name:       m.Name,
		AlwaysShow: m.AlwaysShow == 1,
		Hidden:     m.Hidden == 1,
		Meta: &v1.RouterMeta{
			Title:       m.Title,
			Icon:        m.Icon,
			Affix:       m.Affix == 1,
			NoKeepAlive: m.KeepAlive == 0,
		},
	}
	if children := childrenMap[m.Id]; len(children) > 0 {
		router.Children = make([]*v1.Router, 0, len(children))
		for _, child := range children {
			router.Children = append(router.Children, toRouter(child, childrenMap))
		}
	}
	return router
}

// List 查询全部菜单 (管理页用, 已按 sort 排序, 前端按 parentId 构建树)。
func (s *sMenu) List(ctx context.Context, req *v1.MenuListReq) (res *v1.MenuListRes, err error) {
	res = &v1.MenuListRes{List: make([]*entity.Menu, 0)}
	err = dao.Menu.Ctx(ctx).
		OrderAsc(dao.Menu.Columns().Sort).
		OrderAsc(dao.Menu.Columns().Id).
		Scan(&res.List)
	if err != nil {
		return nil, gerror.Wrap(err, "查询菜单列表失败")
	}
	return res, nil
}

// Create 新增菜单。
func (s *sMenu) Create(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error) {
	if err = checkParentValid(ctx, req.ParentId); err != nil {
		return nil, err
	}
	if _, err = dao.Menu.Ctx(ctx).Data(do.Menu{
		ParentId:   req.ParentId,
		MenuType:   req.MenuType,
		Title:      req.Title,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Redirect:   req.Redirect,
		AlwaysShow: req.AlwaysShow,
		Hidden:     req.Hidden,
		Icon:       req.Icon,
		KeepAlive:  req.KeepAlive,
		Affix:      req.Affix,
		Permission: req.Permission,
		Sort:       req.Sort,
		Status:     req.Status,
		Remark:     req.Remark,
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "新增菜单失败")
	}
	return &v1.MenuCreateRes{}, nil
}

// Update 修改菜单 (不允许选择自身为上级)。
func (s *sMenu) Update(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	if req.ParentId == req.Id {
		return nil, gerror.New("上级菜单不能选择自己")
	}
	if err = checkParentValid(ctx, req.ParentId); err != nil {
		return nil, err
	}
	record, err := dao.Menu.Ctx(ctx).WherePri(req.Id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询菜单失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("菜单不存在")
	}

	if _, err = dao.Menu.Ctx(ctx).WherePri(req.Id).Data(do.Menu{
		ParentId:   req.ParentId,
		MenuType:   req.MenuType,
		Title:      req.Title,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Redirect:   req.Redirect,
		AlwaysShow: req.AlwaysShow,
		Hidden:     req.Hidden,
		Icon:       req.Icon,
		KeepAlive:  req.KeepAlive,
		Affix:      req.Affix,
		Permission: req.Permission,
		Sort:       req.Sort,
		Status:     req.Status,
		Remark:     req.Remark,
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "修改菜单失败")
	}
	return &v1.MenuUpdateRes{}, nil
}

// Delete 批量删除菜单 (存在子菜单时拒绝), 同步清理角色菜单关联。
func (s *sMenu) Delete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error) {
	// 校验: 任一菜单存在子菜单则拒绝删除
	for _, id := range req.Ids {
		count, err := dao.Menu.Ctx(ctx).
			Where(dao.Menu.Columns().ParentId, id).
			Count()
		if err != nil {
			return nil, gerror.Wrap(err, "查询子菜单失败")
		}
		if count > 0 {
			return nil, gerror.Newf("菜单 [%d] 存在子菜单, 请先删除子菜单", id)
		}
	}

	err = dao.Menu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.Menu.Ctx(ctx).WhereIn(dao.Menu.Columns().Id, req.Ids).Delete(); err != nil {
			return gerror.Wrap(err, "删除菜单失败")
		}
		if _, err = dao.RoleMenu.Ctx(ctx).
			WhereIn(dao.RoleMenu.Columns().MenuId, req.Ids).
			Delete(); err != nil {
			return gerror.Wrap(err, "删除角色菜单关联失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.MenuDeleteRes{}, nil
}

// checkParentValid 校验上级菜单存在 (parentId=0 表示顶级, 直接通过)
func checkParentValid(ctx context.Context, parentId uint64) error {
	if parentId == 0 {
		return nil
	}
	count, err := dao.Menu.Ctx(ctx).WherePri(parentId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询上级菜单失败")
	}
	if count == 0 {
		return gerror.Newf("上级菜单 [%d] 不存在", parentId)
	}
	return nil
}

// containsString 字符串包含判断
func containsString(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// 确保 sMenu 实现接口 (编译期校验)
var _ service.IMenu = (*sMenu)(nil)
