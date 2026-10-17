# 部署指南

本项目支持三种部署方式，按需选择：

| 方式 | 适用场景 | 说明 |
| --- | --- | --- |
| 源码部署 | 开发 / 联调 | 前端 devServer + 后端 go run，接口代理转发 |
| 单二进制部署 | 生产推荐 | 前端资源内嵌，一个可执行文件 + 一份配置文件 |
| Docker Compose | 容器化环境 | 应用 + MySQL 一键编排 |

> 三种方式的默认管理员账号均为 `admin / 123456`，登录后请立即修改密码。

---

## 一、源码部署（开发模式）

### 环境要求

- Go >= 1.26
- Node.js >= 22（需支持 `require(esm)`，rspack 构建依赖）
- pnpm >= 8
- MySQL >= 5.7（推荐 8.0）

### 步骤

```bash
# 1. 初始化数据库 (建库 + RBAC 表 + 种子数据)
#    参数可按环境覆盖: make init-db DB_HOST=... DB_PORT=... DB_USER=... DB_PASSWORD=... DB_NAME=...
make init-db

# 2. 修改后端配置 (数据库连接)
#    编辑 server/manifest/config/config.yaml 中 database.default.link
#    link: "mysql:root:你的密码@tcp(127.0.0.1:3306)/gf_vue_admin_better?charset=utf8mb4&parseTime=true"

# 3. 启动后端 (默认监听 :8000)
make dev-server

# 4. 启动前端 (默认监听 :8091, /api 自动代理到 :8000)
make dev-web
```

浏览器访问 `http://127.0.0.1:8091` 开始开发。

---

## 二、单二进制部署（生产推荐）

将前端构建产物经 `gf pack` 内嵌进 GoFrame 二进制，产出**单文件可执行程序**，
无需 Nginx、无需分发静态目录，天然支持 SPA 路由（非 /api 请求自动回退 index.html）。

### 打包

```bash
# 一键完成: 前端构建 -> 资源内嵌 -> 单二进制产出
# 产物: bin/gf-vue-admin-better
make build
```

### 部署

```bash
# 服务器目录结构
/opt/gfab/
├── gf-vue-admin-better          # 单二进制 (含前端资源)
└── manifest/config/config.yaml  # 配置文件 (数据库连接等)

# 启动 (GF 从工作目录读取 manifest/config/config.yaml)
cd /opt/gfab && ./gf-vue-admin-better
```

生产环境配置建议：

```yaml
# manifest/config/config.yaml
database:
  logger:
    level: "error"          # 关闭 SQL 全量日志
  default:
    link: "mysql:root:密码@tcp(数据库地址:3306)/gf_vue_admin_better?charset=utf8mb4&parseTime=true"
cors:
  enabled: false            # 同源部署无需 CORS, 前后端分离独立域名时再开启
```

也可用 systemd 托管进程：

```ini
# /etc/systemd/system/gfab.service
[Unit]
Description=gf-vue-admin-better
After=network.target mysql.service

[Service]
WorkingDirectory=/opt/gfab
ExecStart=/opt/gfab/gf-vue-admin-better
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable --now gfab
```

---

## 三、Docker Compose 一键部署

内置应用 + MySQL 编排，首次启动自动执行 `server/resource/sql/init.sql` 完成数据库初始化。

```bash
# 构建镜像并启动 (构建上下文为项目根目录)
make docker-up
# 等价于: docker compose -f docker/docker-compose.yml up -d --build
```

启动完成后访问 `http://服务器IP:8000`。

自定义数据库密码：编辑 `docker/docker-compose.yml` 中 `MYSQL_ROOT_PASSWORD` 与
`docker/config.yaml` 的 `link`（两者保持一致）后重新 `docker-up`。
> 注意：GF 不支持通过环境变量覆盖数据库连接，容器内配置统一走挂载的 `docker/config.yaml`。

常用操作：

```bash
make docker-down      # 停止并移除容器
make docker-build     # 仅构建镜像
```

---

## 四、Nginx 反向代理（可选）

若需独立域名部署前端与后端（不使用单二进制托管），参考：

```nginx
server {
    listen 80;
    server_name admin.example.com;

    # 前端静态资源 (web/dist)
    location / {
        root /opt/gfab/web/dist;
        try_files $uri $uri/ /index.html;   # SPA 回退
    }

    # 后端 API 反向代理
    location /api/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

此时后端配置 `cors.enabled: true` 放行跨域。

---

## 常见问题

**Q: 页面刷新 404？**
单二进制模式已内置 SPA 回退，无需处理；Nginx 模式需配置 `try_files`（见上文）。

**Q: 忘记 admin 密码？**
重新执行初始化 SQL（会重置权限数据），或生成新哈希后更新库：
```bash
make gen-password PASSWORD=新密码     # 生成 bcrypt 哈希
mysql> UPDATE sys_user SET password='生成的哈希' WHERE id=1;
```

**Q: 修改前端后如何更新单二进制？**
重新执行 `make build` 即可（全流程自动）。
