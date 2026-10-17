# ==============================================================================
# gf-vue-admin-better 全栈镜像构建 (多阶段: 前端构建 -> 后端构建 -> 运行)
# 构建上下文为项目根目录: docker build -t gf-vue-admin-better -f docker/Dockerfile .
# ==============================================================================

# ============ 阶段一: 构建前端 ============
FROM node:22-alpine AS web-builder
WORKDIR /build/web

# 先装依赖以利用层缓存
COPY web/package.json web/pnpm-lock.yaml ./
RUN corepack enable && \
    pnpm config set registry https://registry.npmmirror.com && \
    pnpm install --no-frozen-lockfile

# 再拷源码构建 (产出 dist)
COPY web/ .
RUN pnpm build

# ============ 阶段二: 构建后端单二进制 ============
FROM golang:1.26-alpine AS server-builder
WORKDIR /build/server

ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# 安装 gf CLI (用于 gf pack 内嵌前端资源)
RUN go install github.com/gogf/gf/cmd/gf/v2@v2.10.3

# 先下载依赖以利用层缓存
COPY server/go.mod server/go.sum ./
RUN go mod download

# 拷贝后端源码, 注入前端构建产物并打包
COPY server/ .
COPY --from=web-builder /build/web/dist ./resource/public/dist
RUN gf pack resource/public/dist internal/packed/packed.go -n packed -p=resource/public/dist && \
    go build -trimpath -ldflags="-s -w" -o /out/gf-vue-admin-better .

# ============ 阶段三: 运行时 ============
FROM alpine:3.20
WORKDIR /app

# 时区与证书
RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 单二进制 (已内嵌前端资源) + 外置配置文件 (数据库连接等按环境修改)
COPY --from=server-builder /out/gf-vue-admin-better ./
COPY server/manifest/config/config.yaml ./manifest/config/

EXPOSE 8000
ENTRYPOINT ["./gf-vue-admin-better"]
