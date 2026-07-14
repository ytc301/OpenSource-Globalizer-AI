.PHONY: build run test lint clean deps docker-build docker-up docker-down

# 项目变量
BINARY_NAME=globalizer
BUILD_DIR=bin
MAIN_PATH=./cmd/globalizer
GO=go
GOFLAGS=-ldflags="-s -w"

# 默认目标
all: deps build

# 安装依赖
deps:
	$(GO) mod tidy
	$(GO) mod download

# 编译
build:
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# 交叉编译 (Linux AMD64, 用于 GitHub Action)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)

# 运行
run:
	$(GO) run ./$(MAIN_PATH) $(ARGS)

# 测试
test:
	$(GO) test -v -race -cover ./...

# 测试覆盖率报告
test-cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	golangci-lint run ./...

# 代码格式化
fmt:
	$(GO) fmt ./...

# 清理
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Docker 构建
docker-build:
	docker build -t opensource-globalizer:latest .

# Docker 启动
docker-up:
	docker compose up -d

# Docker 停止
docker-down:
	docker compose down

# 安装到本地
install:
	$(GO) install ./$(MAIN_PATH)

# 帮助
help:
	@echo "OpenSource Globalizer AI — Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make deps         安装依赖"
	@echo "  make build        编译"
	@echo "  make build-linux  交叉编译 (Linux)"
	@echo "  make run          运行"
	@echo "  make test         运行测试"
	@echo "  make test-cover   测试覆盖率报告"
	@echo "  make lint         代码检查"
	@echo "  make fmt          代码格式化"
	@echo "  make clean        清理构建产物"
	@echo "  make install      安装到本地 GOPATH"
	@echo "  make docker-build Docker 构建"
	@echo "  make docker-up    Docker 启动"
