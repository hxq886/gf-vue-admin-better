// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure of table sys_menu for DAO operations like Where/Data.
type Menu struct {
	g.Meta     `orm:"table:sys_menu, do:true"`
	Id         any         // 菜单ID
	ParentId   any         // 父菜单ID(顶级为0)
	MenuType   any         // 类型(M目录 C菜单 F按钮)
	Title      any         // 菜单标题
	Name       any         // 路由name(唯一)
	Path       any         // 路由path
	Component  any         // 组件(Layout/EmptyLayout/视图相对路径)
	Redirect   any         // 重定向地址(noRedirect/子路由path)
	AlwaysShow any         // 始终显示根路由
	Hidden     any         // 隐藏菜单(1隐藏)
	Icon       any         // 菜单图标(FontAwesome风格名称)
	KeepAlive  any         // 是否缓存(1缓存)
	Affix      any         // 是否固定在标签栏
	Permission any         // 权限标识(按钮: system:user:create)
	Sort       any         // 显示顺序
	Status     any         // 状态(1启用 0停用)
	Remark     any         // 备注
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}
