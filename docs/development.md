# 二次开发指南

面向基于本脚手架做业务开发的工程师，涵盖工程结构、开发流程、代码生成、权限体系扩展。

---

## 一、工程总览

```
gf-vue-admin-better/
├── web/                          # Vue3-Admin-Better 前端 (rspack 构建, 原生能力保留)
│   ├── src/api/                  # 接口封装 (user.js/router.js/system.js)
│   ├── src/config/               # 环境与行为配置 (net/setting/theme)
│   ├── src/views/vab/system/     # 系统管理页面 (用户/角色/菜单/字典)
│   ├── .env.development          # 开发环境变量 (mock 开关等)
│   └── rspack.config.js          # 构建配置 (含 /api 代理)
├── server/                       # GoFrame v2 后端 (官方标准分层)
│   ├── api/                      # 接口定义层: 请求/响应结构体 + 路由元信息
│   │   ├── auth/v1               # 登录登出
│   │   └── system/v1             # 用户/角色/菜单/字典
│   ├── internal/
│   │   ├── cmd/                  # 启动入口: 路由注册/静态托管/SPA回退
│   │   ├── controller/           # 控制器层: 薄委托, 校验 service 装配
│   │   ├── service/              # 服务接口层: 定义接口 + 注册函数 (解耦依赖)
│   │   ├── logic/                # 业务实现层: auth/user/role/menu/dict
│   │   ├── dao/                  # 数据访问层 (gf gen dao 生成, 勿手改)
│   │   ├── model/                # do(写入模型) / entity(表实体)
│   │   ├── middleware/           # 中间件: 统一响应/JWT/CORS
│   │   └── packed/               # gf pack 资源包 (单二进制用)
│   ├── resource/sql/init.sql     # RBAC 初始化脚本
│   └── hack/config.yaml          # gf gen dao 配置
├── docs/                         # 文档
├── docker/                       # Dockerfile + docker-compose
├── Makefile                      # 根目录一键脚本
└── README.md
```

**分层调用链**: `cmd(路由) -> controller -> service(接口) -> logic(实现) -> dao -> DB`
- 依赖倒置: logic 实现接口并 `init()` 注册, controller 只面向 service 接口编程
- 新增业务按此链路扩展, 各层职责单一

---

## 二、开发环境准备

```bash
make init-db      # 初始化数据库 (可覆盖 DB_HOST/DB_PORT/DB_PASSWORD)
# 修改 server/manifest/config/config.yaml 的 database.default.link
make dev-server   # 终端1: 后端 :8000
make dev-web      # 终端2: 前端 :8091 (代理 /api -> :8000)
```

> 开发模式下前端 mock 默认关闭（`web/.env.development` 中 `VUE_APP_MOCK_ENABLE=false`）。
> 如需纯前端演示，将其改为 `true` 后重启 `dev-web`。

---

## 三、新增一个业务模块（标准流程）

以"文章管理"为例，演示从建表到接口的全流程。

### 1. 建表

```sql
CREATE TABLE `sys_article` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `title`      varchar(128)  NOT NULL DEFAULT '' COMMENT '标题',
  `content`    text          COMMENT '内容',
  `status`     tinyint(1)    NOT NULL DEFAULT 1  COMMENT '状态(1发布 0下架)',
  `created_at` datetime      COMMENT '创建时间',
  `updated_at` datetime      COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';
```

> 表统一加 `sys_` 前缀, gen dao 时自动去除生成 `Article` 模型。

### 2. 生成 dao / model / entity

```bash
# server/hack/config.yaml 已配置 removePrefix "sys_"
make gen-dao     # 等价于 cd server && make dao
```

### 3. 定义 API（api/system/v1/article.go）

```go
type ArticleListReq struct {
    g.Meta `path:"/article/list" method:"get" tags:"文章管理" summary:"文章列表"`
    Page   int `json:"page" dc:"页码"`
    PageSize int `json:"pageSize" dc:"每页条数"`
}
type ArticleListRes struct {
    List  []*entity.Article `json:"list" dc:"文章列表"`
    Total int               `json:"total" dc:"总数"`
}
```

### 4. 实现 logic + service 接口 + controller + 注册路由

参考现有 `internal/logic/dict/dict.go` 的结构：
- `internal/service/article.go`: 接口定义 + `RegisterArticle`
- `internal/logic/article/article.go`: 实现 + `init()` 注册 + 编译期断言
- `internal/logic/logic.go`: 追加空导入
- `internal/controller/article/`: 薄委托
- `internal/cmd/cmd.go`: 鉴权分组内追加 `article.New()`

### 5. 挂菜单与按钮权限

在"菜单管理"页面配置（或写进 init.sql）：
- 菜单: path=`article`, component=`vab/article/index`（对应前端页面文件）
- 按钮: 权限标识如 `system:article:create`，前端用 `checkPermission` 控制显隐

### 6. 前端页面

新增 `web/src/views/vab/article/index.vue`，接口封装加入 `web/src/api/system.js`，
按钮显隐参考 `web/src/views/vab/system/user/index.vue` 的写法。

---

## 四、权限体系说明

- **RBAC 模型**: 用户 --(sys_user_role)-- 角色 --(sys_role_menu)-- 菜单/按钮
- **动态路由**: 登录后前端调 `GET /api/menu/getRouters`，后端按角色过滤返回树形路由
- **超管角色**: `role_key = admin` 的角色拥有全部路由与按钮权限（consts.RoleKeyAdmin）
- **按钮级权限**: 菜单类型 F + 权限标识（如 `system:user:delete`），
  后端 `GetInfo` 返回 `permissions` 数组，前端 `checkPermission` / `v-permission` 精确匹配
- **内置保护**: id=1 的用户与角色为内置超管，接口层拒绝修改/删除

### JWT 鉴权要点

- 令牌经请求头 `accessToken` 携带（裸 token，后端自动适配 `Bearer` 前缀）
- 载荷键为 `uid`（**不可改为 `id`**：gf-jwt 会把 IdentityKey 写入路由参数，
  与请求体 `id` 字段冲突导致更新错行，详见 `internal/middleware/auth.go` 注释）
- 登出走黑名单机制: logout 后旧令牌立即失效
- 401 业务码触发前端自动登出，403 跳转无权限页

---

## 五、前端定制边界（重要）

为保持上游升级能力，**禁止修改以下内容**，均通过配置与外围扩展实现：

| 需求 | 正确做法 |
| --- | --- |
| 改接口地址 / 错误码 | `web/src/config/net.config.js` |
| 改路由模式 / 站点标题 | `web/src/config/setting.config.js` |
| 对接后端接口 | `web/src/api/*.js`（只改 url 与 method） |
| 新增页面 | `web/src/views/vab/**` 新增目录 |
| 开关 mock | `.env.development` / `.env.production` |
| 构建代理 | `web/rspack.config.js` 的 devServer.proxy |

---

## 六、常用命令速查

```bash
make help           # 查看全部可用命令
make dev-server     # 启动后端
make dev-web        # 启动前端
make init-db        # 重置数据库(危险, 会清空重灌)
make gen-dao        # 按表结构重新生成 dao/model
make gen-password PASSWORD=xxx   # 生成 bcrypt 密码哈希
make build          # 全量打包单二进制
make docker-up      # Docker Compose 启动
make clean          # 清理构建产物
```
