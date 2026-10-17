// Package user 用户业务实现: 用户信息 / 用户CRUD / 角色关联。
package user

import (
	"context"

	gfjwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"

	v1 "gf-vue-admin-better/api/system/v1"
	"gf-vue-admin-better/internal/consts"
	"gf-vue-admin-better/internal/dao"
	"gf-vue-admin-better/internal/model/do"
	"gf-vue-admin-better/internal/service"
)

// sUser 用户业务实现
type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

// GetUserId 从请求上下文中提取当前登录用户ID (由 JWT 中间件写入载荷)。
func GetUserId(ctx context.Context) (uint64, error) {
	claims := gfjwt.ExtractClaims(ctx)
	userId := gconv.Uint64(claims["uid"])
	if userId == 0 {
		return 0, gerror.New("获取用户身份失败")
	}
	return userId, nil
}

// GetUserRoles 查询指定用户的角色权限字符数组 (供鉴权与路由构建复用)。
func GetUserRoles(ctx context.Context, userId uint64) (roles []string, err error) {
	values, err := dao.Role.Ctx(ctx).As("r").
		InnerJoin(dao.UserRole.Table()+" ur", "r.id=ur.role_id").
		Fields("r.role_key").
		Where("ur.user_id", userId).
		Where("r.status", 1).
		Array()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户角色失败")
	}
	return values.Strings(), nil
}

// GetInfo 获取当前登录用户信息: 昵称 + 角色 + 按钮权限标识。
// 超管角色(admin)返回全部按钮权限, 其余按角色-菜单关联查询。
func (s *sUser) GetInfo(ctx context.Context) (res *v1.GetInfoRes, err error) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	record, err := dao.User.Ctx(ctx).WherePri(userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("登录用户不存在或已被删除")
	}

	roles, err := GetUserRoles(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 按钮权限标识
	var perms []string
	if gstrContains(roles, consts.RoleKeyAdmin) {
		// 超管: 返回全部按钮权限
		perms, err = allPermissions(ctx)
	} else {
		perms, err = rolePermissions(ctx, userId)
	}
	if err != nil {
		return nil, err
	}

	return &v1.GetInfoRes{
		Username:    record[dao.User.Columns().Username].String(),
		Nickname:    record[dao.User.Columns().Nickname].String(),
		Avatar:      record[dao.User.Columns().Avatar].String(),
		Roles:       roles,
		Permissions: perms,
	}, nil
}

// allPermissions 查询全部启用的按钮权限标识
func allPermissions(ctx context.Context) ([]string, error) {
	values, err := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().MenuType, consts.MenuTypeButton).
		Where(dao.Menu.Columns().Status, 1).
		WhereNot(dao.Menu.Columns().Permission, "").
		Array(dao.Menu.Columns().Permission)
	if err != nil {
		return nil, gerror.Wrap(err, "查询权限标识失败")
	}
	return values.Strings(), nil
}

// rolePermissions 按用户角色关联查询按钮权限标识
func rolePermissions(ctx context.Context, userId uint64) ([]string, error) {
	values, err := dao.Menu.Ctx(ctx).As("m").
		InnerJoin(dao.RoleMenu.Table()+" rm", "m.id=rm.menu_id").
		InnerJoin(dao.UserRole.Table()+" ur", "rm.role_id=ur.role_id").
		Fields("m.permission").
		Where("ur.user_id", userId).
		Where("m.menu_type", consts.MenuTypeButton).
		Where("m.status", 1).
		WhereNot("m.permission", "").
		Array()
	if err != nil {
		return nil, gerror.Wrap(err, "查询权限标识失败")
	}
	return values.Strings(), nil
}

// List 分页查询用户列表 (含角色ID回显)。
func (s *sUser) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	res = &v1.UserListRes{List: make([]*v1.UserItem, 0)}

	m := dao.User.Ctx(ctx)
	if req.Username != "" {
		m = m.WhereLike(dao.User.Columns().Username, "%"+req.Username+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.User.Columns().Status, *req.Status)
	}

	if res.Total, err = m.Count(); err != nil {
		return nil, gerror.Wrap(err, "查询用户总数失败")
	}
	records, err := m.
		OrderAsc(dao.User.Columns().Id).
		Page(req.Page, req.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户列表失败")
	}
	if records.IsEmpty() {
		return res, nil
	}

	// 一次性查询所有用户的角色关联, 避免循环查询
	roleMap, err := userRoleMap(ctx, records)
	if err != nil {
		return nil, err
	}

	if err = records.Structs(&res.List); err != nil {
		return nil, gerror.Wrap(err, "用户数据转换失败")
	}
	for _, item := range res.List {
		if item.RoleIds == nil {
			item.RoleIds = []uint64{}
		}
		if roleMap[item.Id] != nil {
			item.RoleIds = roleMap[item.Id]
		}
	}
	return res, nil
}

// userRoleMap 查询一批用户的角色关联映射 {userId: [roleId...]}
func userRoleMap(ctx context.Context, users gdb.Result) (map[uint64][]uint64, error) {
	roleMap := make(map[uint64][]uint64)
	if users.IsEmpty() {
		return roleMap, nil
	}
	ids := users.Array(dao.User.Columns().Id).Uint64s()
	records, err := dao.UserRole.Ctx(ctx).WhereIn(dao.UserRole.Columns().UserId, ids).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户角色关联失败")
	}
	for _, record := range records {
		roleMap[record[dao.UserRole.Columns().UserId].Uint64()] = append(
			roleMap[record[dao.UserRole.Columns().UserId].Uint64()],
			record[dao.UserRole.Columns().RoleId].Uint64(),
		)
	}
	return roleMap, nil
}

// Create 新增用户: 账号唯一性校验 + 密码加密 + 角色授权 (事务)。
func (s *sUser) Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	// 账号唯一性校验
	count, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().Username, req.Username).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询账号失败")
	}
	if count > 0 {
		return nil, gerror.Newf("账号 [%s] 已存在", req.Username)
	}

	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "密码加密失败")
	}

	// 事务: 新增用户 + 授权角色
	err = dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := dao.User.Ctx(ctx).Data(do.User{
			Username: req.Username,
			Nickname: req.Nickname,
			Password: string(hash),
			Avatar:   req.Avatar,
			Email:    req.Email,
			Mobile:   req.Mobile,
			Status:   req.Status,
			Remark:   req.Remark,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "新增用户失败")
		}
		userId, err := result.LastInsertId()
		if err != nil {
			return gerror.Wrap(err, "获取用户ID失败")
		}
		return grantUserRoles(ctx, uint64(userId), req.RoleIds)
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserCreateRes{}, nil
}

// Update 修改用户基本信息并重新授权角色 (账号与密码不可在此修改)。
func (s *sUser) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	record, err := dao.User.Ctx(ctx).WherePri(req.Id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("用户不存在")
	}

	err = dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.User.Ctx(ctx).WherePri(req.Id).Data(do.User{
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Mobile:   req.Mobile,
			Status:   req.Status,
			Remark:   req.Remark,
		}).Update(); err != nil {
			return gerror.Wrap(err, "修改用户失败")
		}
		// 角色重新授权: 先清空旧关联
		if _, err = dao.UserRole.Ctx(ctx).
			Where(dao.UserRole.Columns().UserId, req.Id).
			Delete(); err != nil {
			return gerror.Wrap(err, "清除用户角色失败")
		}
		return grantUserRoles(ctx, req.Id, req.RoleIds)
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserUpdateRes{}, nil
}

// Password 重置用户密码。
func (s *sUser) Password(ctx context.Context, req *v1.UserPasswordReq) (res *v1.UserPasswordRes, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "密码加密失败")
	}
	if _, err = dao.User.Ctx(ctx).WherePri(req.Id).Data(do.User{
		Password: string(hash),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "重置密码失败")
	}
	return &v1.UserPasswordRes{}, nil
}

// Delete 批量删除用户 (内置管理员受保护), 同步删除角色关联。
func (s *sUser) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	for _, id := range req.Ids {
		if id == consts.UserIdAdmin {
			return nil, gerror.New("内置超级管理员不允许删除")
		}
	}
	err = dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = dao.User.Ctx(ctx).WhereIn(dao.User.Columns().Id, req.Ids).Delete(); err != nil {
			return gerror.Wrap(err, "删除用户失败")
		}
		if _, err = dao.UserRole.Ctx(ctx).
			WhereIn(dao.UserRole.Columns().UserId, req.Ids).
			Delete(); err != nil {
			return gerror.Wrap(err, "删除用户角色关联失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserDeleteRes{}, nil
}

// grantUserRoles 写入用户-角色关联 (空数组跳过)
func grantUserRoles(ctx context.Context, userId uint64, roleIds []uint64) error {
	if len(roleIds) == 0 {
		return nil
	}
	list := g.List{}
	for _, roleId := range roleIds {
		list = append(list, g.Map{
			dao.UserRole.Columns().UserId: userId,
			dao.UserRole.Columns().RoleId: roleId,
		})
	}
	if _, err := dao.UserRole.Ctx(ctx).Data(list).Insert(); err != nil {
		return gerror.Wrap(err, "授权用户角色失败")
	}
	return nil
}

// gstrContains 简易包含判断 (避免为一次调用引入依赖)
func gstrContains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// 确保 sUser 实现接口 (编译期校验)
var _ service.IUser = (*sUser)(nil)
