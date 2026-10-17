# 编译部署详细教程

> 本文档覆盖从零开始的环境准备、源码获取、数据库初始化、开发运行，到生产编译（单二进制）、服务器部署、Docker 部署、交叉编译与常见问题排查。
>
> 全文命令均已在 Linux (Ubuntu/CentOS/Debian) 实测通过；MacOS 命令基本一致，Windows 建议使用 WSL2。

---

## 目录

1. [环境要求总览](#一环境要求总览)
2. [环境安装](#二环境安装)
3. [获取源码](#三获取源码)
4. [数据库准备](#四数据库准备)
5. [开发模式运行](#五开发模式运行)
6. [生产编译（单二进制打包）](#六生产编译单二进制打包)
7. [服务器部署单二进制](#七服务器部署单二进制)
8. [Docker 部署](#八docker-部署)
9. [交叉编译](#九交叉编译)
10. [版本升级与重新打包](#十版本升级与重新打包)
11. [常见问题排查](#十一常见问题排查)

---

## 一、环境要求总览

| 工具 | 最低版本 | 用途 | 说明 |
| --- | --- | --- | --- |
| Go | **1.26+** | 编译后端 | `server/go.mod` 指定 go 1.26 |
| Node.js | **22.12+** | 构建前端 | 需支持 `require(esm)`（rspack 2.x 为纯 ESM 包，Node 20 及以下会报 `ERR_REQUIRE_ESM`） |
| pnpm | 8+ | 前端包管理 | `corepack enable` 或 npm 全局安装 |
| MySQL | 5.7+（推荐 8.0） | 数据存储 | 字符集需 utf8mb4 |
| gf CLI | v2.10.3 | 资源打包/代码生成 | 仅打包时需要，`make pack` 会自动安装 |
| GNU Make | 3.81+ | 执行脚本 | Linux/macOS 一般自带 |
| Docker | 20+（含 compose 插件） | 容器部署 | 仅 Docker 方案需要 |

查看本机版本：

```bash
go version          # go version go1.26.x linux/amd64
node -v             # v22.x.x
pnpm -v             # 8.x 或更高
mysql --version     # mysql  Ver 8.0.x / 5.7.x
gf -v               # v2.10.x
make -v             # GNU Make 4.x
docker --version    # Docker version 24.x
```

---

## 二、环境安装

### 2.1 安装 Go（>= 1.26）

**Linux（官方二进制）**：

```bash
# 下载 (以 1.26.0 为例, 请到 https://go.dev/dl/ 确认最新版本号)
wget https://golang.google.cn/dl/go1.26.0.linux-amd64.tar.gz

# 解压到 /usr/local (覆盖旧版本先 rm -rf /usr/local/go)
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz

# 配置环境变量 (写入 ~/.bashrc 或 ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

**国内代理（必配，否则 go mod 下载会超时）**：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

**MacOS**：`brew install go`

**验证**：`go version` 输出 >= 1.26 即可。

### 2.2 安装 Node.js（>= 22.12）

版本要求说明：本项目前端构建工具 rspack 2.x 是纯 ESM 包，`rspack.js` 以 CommonJS `require()` 方式加载它。Node 22.12+ 原生支持 `require(esm)`；Node 20 及以下会报错 `ERR_REQUIRE_ESM`。**请直接安装 Node 22 LTS。**

**Linux（官方二进制）**：

```bash
# 下载 Node 22 (以 22.14.0 为例, 最新版见 https://nodejs.org/)
wget https://cdn.npmmirror.com/binaries/node/v22.14.0/node-v22.14.0-linux-x64.tar.xz

# 解压到 /usr/local
sudo rm -rf /usr/local/{bin/node,bin/npm,bin/npx,lib/node_modules}
sudo tar -C /usr/local --strip-components=1 -xf node-v22.14.0-linux-x64.tar.xz
rm node-v22.14.0-linux-x64.tar.xz

node -v && npm -v
```

**MacOS**：`brew install node@22`

**验证**：

```bash
# 确认 require(esm) 可用
node -e "require('node:path'); console.log(process.version)"
# 版本 >= v22.12.0 即可
```

### 2.3 安装 pnpm

```bash
# 方式一: npm 全局安装 (推荐, 已配置国内镜像)
npm i -g pnpm --registry=https://registry.npmmirror.com

# 方式二: corepack (Node 自带)
corepack enable

# 配置国内源 (加速依赖下载)
pnpm config set registry https://registry.npmmirror.com

pnpm -v
```

### 2.4 安装 MySQL（5.7+ / 8.0）

已有 MySQL 可跳过。以下为 Docker 快速起一个（仅开发用）：

```bash
docker run -d --name gfab-mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=gf_vue_admin_better \
  -p 3306:3306 \
  mysql:8.0
```

裸机安装请参考各发行版文档（Ubuntu: `apt install mysql-server`；CentOS: 官方 yum 源），关键是确保：

- 字符集为 `utf8mb4`
- 拥有一个可登录的账号（如 root，或专门创建的业务账号）

```sql
-- 创建专用账号 (可选, 生产推荐)
CREATE DATABASE gf_vue_admin_better DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER 'gfab'@'%' IDENTIFIED BY '你的强密码';
GRANT ALL PRIVILEGES ON gf_vue_admin_better.* TO 'gfab'@'%';
FLUSH PRIVILEGES;
```

### 2.5 安装 gf CLI（可选，打包时自动装）

gf CLI 用于 `gf pack`（资源内嵌）与 `gf gen dao`（代码生成）：

```bash
# 方式一: go install (需 Go 1.26+)
go install github.com/gogf/gf/cmd/gf/v2@latest

# 方式二: 脚本安装 (Linux/MacOS)
wget https://goframe.org/cli.sh && chmod +x cli.sh && ./cli.sh

gf -v
```

> 提示：项目根 Makefile 的 `make pack` 依赖 `server/hack/hack-cli.mk` 的 `cli.install` 目标，检测不到 gf 命令时会**自动安装**，因此也可跳过本步。

---

## 三、获取源码

```bash
# Gitee (国内快)
git clone https://gitee.com/hxq886123/gf-vue-admin-better.git

#Gitcode (国内快)
git https://gitcode.com/hxq886/gf-vue-admin-better.git

# GitHub
git clone https://github.com/hxq886/gf-vue-admin-better.git

cd gf-vue-admin-better
```

目录结构：

```
gf-vue-admin-better/
├── web/          # 前端源码 (Vue3-Admin-Better)
├── server/       # 后端源码 (GoFrame v2)
├── docs/         # 文档 (部署/开发/API)
├── docker/       # Dockerfile / docker-compose / 容器专用配置
├── Makefile      # 根一键脚本 (make help 查看全部命令)
└── LICENSE       # MIT
```

---

## 四、数据库准备

### 方式一：make 一键初始化（推荐）

```bash
# 默认连接 127.0.0.1:3306, root 空密码
make init-db

# 自定义连接参数 (按需覆盖)
make init-db DB_HOST=127.0.0.1 DB_PORT=3306 DB_USER=root DB_PASSWORD=yourpass
```

该命令做两件事：

1. 建库：`CREATE DATABASE IF NOT EXISTS gf_vue_admin_better ...`
2. 导入 [server/resource/sql/init.sql](../server/resource/sql/init.sql)：7 张 RBAC 表（用户/角色/菜单/角色菜单/用户角色/字典类型/字典数据）+ 种子数据 + 内置管理员 `admin/123456`

### 方式二：手动执行

```bash
mysql -h127.0.0.1 -P3306 -uroot -p -e \
  "CREATE DATABASE gf_vue_admin_better DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

mysql -h127.0.0.1 -P3306 -uroot -p gf_vue_admin_better < server/resource/sql/init.sql
```

### 验证初始化结果

```bash
mysql -h127.0.0.1 -uroot -p -e "USE gf_vue_admin_better; SHOW TABLES; SELECT id,username,nickname FROM sys_user;"
```

预期输出包含 7 张 `sys_*` 表，以及一行 `id=1, username=admin` 记录。

---

## 五、开发模式运行

开发模式下前后端独立进程：后端 GoFrame 直接 `go run`（改动即重启需手动），前端 rspack devServer 热更新并代理 `/api` 到后端。

### 5.1 修改后端数据库连接

编辑 [server/manifest/config/config.yaml](../server/manifest/config/config.yaml)：

```yaml
database:
  logger:
    level: "all"            # 开发期可改 "all" 观察 SQL
  default:
    link: "mysql:root:yourpass@tcp(127.0.0.1:3306)/gf_vue_admin_better?charset=utf8mb4&parseTime=true"
```

> `parseTime=true` 必须保留，否则时间字段扫描报错。

### 5.2 启动后端（终端 1）

```bash
make dev-server        # 等价: cd server && go run main.go
```

首次运行会自动下载 Go 依赖（已配 goproxy.cn，无需科学上网）。看到以下输出即成功：

```
SERVER ADDRESS  :8000
ROUTE TABLE: ... (POST /api/login, GET /api/user/getInfo ...)
```

### 5.3 启动前端（终端 2）

```bash
make dev-web           # 等价: cd web && pnpm install && pnpm serve
```

> 也可以一条命令并行启动两个服务：`make dev`

首次启动安装依赖（约 1-2 分钟）。完成后访问：

```
http://127.0.0.1:8091
```

> 跨域说明：开发态前端 8091、后端 8000，[web/rspack.config.js](../web/rspack.config.js) 中 devServer 已配置 `/api` 代理到 `http://127.0.0.1:8000`，无需额外处理跨域。

### 5.4 登录验证

浏览器打开 `http://127.0.0.1:8091`，使用内置账号登录：

- 用户名：`admin`
- 密码：`123456`

登录后左侧应出现"系统管理"菜单（用户/角色/菜单/字典管理四个子页面均可正常访问）。

> mock 说明：前端请求默认直连真实后端（`.env.development` 中 `VUE_APP_MOCK_ENABLE=false`）。如需纯前端演示（无后端），改为 `true` 重启前端即可。

---

## 六、生产编译（单二进制打包）

### 6.1 打包原理

```
web (Vue3 源码)
  │  pnpm build                  ① 前端构建 → web/dist/
  ▼
server/resource/public/dist      ② 复制静态产物到后端资源目录
  │  gf pack                     ③ 资源打包为 Go 源码 → internal/packed/packed.go (gres 压缩内嵌)
  ▼
go build                         ④ 编译 → 单二进制 bin/gf-vue-admin-better (~31MB)
```

产物特性：**一个可执行文件同时提供 API 服务与前端页面**，非 `/api` 请求自动回退 `index.html`（SPA 刷新不 404），无需 Nginx。

### 6.2 一键打包（推荐）

在**项目根目录**执行：

```bash
make build
```

该命令依次完成：前端构建 → 产物复制 → gf pack 内嵌 → 静态编译。完成后输出：

```
>>> 单二进制打包完成: bin/gf-vue-admin-better (含前端静态资源, 直接运行即可)
```

产物约 31MB。**注意：打包会覆盖 `server/internal/packed/packed.go`（占位文件 → 2MB+ 资源文件），该文件已在 .gitignore 之外单独说明，提交代码前可执行 `git checkout -- server/internal/packed/packed.go` 还原为占位版本。**

### 6.3 分步执行（理解每一步）

```bash
# ① 构建前端 (在 web/ 目录)
cd web
pnpm install --registry=https://registry.npmmirror.com
pnpm build                       # 产物: web/dist/

# ② 复制产物到后端资源目录
cd ..
rm -rf server/resource/public/dist
mkdir -p server/resource/public/dist
cp -r web/dist/. server/resource/public/dist/

# ③ gf pack 内嵌 (在 server/ 目录; -p 前缀保证资源键名与代码读取路径对齐, 不可省略)
cd server
make pack
# 等价: gf pack resource/public/dist internal/packed/packed.go -n packed -p=resource/public/dist

# ④ 静态编译 (无 libc 依赖, 任意 Linux 可运行; -s -w 去符号表减小体积)
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o ../bin/gf-vue-admin-better .
```

### 6.4 本地快速验证产物

```bash
# 启动 (会读取工作目录下 manifest/config/config.yaml)
cd server && ../bin/gf-vue-admin-better
```

```bash
# 验证 1: 首页 (应返回 index.html)
curl -s http://127.0.0.1:8000/ | head -c 100

# 验证 2: SPA 刷新回退 (应同样返回 index.html, 而非 404)
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:8000/system/user   # 200

# 验证 3: 静态资源
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:8000/js/app.*.js  # 200

# 验证 4: API 登录 (应返回 code:200 + accessToken)
curl -s -X POST http://127.0.0.1:8000/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}'

# 验证 5: 未知 API 返回 JSON 404 (而非 HTML)
curl -s http://127.0.0.1:8000/api/notexist   # {"code":404,...}
```

浏览器直接访问 `http://127.0.0.1:8000` 应看到登录页。

---

## 七、服务器部署单二进制

### 7.1 目标服务器要求

- Linux x86_64（任意发行版：CentOS/Ubuntu/Debian/麒麟...）
- **无需安装 Go / Node / Nginx**（二进制自包含前端）
- 需要 MySQL（本机或远程均可）
- 放行 8000 端口（或在负载均衡/安全组中配置）

### 7.2 上传文件

服务器只需两个文件（二进制 + 配置）：

```bash
# 本地打包完成后上传
scp bin/gf-vue-admin-better root@服务器IP:/opt/gfab/
scp server/manifest/config/config.yaml root@服务器IP:/opt/gfab/manifest/config/config.yaml
```

目标结构：

```
/opt/gfab/
├── gf-vue-admin-better              # 单二进制 (含前端)
└── manifest/
    └── config/
        └── config.yaml              # 生产配置
```

> GF 约定：程序从**工作目录**下 `manifest/config/config.yaml` 读取配置，因此必须以 `/opt/gfab` 为工作目录启动。

### 7.3 修改生产配置

服务器上编辑 `/opt/gfab/manifest/config/config.yaml`：

```yaml
server:
  address:     ":8000"        # 监听端口, 改端口只需改这里

database:
  logger:
    level: "error"            # 生产必改: 关闭 SQL 全量日志
  default:
    link: "mysql:gfab:强密码@tcp(数据库地址:3306)/gf_vue_admin_better?charset=utf8mb4&parseTime=true"

cors:
  enabled: false              # 同源部署 (二进制同时托管前端+API) 无需跨域

jwt:
  key: "请替换为32位以上随机字符串"   # 生产必改!
  timeout: "24h"
```

> 重要：GF **不支持**通过环境变量覆盖 `database.default.link`（实测 `GF_DATABASE_DEFAULT_LINK` 无效），数据库连接必须写入配置文件。Docker 场景通过挂载配置文件实现（见第八节）。

### 7.4 初始化服务器数据库

```bash
# 在服务器上 (或任何能连到数据库的机器)
mysql -h数据库地址 -uroot -p -e \
  "CREATE DATABASE gf_vue_admin_better DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
mysql -h数据库地址 -uroot -p gf_vue_admin_better < server/resource/sql/init.sql
```

### 7.5 启动与验证

```bash
cd /opt/gfab
chmod +x gf-vue-admin-better
./gf-vue-admin-better
```

看到路由表输出后，浏览器访问 `http://服务器IP:8000`，admin/123456 登录。验证 SPA 刷新：进入"系统管理 → 用户管理"后按 F5，页面应正常（不 404）。

### 7.6 systemd 常驻托管

前台方式退出 SSH 即停止，生产必须用 systemd：

```bash
sudo tee /etc/systemd/system/gfab.service > /dev/null <<'EOF'
[Unit]
Description=gf-vue-admin-better
After=network.target mysql.service

[Service]
Type=simple
WorkingDirectory=/opt/gfab
ExecStart=/opt/gfab/gf-vue-admin-better
Restart=always
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now gfab     # 开机自启 + 立即启动
sudo systemctl status gfab           # 查看状态 (active (running))
```

常用运维命令：

```bash
sudo systemctl restart gfab          # 重启
sudo systemctl stop gfab             # 停止
journalctl -u gfab -f                # 实时日志
journalctl -u gfab --since today     # 今日日志
```

### 7.7 防火墙放行（按需）

```bash
# CentOS firewalld
sudo firewall-cmd --permanent --add-port=8000/tcp && sudo firewall-cmd --reload

# Ubuntu ufw
sudo ufw allow 8000/tcp

# 云服务器还需在控制台安全组放行 8000
```

### 7.8 修改默认密码（必做）

```bash
# 生成新密码哈希 (在打包机上执行)
make gen-password PASSWORD=你的新密码
# 输出: $2a$10$xxxxxxxx...

# 更新数据库
mysql -h数据库地址 -uroot -p -e \
  "UPDATE gf_vue_admin_better.sys_user SET password='上面生成的哈希' WHERE id=1;"
```

---

## 八、Docker 部署

### 8.1 文件说明

```
docker/
├── Dockerfile           # 多阶段构建: node22 前端构建 → golang1.26 后端构建 → alpine 运行时
├── docker-compose.yml   # 应用 + MySQL 编排 (首次启动自动执行 init.sql)
└── config.yaml          # 容器专用配置 (数据库指向 compose 内的 mysql 服务)
```

### 8.2 一键启动

```bash
# 项目根目录执行
make docker-up
# 等价: docker compose -f docker/docker-compose.yml up -d --build
```

流程：

1. 阶段一 `node:22-alpine`：pnpm 安装依赖 + `pnpm build` 产出前端
2. 阶段二 `golang:1.26-alpine`：`gf pack` 内嵌前端 → `go build` 单二进制
3. 阶段三 `alpine:3.20`：仅携带二进制 + 配置，时区设为 Asia/Shanghai
4. MySQL 容器首次启动时自动执行 `server/resource/sql/init.sql` 完成建表与种子数据
5. 应用容器等待 MySQL 健康检查通过后启动，配置挂载 `docker/config.yaml`（连接 `mysql:3306`）

访问 `http://服务器IP:8000`，admin/123456 登录。

### 8.3 自定义数据库密码

修改两处保持一致：

```yaml
# docker/docker-compose.yml
MYSQL_ROOT_PASSWORD: your_strong_pass

# docker/config.yaml
link: "mysql:root:your_strong_pass@tcp(mysql:3306)/gf_vue_admin_better?..."
```

重新启动：`make docker-down && make docker-up`（注意：已初始化过的数据卷需 `docker volume rm` 后才会重跑 init.sql）。

### 8.4 常用操作

```bash
make docker-down                 # 停止并移除容器 (数据卷保留)
docker compose -f docker/docker-compose.yml logs -f app    # 看应用日志
docker compose -f docker/docker-compose.yml exec mysql \
  mysql -uroot -p123456 gf_vue_admin_better                 # 进数据库
make docker-build                # 仅构建镜像不启动
```

### 8.5 仅部署应用到已有数据库

不需要 compose 里的 MySQL 时，注释掉 `mysql` 服务与 `volumes`，并把 `docker/config.yaml` 的 link 改为外部数据库地址，然后 `make docker-up`。

---

## 九、交叉编译

在打包机（如开发用的 Mac/Windows）上编译出目标平台二进制：

```bash
# 前端与 gf pack 步骤不变 (先完成 6.3 的 ①②③ 步)
cd server

# Linux x86_64 (常见云服务器)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ../bin/gf-vue-admin-better-linux-amd64 .

# Linux ARM64 (如树莓派/国产 ARM 服务器)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o ../bin/gf-vue-admin-better-linux-arm64 .

# Mac Apple Silicon
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../bin/gf-vue-admin-better-darwin-arm64 .

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../bin/gf-vue-admin-better-windows-amd64.exe .
```

> 静态资源已通过 gf pack 内嵌，交叉编译产物同样"开箱即用"（Windows 下配置路径为 `manifest\config\config.yaml`）。

---

## 十、版本升级与重新打包

### 前端有改动

```bash
make build          # 全流程重跑, 生成新二进制
# 服务器替换二进制后
sudo systemctl restart gfab
```

### 仅后端有改动

```bash
# 跳过前端构建, 直接复用已存在的 resource/public/dist
cd server
make pack && CGO_ENABLED=0 go build -o ../bin/gf-vue-admin-better .
```

### 数据库结构有改动

```bash
# 1. 修改表结构 SQL
# 2. 重新生成 dao/model/entity
make gen-dao        # 等价: cd server && make dao (依赖 server/hack/config.yaml 的 link 配置)
# 3. 重新编译部署
make build
```

### 还原 packed.go 为占位（提交代码前）

```bash
git checkout -- server/internal/packed/packed.go
# 或手动写入占位内容:
# // Package packed 内嵌静态资源包占位文件。
# package packed
```

---

## 十一、常见问题排查

### Q1: 前端构建报 `ERR_REQUIRE_ESM` / `require() of ES Module`

```
Error [ERR_REQUIRE_ESM]: require() of ES Module .../@rspack/core/... not supported
```

**原因**：Node 版本低于 22.12，不支持 `require()` 加载 ESM 包。

**解决**：升级 Node 至 22 LTS：`node -v` 确认 >= 22.12.0。

### Q2: `go mod download` 卡住/超时

**解决**：配置代理 `go env -w GOPROXY=https://goproxy.cn,direct` 后重试。本项目依赖均可在国内镜像拉取。

### Q3: 启动报 `configuration not found` 或连接信息为空

**原因**：工作目录不对，GF 找不到 `manifest/config/config.yaml`。

**解决**：必须在 `/opt/gfab`（即 `manifest` 的上级目录）启动：

```bash
cd /opt/gfab && ./gf-vue-admin-better
```

systemd 场景检查 `WorkingDirectory=/opt/gfab` 配置。

### Q4: 登录返回 `Error 1045 Access denied` / `Error 1049 Unknown database`

**原因**：数据库账号密码错误 / 库未创建。

**解决**：
- 1045：核对 `config.yaml` 的 link 中账号密码，注意格式 `mysql:用户:密码@tcp(host:port)/库名`
- 1049：执行第四节初始化数据库

### Q5: 登录返回 `Error 1054 Unknown column` 或表不存在

**原因**：init.sql 未导入或导入了旧版本。

**解决**：重新导入 `server/resource/sql/init.sql`（注意会重置权限数据）。

### Q6: 页面刷新 404

**原因**：走了 Nginx/其他静态服务器托管前端，未配置 SPA 回退。

**解决**：单二进制模式自带回退无需处理；Nginx 场景增加：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

### Q7: 登录后立即被踢回登录页（401 死循环）

**排查顺序**：
1. F12 → Network → `POST /api/login` 响应是否 `code:200`？非 200 见 Q4/Q5
2. `GET /api/user/getInfo` 请求头是否携带 `accessToken`？未携带则检查 `web/src/config/net.config.js` 的 `tokenName: "accessToken"`
3. 确认前端构建时 `.env.production` 中 `VUE_APP_MOCK_ENABLE=false`（生产禁用 mock，否则登录接口被 mock 截胡返回假 token）

### Q8: 端口被占用 `address already in use`

```bash
# 查占用
sudo lsof -i:8000    # 或 ss -lntp | grep 8000
# 改端口: 修改 config.yaml 的 server.address 为 ":8001" 后重启
```

### Q9: `gf: command not found`（手动分步打包时）

```bash
make pack    # 内置 cli.install 会自动安装 gf CLI
# 或手动: go install github.com/gogf/gf/cmd/gf/v2@latest
```

### Q10: Docker 构建前端阶段失败 `pnpm install` 超时

**解决**：Dockerfile 已配置 npmmirror 镜像源；仍慢可检查 Docker 网络或给 daemon 配置代理。构建缓存优化：未改 `package.json` 时重复构建会命中缓存层，秒级完成依赖安装。

### Q11: 服务器时区不对（时间差 8 小时）

**解决**：单二进制模式设置环境变量 `TZ=Asia/Shanghai`（systemd 单元加 `Environment=TZ=Asia/Shanghai`）；Docker 模式镜像已内置上海时区。

### Q12: 如何确认二进制里真的内嵌了前端？

```bash
# 体积 ~31MB 且启动后访问 / 返回 index.html 即已内嵌
# 若只有 ~29MB 且访问 / 返回 404/空白, 说明 gf pack 步骤被跳过
ls -la bin/gf-vue-admin-better
curl -s http://127.0.0.1:8000/ | head -c 50
```

---

> 更多文档：[部署方案对比](deploy.md) | [二次开发指南](development.md) | [接口文档](api.md)
