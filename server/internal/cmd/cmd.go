// Package cmd 服务启动入口: 路由注册 / 中间件挂载 / 静态资源托管。
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"

	"gf-vue-admin-better/internal/controller/auth"
	"gf-vue-admin-better/internal/controller/dict"
	"gf-vue-admin-better/internal/controller/menu"
	"gf-vue-admin-better/internal/controller/role"
	"gf-vue-admin-better/internal/controller/user"

	// 触发业务实现注册 (service 接口就绪)
	_ "gf-vue-admin-better/internal/logic"

	"gf-vue-admin-better/internal/middleware"
)

// distDir 前端构建产物目录 (源码部署时读取磁盘, 单二进制部署时读取内嵌资源)
var distDir = gfile.Join("resource", "public", "dist")

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化 JWT 鉴权实例 (读取 config.yaml 的 jwt 节点)
			middleware.InitAuth()

			s := g.Server()

			// ------------------------------------------------------------------
			// API 路由: 统一 /api 前缀, 与前端 baseURL 配置对齐
			// ------------------------------------------------------------------
			s.Group("/api", func(group *ghttp.RouterGroup) {
				// 跨域放行 (config.yaml 的 cors.enabled 控制, 生产同源部署可关闭)
				if g.Cfg().MustGet(ctx, "cors.enabled", true).Bool() {
					group.Middleware(middleware.CORS)
				}
				group.Middleware(middleware.Response) // 统一响应 {code,msg,data}

				// 公开接口: 登录 / 登出 (无需令牌)
				group.Bind(auth.New())

				// 鉴权接口: 携带有效 JWT 令牌方可访问
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.MiddlewareAuth)
					group.Bind(
						user.New(), // 用户: getInfo/list/create/update/password/delete
						role.New(), // 角色: list/create/update/delete
						menu.New(), // 菜单: getRouters/list/create/update/delete
						dict.New(), // 字典: type/*/data/*
					)
				})
			})

			// ------------------------------------------------------------------
			// 前端静态资源托管 + SPA 404 回退
			// 开发模式: 前端由 rspack devServer(8091) 独立运行, 代理 /api 到本服务;
			// 生产模式: 托管 resource/public/dist, 非接口路径刷新统一回退 index.html。
			// ------------------------------------------------------------------
			setupStatic(s)

			s.Run()
			return nil
		},
	}
)

// setupStatic 托管前端构建产物并注册 SPA 404 回退。
// 目录不存在时跳过托管 (开发模式), 不影响 API 服务。
func setupStatic(s *ghttp.Server) {
	if gfile.Exists(distDir) || gres.Contains(distDir) {
		s.SetServerRoot(distDir)
	}
	// SPA 路由回退: 所有非 /api 请求的 404 统一返回 index.html, 解决前端刷新 404
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		if gstr.HasPrefix(r.URL.Path, "/api") {
			r.Response.WriteJson(g.Map{"code": 404, "msg": "接口不存在", "data": nil})
			return
		}
		indexPath := gfile.Join(distDir, "index.html")
		if gfile.Exists(indexPath) || gres.Contains(indexPath) {
			r.Response.ServeFile(indexPath)
			return
		}
		r.Response.WriteJson(g.Map{"code": 404, "msg": "资源不存在", "data": nil})
	})
}
