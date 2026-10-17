# gf-vue-admin-better

基于 **GoFrame V2**（后端）+ **Vue3-Admin-Better**（前端）打造的标准化、可商用、开箱即用的前后端分离中后台脚手。


- 纯前后端分离架构，工程完全解耦，前端保留原生升级能力
- 生产模式前端资源经 `gf pack` 内嵌，**单二进制一键部署**，无需 Nginx
- 内置完整 RBAC 权限体系与初始化 SQL，`admin / 123456` 开箱即用
- MIT 开源协议，可自由商用

## 功能总览

| 模块 | 能力 |
| --- | --- |
| 登录鉴权 | JWT 令牌（登出黑名单机制）、令牌过期自动登出 |
| 权限体系 | RBAC 用户/角色/菜单、后端驱动动态路由、按钮级权限标识 |
| 系统管理 | 用户管理、角色管理、菜单管理、字典管理（CRUD 完整） |
| 工程能力 | GF 代码生成器（dao/model/entity）、Makefile 一键脚本 |
| 部署能力 | 源码部署 / 单二进制部署 / Docker Compose 一键部署 |
| 安全能力 | bcrypt 密码、请求参数校验、CORS 开关、内置超管保护 |

## 技术栈

- 后端：GoFrame v2（最新稳定版）+ gf-jwt + MySQL（>=5.7，推荐 8.0）
- 前端：Vue 3 + Vue Router + Vuex + Element Plus + rspack（Node >= 22）

## 快速开始

```bash
# 1. 克隆
git clone https://gitee.com/your-name/gf-vue-admin-better.git
cd gf-vue-admin-better

# 2. 初始化数据库 (建库 + RBAC 表 + 种子数据)
#    默认连接 127.0.0.1:3306 root 空密码, 可覆盖:
#    make init-db DB_HOST=... DB_PORT=... DB_USER=... DB_PASSWORD=...
make init-db

# 3. 修改后端数据库连接
#    编辑 server/manifest/config/config.yaml 的 database.default.link

# 4. 启动 (两个终端, 或 make dev 并行)
make dev-server    # 后端 http://127.0.0.1:8000
make dev-web       # 前端 http://127.0.0.1:8091 (/api 代理到后端)
```

浏览器访问 `http://127.0.0.1:8091`，使用 `admin / 123456` 登录。

## 一键打包单二进制

```bash
make build
# 产物 bin/gf-vue-admin-better (前端资源已内嵌)
# 服务器上仅需: 二进制 + manifest/config/config.yaml, 直接运行即可
```

## Docker 一键部署

```bash
make docker-up    # 应用 + MySQL 编排, 首次启动自动初始化数据库
```

## 常用命令

```bash
make help                          # 全部命令列表
make gen-dao                       # 按表结构生成 dao/model/entity
make gen-password PASSWORD=xxx     # 生成 bcrypt 密码哈希
make docker-down                   # 停止 Docker 服务
make clean                         # 清理构建产物
```

## 目录结构

```
gf-vue-admin-better/
├── web/                        # Vue3-Admin-Better 前端 (rspack)
├── server/                     # GoFrame v2 后端 (官方标准分层)
│   ├── api/                    # 接口定义 (登录/用户/角色/菜单/字典)
│   ├── internal/               # cmd/controller/service/logic/dao/model
│   ├── resource/sql/init.sql   # RBAC 初始化脚本
│   └── Makefile                # 后端专属脚本 (pack/gen-dao/init-db)
├── docs/                       # 部署 / 二次开发 / 接口文档
├── docker/                     # Dockerfile + docker-compose
├── Makefile                    # 根目录一键脚本
└── LICENSE                     # MIT
```

## 文档

- [编译部署详细教程](docs/build-deploy.md)：环境安装、单二进制打包、服务器部署、Docker、交叉编译、故障排查
- [部署指南](docs/deploy.md)：源码 / 单二进制 / Docker Compose / Nginx
- [二次开发指南](docs/development.md)：工程结构、新增业务模块全流程、权限体系
- [接口文档](docs/api.md)：统一响应结构、全部接口契约

## 说明

- 前端基于开源项目 [vue-admin-better](https://gitee.com/chu1204505056/vue-admin-better)（原作者保留版权），本项目仅做配置层适配，未深度修改原生源码
- 内置账号 `admin` 为超管（受接口保护不可删除），上线后请立即修改密码
