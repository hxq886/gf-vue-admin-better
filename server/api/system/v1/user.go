// Package v1 定义系统用户相关接口契约。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GetInfoReq 获取当前登录用户信息
type GetInfoReq struct {
	g.Meta `path:"/user/getInfo" method:"get" tags:"用户管理" summary:"获取登录用户信息"`
}

// GetInfoRes 用户信息响应 (roles/permissions 驱动前端动态路由与按钮权限)
type GetInfoRes struct {
	Username    string   `json:"username"    dc:"登录账号"`
	Nickname    string   `json:"nickname"    dc:"用户昵称"`
	Avatar      string   `json:"avatar"      dc:"头像地址"`
	Roles       []string `json:"roles"       dc:"角色权限字符数组"`
	Permissions []string `json:"permissions" dc:"按钮权限标识数组"`
}

// UserListReq 用户列表查询
type UserListReq struct {
	g.Meta   `path:"/user/list" method:"get" tags:"用户管理" summary:"用户列表"`
	Page     int    `json:"page"     dc:"页码"        d:"1"`
	PageSize int    `json:"pageSize" dc:"每页条数"    d:"10"`
	Username string `json:"username" dc:"账号(模糊查询)"`
	Status   *int   `json:"status"   dc:"状态(1启用 0停用)"`
}

// UserItem 用户列表项
type UserItem struct {
	Id        uint64      `json:"id"        dc:"用户ID"`
	Username  string      `json:"username"  dc:"登录账号"`
	Nickname  string      `json:"nickname"  dc:"用户昵称"`
	Avatar    string      `json:"avatar"    dc:"头像"`
	Email     string      `json:"email"     dc:"邮箱"`
	Mobile    string      `json:"mobile"    dc:"手机号"`
	Status    int         `json:"status"    dc:"状态(1启用 0停用)"`
	Remark    string      `json:"remark"    dc:"备注"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"创建时间"`
	RoleIds   []uint64    `json:"roleIds"   dc:"拥有角色ID数组"`
}

// UserListRes 用户列表响应
type UserListRes struct {
	Total int        `json:"total" dc:"总记录数"`
	List  []*UserItem `json:"list"  dc:"用户列表"`
}

// UserCreateReq 新增用户
type UserCreateReq struct {
	g.Meta   `path:"/user/create" method:"post" tags:"用户管理" summary:"新增用户"`
	Username string   `json:"username" v:"required|length:2,64#请输入用户名|用户名长度为2-64位" dc:"登录账号"`
	Nickname string   `json:"nickname" v:"required#请输入昵称" dc:"用户昵称"`
	Password string   `json:"password" v:"required|length:6,32#请输入密码|密码长度为6-32位" dc:"登录密码"`
	Avatar   string   `json:"avatar"   dc:"头像地址"`
	Email    string   `json:"email"    v:"email#邮箱格式不正确" dc:"邮箱"`
	Mobile   string   `json:"mobile"   dc:"手机号"`
	Status   *int     `json:"status"   dc:"状态(1启用 0停用)" d:"1"`
	Remark   string   `json:"remark"   dc:"备注"`
	RoleIds  []uint64 `json:"roleIds"  dc:"角色ID数组"`
}

// UserCreateRes 新增用户响应 (无业务数据)
type UserCreateRes struct{}

// UserUpdateReq 修改用户 (账号不可修改, 密码请走重置接口)
type UserUpdateReq struct {
	g.Meta   `path:"/user/update" method:"put" tags:"用户管理" summary:"修改用户"`
	Id       uint64   `json:"id"       v:"required#缺少用户ID" dc:"用户ID"`
	Nickname string   `json:"nickname" v:"required#请输入昵称" dc:"用户昵称"`
	Avatar   string   `json:"avatar"   dc:"头像地址"`
	Email    string   `json:"email"    v:"email#邮箱格式不正确" dc:"邮箱"`
	Mobile   string   `json:"mobile"   dc:"手机号"`
	Status   *int     `json:"status"   dc:"状态(1启用 0停用)"`
	Remark   string   `json:"remark"   dc:"备注"`
	RoleIds  []uint64 `json:"roleIds"  dc:"角色ID数组"`
}

// UserUpdateRes 修改用户响应 (无业务数据)
type UserUpdateRes struct{}

// UserPasswordReq 重置用户密码
type UserPasswordReq struct {
	g.Meta   `path:"/user/password" method:"put" tags:"用户管理" summary:"重置密码"`
	Id       uint64 `json:"id"       v:"required#缺少用户ID" dc:"用户ID"`
	Password string `json:"password" v:"required|length:6,32#请输入新密码|密码长度为6-32位" dc:"新密码"`
}

// UserPasswordRes 重置密码响应 (无业务数据)
type UserPasswordRes struct{}

// UserDeleteReq 删除用户 (批量)
type UserDeleteReq struct {
	g.Meta `path:"/user/delete" method:"delete" tags:"用户管理" summary:"删除用户"`
	Ids    []uint64 `json:"ids" v:"required#请选择要删除的用户" dc:"用户ID数组"`
}

// UserDeleteRes 删除用户响应 (无业务数据)
type UserDeleteRes struct{}
