# 🌍 OpenSource Globalizer AI

> 基于 AI 的开源项目本地化与维护助手
>
> 面向开源项目的 AI 国际化与维护助手

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## OpenSource Globalizer AI 是什么？

OpenSource Globalizer AI 通过自动化以下工作，帮助**开源维护者构建真正的全球化社区**：

- 📖 **README / 文档翻译** — 通过 goldmark AST 在 10+ 种语言中保留 Markdown 结构
- 🔄 **GitHub Action 集成** — push 时自动翻译，自动打开 PR
- 🏷️ **Issue 分类与翻译** (V2) — 检测语言、自动打标签、为非母语维护者翻译
- 📦 **发布说明生成** (V3) — 从变更日志生成多语言发布说明

全部由 AI 驱动。全部集成到你已在使用的 GitHub 工作流中。

---

## 为什么？

开源本质上是全球化的。你的用户使用中文、日本語、한국어、Español、Français、Deutsch……但大多数维护者无法手动翻译每一个 README、每一个 Issue、每一条发布说明。

> **OpenSource Globalizer AI 将本地化工作量从数小时缩短至数秒。**

---

## 功能特性

| 功能 | 版本 | 状态 | 描述 |
|---------|-----|--------|-------------|
| 📖 **README 翻译器** | v0.1 | ✅ 已发布 | 将 README.md 翻译为多种语言，goldmark AST 保留所有格式 |
| 🌐 **HTTP API** | v0.1 | ✅ 已发布 | 基于 Gin 的 REST API，POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ 已发布 | push 时自动翻译，自动创建 PR |
| 🏷️ **Issue 助手** | v0.3 | 📋 计划中 | 检测 Issue 语言，自动分类、翻译 |
| 📦 **发布助手** | v0.4 | 📋 计划中 | 生成多语言发布说明 |
| 🤖 **GitHub App** | v1.0 | 📋 计划中 | 完整的机器人集成，支持 PR 评论与审查 |

---

## 快速开始

> 📖 完整安装指南见 **[INSTALL.md](INSTALL.md)** — 从零到首次翻译，5 分钟上手。

### 一行命令安装

```bash
# 下载预编译二进制（macOS/Linux/Windows）或：
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker 运行（无需 Go 环境）

> **前提**：`-v $(pwd):/workspace` 会把当前目录挂载进容器，所以**必须在 README.md 所在目录执行**。
> macOS Docker Desktop 默认只共享 `/Users` 等目录 — 若在 `/tmp` 等非共享目录执行，容器内看不到文件（报 `no such file or directory`）。
> 解决：在 Docker Desktop → Settings → Resources → File Sharing 中添加该目录，或改用 `/Users/...` 下的路径。

```bash
# 拉取镜像
docker pull ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0

# CLI 模式：翻译当前目录 README（需在 README.md 所在目录执行）
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 \
  translate README.md --lang zh-CN,ja

# CLI 模式：自定义 API 地址和模型
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 \
  translate README.md --lang zh-CN --base-url https://api.example.com/v1 -m gpt-5-mini

# HTTP API 模式（无需挂载）
docker run -d -p 8080:8080 -e OPENAI_API_KEY="sk-xxx" \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 serve
```


### 一行命令翻译

```bash
# 环境变量方式
export OPENAI_API_KEY="sk-..."
globalizer translate README.md --lang zh-CN,ja,ko,es

# 命令行参数方式
globalizer translate README.md --lang zh-CN,ja,ko,es --api-key "sk-..."

# 指定 API 地址和模型
globalizer translate README.md --lang zh-CN -m gpt-5-mini \
  --base-url https://api.openai.com/v1 --api-key "sk-..."
```


```
📖 源文件: README.md
🌍 目标语言: zh-CN, ja, ko, es
🤖 模型: gpt-4o
  ✅ docs/README.zh-CN.md
  ✅ docs/README.ja.md
  ✅ docs/README.ko.md
  ✅ docs/README.es.md
✨ 翻译完成！共 4 个语言版本
```


### 启动 HTTP API

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


---

## 架构

```
                 GitHub
                    |
            GitHub Action / CLI
                    |
        ┌───────────────────────┐
        |  OpenSource Globalizer |
        |────────────────────────|
        |  ┌─────────────────┐   |
        |  │  Gin HTTP API   │   |
        |  │  (serve 命令)    │   |
        |  ├─────────────────┤   |
        |  │  Translator     │───|── OpenAI API (GPT-4o)
        |  │  (goldmark AST) │   |
        |  ├─────────────────┤   |
        |  │  GORM + SQLite  │   |  ← 翻译缓存
        |  └─────────────────┘   |
        └───────────────────────┘
```


完整设计请参阅 [docs/architecture.md](docs/architecture.md)。

---

## 技术栈

| 层 | 技术 |
|-------|-----------|
| **语言** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark（AST 级解析） |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **配置** | viper（环境变量 + YAML 合并） |
| **日志** | zap（结构化） |
| **GitHub** | go-github、GitHub Actions |
| **部署** | Docker、单二进制分发 |

---

## 项目结构

```
opensource-globalizer/
├── cmd/
│   └── globalizer/            # CLI 入口 (cobra + zap)
│       ├── main.go            # 根命令
│       ├── translate.go       # translate 子命令
│       ├── serve.go           # HTTP API 服务 (Gin)
│       └── commands.go        # version / languages 辅助命令
├── internal/
│   ├── handler/               # Gin HTTP Handler
│   ├── translator/            # 翻译引擎 (goldmark → AI → 重组)
│   ├── ai/                    # AI Provider 接口 + OpenAI 实现 + Mock
│   ├── github/                # GitHub Client 接口 + Mock
│   └── store/                 # GORM + SQLite 翻译缓存
├── pkg/
│   ├── markdown/              # goldmark AST 解析 + 片段管理
│   └── config/                # viper 配置管理
├── docs/
│   ├── srs.md                 # 软件需求规格
│   ├── architecture.md        # 架构设计文档
│   ├── api.md                 # API 接口设计
│   └── roadmap.md             # 版本路线图
├── github-action/
│   └── action.yml             # GitHub Action 定义
├── configs/
│   └── config.example.yaml    # 配置模板
├── tests/
├── docker-compose.yml         # 单容器部署
├── Dockerfile                 # 多阶段构建 (Alpine)
├── Makefile
└── go.mod
```


---

## 版本路线图

| 版本 | 时间线 | 交付内容 | 状态 |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07（第 1-2 周） | CLI README 翻译器 + HTTP API | ✅ 已发布 |
| **v0.2.0** | 2026-07（第 3-4 周） | GitHub Action + 自动 PR + Docker 镜像 | ✅ 已发布 |
| **v0.3.0** | 2026-08 | Issue 语言检测 + 翻译 + 标签 | 📋 计划中 |
| **v0.4.0** | 2026-09 | 多语言发布说明生成 | 📋 计划中 |
| **v1.0.0** | 2026-10 | GitHub App + 仪表盘 + 多 AI 提供商 | 📋 计划中 |

详细里程碑请参阅 [docs/roadmap.md](docs/roadmap.md)。

---

## 贡献

欢迎贡献！请先阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

```bash
git clone https://github.com/ytc301/opensource-globalizer.git
make deps
make test
make build
```


---

## 许可证

MIT © 2026 OpenSource Globalizer AI 贡献者