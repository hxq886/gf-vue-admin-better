// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MenuDao is the data access object for the table sys_menu.
type MenuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MenuColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MenuColumns defines and stores column names for the table sys_menu.
type MenuColumns struct {
	Id         string // 菜单ID
	ParentId   string // 父菜单ID(顶级为0)
	MenuType   string // 类型(M目录 C菜单 F按钮)
	Title      string // 菜单标题
	Name       string // 路由name(唯一)
	Path       string // 路由path
	Component  string // 组件(Layout/EmptyLayout/视图相对路径)
	Redirect   string // 重定向地址(noRedirect/子路由path)
	AlwaysShow string // 始终显示根路由
	Hidden     string // 隐藏菜单(1隐藏)
	Icon       string // 菜单图标(FontAwesome风格名称)
	KeepAlive  string // 是否缓存(1缓存)
	Affix      string // 是否固定在标签栏
	Permission string // 权限标识(按钮: system:user:create)
	Sort       string // 显示顺序
	Status     string // 状态(1启用 0停用)
	Remark     string // 备注
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// menuColumns holds the columns for the table sys_menu.
var menuColumns = MenuColumns{
	Id:         "id",
	ParentId:   "parent_id",
	MenuType:   "menu_type",
	Title:      "title",
	Name:       "name",
	Path:       "path",
	Component:  "component",
	Redirect:   "redirect",
	AlwaysShow: "always_show",
	Hidden:     "hidden",
	Icon:       "icon",
	KeepAlive:  "keep_alive",
	Affix:      "affix",
	Permission: "permission",
	Sort:       "sort",
	Status:     "status",
	Remark:     "remark",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewMenuDao creates and returns a new DAO object for table data access.
func NewMenuDao(handlers ...gdb.ModelHandler) *MenuDao {
	return &MenuDao{
		group:    "default",
		table:    "sys_menu",
		columns:  menuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MenuDao) Columns() MenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MenuDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
