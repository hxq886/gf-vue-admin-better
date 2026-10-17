// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure for table menu.
type Menu struct {
	Id         uint64      `json:"id"         orm:"id"          description:"菜单ID"`                          // 菜单ID
	ParentId   uint64      `json:"parentId"   orm:"parent_id"   description:"父菜单ID(顶级为0)"`                   // 父菜单ID(顶级为0)
	MenuType   string      `json:"menuType"   orm:"menu_type"   description:"类型(M目录 C菜单 F按钮)"`               // 类型(M目录 C菜单 F按钮)
	Title      string      `json:"title"      orm:"title"       description:"菜单标题"`                          // 菜单标题
	Name       string      `json:"name"       orm:"name"        description:"路由name(唯一)"`                    // 路由name(唯一)
	Path       string      `json:"path"       orm:"path"        description:"路由path"`                        // 路由path
	Component  string      `json:"component"  orm:"component"   description:"组件(Layout/EmptyLayout/视图相对路径)"` // 组件(Layout/EmptyLayout/视图相对路径)
	Redirect   string      `json:"redirect"   orm:"redirect"    description:"重定向地址(noRedirect/子路由path)"`     // 重定向地址(noRedirect/子路由path)
	AlwaysShow int         `json:"alwaysShow" orm:"always_show" description:"始终显示根路由"`                       // 始终显示根路由
	Hidden     int         `json:"hidden"     orm:"hidden"      description:"隐藏菜单(1隐藏)"`                     // 隐藏菜单(1隐藏)
	Icon       string      `json:"icon"       orm:"icon"        description:"菜单图标(FontAwesome风格名称)"`         // 菜单图标(FontAwesome风格名称)
	KeepAlive  int         `json:"keepAlive"  orm:"keep_alive"  description:"是否缓存(1缓存)"`                     // 是否缓存(1缓存)
	Affix      int         `json:"affix"      orm:"affix"       description:"是否固定在标签栏"`                      // 是否固定在标签栏
	Permission string      `json:"permission" orm:"permission"  description:"权限标识(按钮: system:user:create)"`  // 权限标识(按钮: system:user:create)
	Sort       int         `json:"sort"       orm:"sort"        description:"显示顺序"`                          // 显示顺序
	Status     int         `json:"status"     orm:"status"      description:"状态(1启用 0停用)"`                   // 状态(1启用 0停用)
	Remark     string      `json:"remark"     orm:"remark"      description:"备注"`                            // 备注
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`                          // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`                          // 更新时间
}
