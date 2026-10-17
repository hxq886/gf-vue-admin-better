// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id        uint64      `json:"id"        orm:"id"         description:"用户ID"`           // 用户ID
	Username  string      `json:"username"  orm:"username"   description:"登录账号"`           // 登录账号
	Nickname  string      `json:"nickname"  orm:"nickname"   description:"用户昵称"`           // 用户昵称
	Password  string      `json:"password"  orm:"password"   description:"密码(bcrypt加密存储)"` // 密码(bcrypt加密存储)
	Avatar    string      `json:"avatar"    orm:"avatar"     description:"头像地址"`           // 头像地址
	Email     string      `json:"email"     orm:"email"      description:"邮箱"`             // 邮箱
	Mobile    string      `json:"mobile"    orm:"mobile"     description:"手机号"`            // 手机号
	Status    int         `json:"status"    orm:"status"     description:"状态(1启用 0停用)"`    // 状态(1启用 0停用)
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`             // 备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`           // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`           // 更新时间
}
