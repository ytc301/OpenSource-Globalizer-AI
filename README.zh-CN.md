# 🌍 OpenSource Globalizer AI

> 面向开源项目的 AI 驱动本地化与维护助手

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## 什么是 OpenSource Globalizer AI？

OpenSource Globalizer AI 通过自动化以下工作，帮助**开源维护者构建真正的全球化社区**：

- 📖 **README / 文档翻译** — 通过 goldmark AST 在 10+ 种语言间保留 Markdown 结构
- 🔄 **GitHub Action 集成** — push 时自动翻译，自动创建 PR
- 🏷️ **Issue 分诊与翻译** (V2) — 检测语言、自动打标签、为非母语维护者提供翻译
- 📦 **发布说明生成** (V3) — 根据变更日志生成多语言发布说明

全部由 AI 驱动。全部集成到您已在使用的 GitHub 工作流中。

---

## 为什么？

开源本质上就是全球性的。您的用户使用中文、日本語、한국어、Español、Français、Deutsch…… 但大多数维护者无法手动翻译每一个 README、每一个 Issue、每一份发布说明。

> **OpenSource Globalizer AI 将本地化工作量从数小时缩短到数秒。**

---

## 功能特性

| 功能 | 版本 | 状态 | 描述 |
|---------|-----|--------|-------------|
| 📖 **README 翻译器** | v0.1 | ✅ 已发布 | 将 README.md 翻译为多种语言，goldmark AST 保留所有格式 |
| 🌐 **HTTP API** | v0.1 | ✅ 已发布 | 基于 Gin 的 REST API，POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ 已发布 | push 时自动翻译，自动创建 PR |
| 🏷️ **Issue 助手** | v0.3 | ✅ 已发布 | 检测 Issue 语言，自动分类，自动回复 + 标签 |
| 📦 **发布助手** | v0.4 | 📋 规划中 | 生成多语言发布说明 |
| 🤖 **GitHub App** | v1.0 | 📋 规划中 | 完整的机器人集成，支持 PR 评论和审查 |

---

## 快速开始

> 📖 完整安装指南：**[INSTALL.md](INSTALL.md)** — 5 分钟快速上手。

### 一行安装

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker（无需 Go 环境）

> **前提条件**：`-v $(pwd):/workspace` 会将当前目录挂载到容器中，因此**请从包含 `README.md` 的目录运行**。
> macOS Docker Desktop 默认只共享 `/Users` 等目录 — 如果您从 `/tmp` 或其他未共享的路径运行，容器将无法看到您的文件（`no such file or directory`）。
> 解决方法：Docker Desktop → Settings → Resources → File Sharing → 添加该目录，或使用 `/Users/...` 下的路径。

```bash
# Pull the image
docker pull ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0

# CLI mode: translate the README in the current directory (run from README.md's directory)
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 \
  translate README.md --lang zh-CN,ja

# CLI mode: custom API endpoint and model
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 \
  translate README.md --lang zh-CN --base-url https://api.example.com/v1 -m gpt-5-mini

# HTTP API mode (no mount needed)
docker run -d -p 8080:8080 -e OPENAI_API_KEY="sk-xxx" \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.2.0 serve
```


### 一条命令完成翻译

```bash
# Via environment variable
export OPENAI_API_KEY="sk-..."
globalizer translate README.md --lang zh-CN,ja,ko,es

# Via CLI flag
globalizer translate README.md --lang zh-CN,ja,ko,es --api-key "sk-..."

# Custom API endpoint and model
globalizer translate README.md --lang zh-CN -m gpt-5-mini \
  --base-url https://api.openai.com/v1 --api-key "sk-..."
```


```
📖 Source: README.md
🌍 Target languages: zh-CN, ja, ko, es
🤖 Model: gpt-4o
  ✅ docs/README.zh-CN.md
  ✅ docs/README.ja.md
  ✅ docs/README.ko.md
  ✅ docs/README.es.md
✨ Done! 4 language versions generated
```


### 启动 HTTP API

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### GitHub Action：自动翻译 + 自动 PR（推荐）

无需本地环境 — push `README.md` 后会自动创建翻译 PR：

```yaml
# .github/workflows/i18n.yml
name: AI Translation

on:
  push:
    branches: [main]       # tag pushes are excluded to keep PR creation reliable
    paths:
      - "README.md"
  workflow_dispatch:      # optional manual trigger

permissions:
  contents: write
  pull-requests: write

jobs:
  translate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: ytc301/OpenSource-Globalizer-AI/github-action@v0.2.0
        with:
          api-key: ${{ secrets.OPENAI_API_KEY }}
          languages: zh-CN,ja,ko
          model: gpt-4o                        # optional: custom model
          base-url: https://api.openai.com/v1  # optional: API-compatible endpoint
          output-dir: .                        # write README.<lang>.md at repo root
```


翻译后的文件将写入仓库根目录，命名为 `README.<lang>.md`（例如 `README.zh-CN.md`、`README.ja.md`）。
**设置方法：**

1. 添加 `OPENAI_API_KEY` 密钥：仓库 **Settings → Secrets and variables → Actions**
2. 启用 PR 创建：仓库 **Settings → Actions → General → Workflow permissions** → 勾选 *Allow GitHub Actions to create and approve pull requests*
3. 推送 `README.md` → Action 自动翻译 → 创建 PR（标题 `🌍 i18n: Auto-translate README to ...`）

> 通过 `base-url` 和 `model` 输入支持兼容 OpenAI 的端点（例如 DeepSeek）。
> 没有 API 密钥用于测试？将 `api-key` 留空并添加 `mock: true` 即可端到端验证流程。
> 示例工作流位于 [.github/workflows/i18n.yml](.github/workflows/i18n.yml)。

### Issue 助手：自动检测、分类、回复与标签（v0.3）

`serve` 模式暴露一个 GitHub webhook 端点，端到端自动处理 Issue：

```bash
export GITHUB_TOKEN="ghp_xxx"                 # 具有 issues:write 权限的 GitHub PAT
export GLOBALIZER_WEBHOOK_SECRET="secret"     # webhook HMAC SHA-256 密钥
globalizer serve                              # → 注册 POST /webhook
```

**设置方法：**

1. 添加 webhook：仓库 **Settings → Webhooks → Add webhook**
   - **Payload URL**：`https://your-server/webhook`
   - **Content type**：`application/json`
   - **Secret**：与 `GLOBALIZER_WEBHOOK_SECRET` 相同的值
   - **Events**：**Issues**（`opened`、`edited`）
2. 设置 `GITHUB_TOKEN`（仓库 **Settings → Secrets and variables → Actions**，或服务器环境变量），需 `issues:write` 权限

当非英语 Issue 被打开时，助手会自动：

1. **检测**语言 → 添加 `lang:xx` 标签
2. **分类** → 添加 `type:bug` / `type:feature` / `type:question` / `type:documentation` 标签
3. **发布**英文摘要作为第一条评论：

```
## 🌐 AI Translation

**语言:** zh-CN

**摘要:** Install fails on Ubuntu 24.04
```

> Webhook 请求通过 HMAC SHA-256（`X-Hub-Signature-256`）校验；无效签名以 `401` 拒绝。
> 配置位于 `.globalizer.yaml` 的 `github.token` / `github.webhook_secret`，或环境变量 `GITHUB_TOKEN` / `GLOBALIZER_WEBHOOK_SECRET`。

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
        |  │  (serve cmd)    │   |
        |  ├─────────────────┤   |
        |  │  Translator     │───|── OpenAI API (GPT-4o)
        |  │  (goldmark AST) │   |
        |  ├─────────────────┤   |
        |  │  GORM + SQLite  │   |  ← translation cache
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
| **GitHub** | go-github, GitHub Actions |
| **部署** | Docker，单二进制分发 |

---

## 项目结构

```
opensource-globalizer/
├── cmd/
│   └── globalizer/            # CLI entry (cobra + zap)
│       ├── main.go            # root command
│       ├── translate.go       # translate subcommand
│       ├── serve.go           # HTTP API server (Gin)
│       └── commands.go        # version / languages helpers
├── internal/
│   ├── handler/               # Gin HTTP handlers
│   ├── translator/            # translation engine (goldmark → AI → rebuild)
│   ├── ai/                    # AI Provider interface + OpenAI impl + Mock
│   ├── github/                # GitHub Client interface + Mock
│   └── store/                 # GORM + SQLite translation cache
├── pkg/
│   ├── markdown/              # goldmark AST parsing + segment management
│   └── config/                # viper configuration
├── docs/
│   ├── srs.md                 # Software Requirements Specification
│   ├── architecture.md        # architecture design
│   ├── api.md                 # API design
│   └── roadmap.md             # version roadmap
├── github-action/
│   └── action.yml             # GitHub Action definition
├── configs/
│   └── config.example.yaml    # config template
├── tests/
├── docker-compose.yml         # single-container deploy
├── Dockerfile                 # multi-stage build (Alpine)
├── Makefile
└── go.mod
```


---

## 版本路线图

| 版本 | 时间线 | 交付物 | 状态 |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07（第 1-2 周） | CLI README 翻译器 + HTTP API | ✅ 已发布 |
| **v0.2.0** | 2026-07（第 3-4 周） | GitHub Action + 自动 PR + Docker 镜像 | ✅ 已发布 |
| **v0.3.0** | 2026-08 | Issue 语言检测 + 翻译 + 标签 | ✅ 已发布 |
| **v0.4.0** | 2026-09 | 发布说明多语言生成 | 📋 规划中 |
| **v1.0.0** | 2026-10 | GitHub App + 仪表盘 + 多 AI 提供商 | 📋 规划中 |

详细里程碑请参阅 [docs/roadmap.md](docs/roadmap.md)。

---

## 贡献

欢迎贡献！请先阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```


---

## 许可证

MIT © 2026 OpenSource Globalizer AI 贡献者