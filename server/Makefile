ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "gf-vue-admin-better"
DOCKER_NAME = "gf-vue-admin-better"

# 数据库连接参数 (可通过环境变量覆盖, 默认与 manifest/config/config.yaml 保持一致)
DB_HOST     ?= 127.0.0.1
DB_PORT     ?= 3306
DB_USER     ?= root
DB_PASSWORD ?=
DB_NAME     ?= gf_vue_admin_better
MYSQL       ?= mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $(if $(DB_PASSWORD),-p$(DB_PASSWORD),)

include ./hack/hack-cli.mk
include ./hack/hack.mk

# ==================== 项目专属目标 ====================

# 本地开发运行 (热重启依赖 gf run, 也可直接 go run)
.PHONY: run
run: cli.install
	@gf run main.go

# 初始化数据库: 建库 + 执行 RBAC 初始化脚本 (含内置管理员账号 admin/123456)
.PHONY: init-db
init-db:
	@$(MYSQL) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
	@$(MYSQL) $(DB_NAME) < $(ROOT_DIR)/resource/sql/init.sql
	@echo ">>> 数据库初始化完成, 内置管理员账号: admin / 123456"

# 生成 bcrypt 密码哈希 (用法: make gen-password PASSWORD=你的明文密码)
.PHONY: gen-password
gen-password:
	@go run ./internal/genhash $(PASSWORD)

# 打包前端静态资源进二进制资源包 (需先在前端目录执行构建, dist 已复制到 resource/public)
# 注意: -p 前缀使内嵌键名为 resource/public/dist/..., 与 cmd.go 的 distDir 对齐
.PHONY: pack
pack: cli.install
	@rm -f $(ROOT_DIR)/internal/packed/packed.go
	@gf pack $(ROOT_DIR)/resource/public/dist $(ROOT_DIR)/internal/packed/packed.go -n packed -p=resource/public/dist

# 清理构建产物
.PHONY: clean
clean:
	@if [ -d $(ROOT_DIR)/temp ]; then rm -rf $(ROOT_DIR)/temp; fi
	@gofmt -l -w $(ROOT_DIR)/main.go $(ROOT_DIR)/internal
