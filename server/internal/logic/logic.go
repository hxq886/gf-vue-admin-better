// Package logic 业务实现聚合入口。
// 通过空白导入触发各业务包的 init 注册, 保证 service 接口在启动时全部就绪。
package logic

import (
	// 认证业务: 登录/登出
	_ "gf-vue-admin-better/internal/logic/auth"
	// 字典业务: 字典类型/字典数据
	_ "gf-vue-admin-better/internal/logic/dict"
	// 菜单业务: 动态路由/菜单管理
	_ "gf-vue-admin-better/internal/logic/menu"
	// 角色业务: 角色管理/菜单授权
	_ "gf-vue-admin-better/internal/logic/role"
	// 用户业务: 用户信息/用户管理
	_ "gf-vue-admin-better/internal/logic/user"
)
