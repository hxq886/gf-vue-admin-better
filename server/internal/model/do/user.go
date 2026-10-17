// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table sys_user for DAO operations like Where/Data.
type User struct {
	g.Meta    `orm:"table:sys_user, do:true"`
	Id        any         // 用户ID
	Username  any         // 登录账号
	Nickname  any         // 用户昵称
	Password  any         // 密码(bcrypt加密存储)
	Avatar    any         // 头像地址
	Email     any         // 邮箱
	Mobile    any         // 手机号
	Status    any         // 状态(1启用 0停用)
	Remark    any         // 备注
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
