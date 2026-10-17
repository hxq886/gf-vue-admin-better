// Package v1 定义系统菜单相关接口契约 (动态路由 + 菜单管理)。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-vue-admin-better/internal/model/entity"
)

// GetRoutersReq 获取当前登录用户的动态路由菜单
type GetRoutersReq struct {
	g.Meta `path:"/menu/getRouters" method:"get" tags:"菜单管理" summary:"获取动态路由"`
}

// RouterMeta 动态路由 meta 字段 (title 必填, 其余按需输出)
type RouterMeta struct {
	Title       string `json:"title"                  dc:"菜单标题"`
	Icon        string `json:"icon,omitempty"         dc:"菜单图标"`
	Affix       bool   `json:"affix,omitempty"        dc:"固定在标签栏"`
	NoKeepAlive bool   `json:"noKeepAlive,omitempty"  dc:"不缓存页面"`
}

// Router 动态路由节点 (字段结构严格适配 vue3-admin-better)
type Router struct {
	Path       string      `json:"path"                 dc:"路由路径"`
	Component  string      `json:"component"            dc:"组件(Layout/EmptyLayout/视图相对路径)"`
	Redirect   string      `json:"redirect,omitempty"   dc:"重定向地址"`
	Name       string      `json:"name,omitempty"       dc:"路由名称"`
	AlwaysShow bool        `json:"alwaysShow,omitempty" dc:"始终显示根路由"`
	Hidden     bool        `json:"hidden,omitempty"     dc:"是否隐藏"`
	Children   []*Router   `json:"children,omitempty"   dc:"子路由"`
	Meta       *RouterMeta `json:"meta,omitempty"       dc:"路由元信息"`
}

// GetRoutersRes 动态路由响应 (data 直接输出路由数组, 前端 setAllRoutes 依赖此结构)
type GetRoutersRes []*Router

// MenuListReq 菜单管理-列表 (返回全部菜单, 前端按 parentId 构建树)
type MenuListReq struct {
	g.Meta `path:"/menu/list" method:"get" tags:"菜单管理" summary:"菜单列表"`
}

// MenuListRes 菜单列表响应
type MenuListRes struct {
	List []*entity.Menu `json:"list" dc:"菜单列表(已按sort排序)"`
}

// MenuCreateReq 新增菜单
type MenuCreateReq struct {
	g.Meta     `path:"/menu/create" method:"post" tags:"菜单管理" summary:"新增菜单"`
	ParentId   uint64 `json:"parentId"   dc:"父菜单ID(顶级为0)"`
	MenuType   string `json:"menuType"   v:"required|in:M,C,F#请选择菜单类型|菜单类型不合法" dc:"类型(M目录 C菜单 F按钮)"`
	Title      string `json:"title"      v:"required#请输入菜单标题" dc:"菜单标题"`
	Name       string `json:"name"       dc:"路由name"`
	Path       string `json:"path"       dc:"路由path"`
	Component  string `json:"component"  dc:"组件"`
	Redirect   string `json:"redirect"   dc:"重定向"`
	AlwaysShow *int   `json:"alwaysShow" dc:"始终显示根路由(1是)"`
	Hidden     *int   `json:"hidden"     dc:"隐藏(1是)"`
	Icon       string `json:"icon"       dc:"图标"`
	KeepAlive  *int   `json:"keepAlive"  dc:"缓存(1是)"`
	Affix      *int   `json:"affix"      dc:"固定标签栏(1是)"`
	Permission string `json:"permission" dc:"权限标识"`
	Sort       *int   `json:"sort"       dc:"显示顺序"`
	Status     *int   `json:"status"     dc:"状态(1启用 0停用)" d:"1"`
	Remark     string `json:"remark"     dc:"备注"`
}

// MenuCreateRes 新增菜单响应 (无业务数据)
type MenuCreateRes struct{}

// MenuUpdateReq 修改菜单
type MenuUpdateReq struct {
	g.Meta     `path:"/menu/update" method:"put" tags:"菜单管理" summary:"修改菜单"`
	Id         uint64 `json:"id"         v:"required#缺少菜单ID" dc:"菜单ID"`
	ParentId   uint64 `json:"parentId"   dc:"父菜单ID(顶级为0)"`
	MenuType   string `json:"menuType"   v:"required|in:M,C,F#请选择菜单类型|菜单类型不合法" dc:"类型(M目录 C菜单 F按钮)"`
	Title      string `json:"title"      v:"required#请输入菜单标题" dc:"菜单标题"`
	Name       string `json:"name"       dc:"路由name"`
	Path       string `json:"path"       dc:"路由path"`
	Component  string `json:"component"  dc:"组件"`
	Redirect   string `json:"redirect"   dc:"重定向"`
	AlwaysShow *int   `json:"alwaysShow" dc:"始终显示根路由(1是)"`
	Hidden     *int   `json:"hidden"     dc:"隐藏(1是)"`
	Icon       string `json:"icon"       dc:"图标"`
	KeepAlive  *int   `json:"keepAlive"  dc:"缓存(1是)"`
	Affix      *int   `json:"affix"      dc:"固定标签栏(1是)"`
	Permission string `json:"permission" dc:"权限标识"`
	Sort       *int   `json:"sort"       dc:"显示顺序"`
	Status     *int   `json:"status"     dc:"状态(1启用 0停用)"`
	Remark     string `json:"remark"     dc:"备注"`
}

// MenuUpdateRes 修改菜单响应 (无业务数据)
type MenuUpdateRes struct{}

// MenuDeleteReq 删除菜单 (批量)
type MenuDeleteReq struct {
	g.Meta `path:"/menu/delete" method:"delete" tags:"菜单管理" summary:"删除菜单"`
	Ids    []uint64 `json:"ids" v:"required#请选择要删除的菜单" dc:"菜单ID数组"`
}

// MenuDeleteRes 删除菜单响应 (无业务数据)
type MenuDeleteRes struct{}
