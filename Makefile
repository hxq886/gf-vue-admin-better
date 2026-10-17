# ==============================================================================
# gf-vue-admin-better 一键工程脚本
# 用法: make <target>  (如 make build / make dev-server / make init-db)
# ==============================================================================

# 数据库连接参数 (可通过环境变量覆盖, 与 server/manifest/config/config.yaml 保持一致)
DB_HOST     ?= 127.0.0.1
DB_PORT     ?= 3306
DB_USER     ?= root
DB_PASSWORD ?=
DB_NAME     ?= gf_vue_admin_better

GO_PROXY    ?= https://goproxy.cn,direct
NPM_REGISTRY?= https://registry.npmmirror.com

.PHONY: help
help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ==================== 开发模式 ====================

.PHONY: dev-server
dev-server: ## 启动后端 (GoFrame, 默认 :8000, 需先初始化数据库)
	cd server && GOPROXY=$(GO_PROXY) go run main.go

.PHONY: dev-web
dev-web: ## 启动前端 (rspack devServer, 默认 :8091, 代理 /api 到后端)
	cd web && pnpm install --registry=$(NPM_REGISTRY) && pnpm serve

.PHONY: dev
dev: ## 并行启动前后端 (开发模式)
	@$(MAKE) -j2 dev-server dev-web

# ==================== 数据库 ====================

.PHONY: init-db
init-db: ## 初始化数据库 (建库 + RBAC 表 + 种子数据, 内置管理员 admin/123456)
	cd server && $(MAKE) init-db DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME)

.PHONY: gen-dao
gen-dao: ## 按数据库重新生成 dao/model/entity (修改表结构后执行)
	cd server && $(MAKE) dao

# ==================== 构建与打包 ====================

.PHONY: build-web
build-web: ## 构建前端产物 (web/dist)
	cd web && pnpm install --registry=$(NPM_REGISTRY) && pnpm build

.PHONY: build
build: ## 全量打包单二进制 (前端构建 -> gf pack 内嵌 -> go build 产出 bin/gf-vue-admin-better)
	$(MAKE) build-web
	rm -rf server/resource/public/dist
	mkdir -p server/resource/public/dist
	cp -r web/dist/. server/resource/public/dist/
	cd server && GOPROXY=$(GO_PROXY) $(MAKE) pack && GOPROXY=$(GO_PROXY) CGO_ENABLED=0 go build -o ../bin/gf-vue-admin-better .
	@echo ">>> 单二进制打包完成: bin/gf-vue-admin-better (含前端静态资源, 直接运行即可)"

.PHONY: build-server
build-server: ## 仅构建后端二进制 (不含前端资源)
	cd server && GOPROXY=$(GO_PROXY) CGO_ENABLED=0 go build -o ../bin/gf-vue-admin-better-server .

.PHONY: gen-password
gen-password: ## 生成 bcrypt 密码哈希 (用法: make gen-password PASSWORD=你的明文密码)
	cd server && $(MAKE) gen-password PASSWORD=$(PASSWORD)

# ==================== Docker ====================

.PHONY: docker-build
docker-build: ## 构建全栈 Docker 镜像 (前端构建 + 单二进制)
	docker build -t gf-vue-admin-better -f docker/Dockerfile .

.PHONY: docker-up
docker-up: ## Docker Compose 一键启动 (应用 + MySQL)
	docker compose -f docker/docker-compose.yml up -d

.PHONY: docker-down
docker-down: ## 停止并移除 Docker Compose 服务
	docker compose -f docker/docker-compose.yml down

# ==================== 清理 ====================

.PHONY: clean
clean: ## 清理构建产物
	rm -rf bin server/resource/public/dist
	cd web && rm -rf node_modules dist
	cd server && rm -rf temp

# 默认目标
.DEFAULT_GOAL := help
