// Package main 程序入口: 触发内嵌资源注册并启动服务。
package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	// 注册 MySQL 数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"gf-vue-admin-better/internal/cmd"

	// 内嵌资源 (gf pack 生成后自动生效)
	_ "gf-vue-admin-better/internal/packed"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
