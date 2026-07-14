# 🌍 OpenSource Globalizer AI

> AI-powered localization and maintenance assistant for open-source projects
>
> 面向开源项目的 AI 国际化与维护助手

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## What is OpenSource Globalizer AI?

OpenSource Globalizer AI helps **open-source maintainers build truly global communities** by automating:

- 📖 **README / documentation translation** — preserve Markdown structure across 10+ languages (via goldmark AST)
- 🔄 **GitHub Action integration** — auto-translate on push, open a PR automatically
- 🏷️ **Issue triage & translation** (V2) — detect language, auto-label, translate for non-native maintainers
- 📦 **Release Notes generation** (V3) — produce multi-language release notes from changelogs

All powered by AI. All integrated into the GitHub workflow you already use.

---

## Why?

Open-source is global by nature. Your users speak 中文, 日本語, 한국어, Español, Français, Deutsch…

But most maintainers cannot manually translate every README, every Issue, every Release Note.

> **OpenSource Globalizer AI reduces localization workload from hours → seconds.**

---

## Features

| Feature | Ver | Status | Description |
|---------|-----|--------|-------------|
| 📖 **README Translator** | v0.1 | 🚧 In Progress | Translate README.md to multiple languages, goldmark AST preserves all formatting |
| 🌐 **HTTP API** | v0.1 | 🚧 In Progress | REST API via Gin, POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | 📋 Planned | Auto-translate on push, create PR automatically |
| 🏷️ **Issue Assistant** | v0.3 | 📋 Planned | Detect issue language, auto-classify, translate |
| 📦 **Release Assistant** | v0.4 | 📋 Planned | Generate multi-language release notes |
| 🤖 **GitHub App** | v1.0 | 📋 Planned | Full bot integration with PR comments and review |

---

## Quick Start

### 1. Install

```bash
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```

### 2. Translate a README

```bash
export OPENAI_API_KEY="sk-..."

globalizer translate README.md --lang zh-CN,ja,ko,es
```

Output:

```
docs/
├── README.zh-CN.md
├── README.ja.md
├── README.ko.md
└── README.es.md
```

### 3. Start HTTP API Server

```bash
globalizer serve
# → http://localhost:8080/api/v1/translate
```

### 4. Add to GitHub Actions

```yaml
# .github/workflows/i18n.yml
name: AI Translation

on:
  push:
    paths:
      - README.md

jobs:
  translate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: opensource-globalizer/action@v1
        with:
          openai-key: ${{ secrets.OPENAI_API_KEY }}
          languages: zh-CN,ja,ko,es
```

---

## Architecture

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

See [docs/architecture.md](docs/architecture.md) for full design.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Language** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark (AST-level parsing) |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **Config** | viper (env + YAML merge) |
| **Logging** | zap (structured) |
| **GitHub** | go-github, GitHub Actions |
| **Deploy** | Docker, 单二进制分发 |

---

## Project Structure

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

## Version Roadmap

| Version | Timeline | Deliverable | Status |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07 (Week 1-2) | CLI README Translator + HTTP API | 🚧 In Progress |
| **v0.2.0** | 2026-07 (Week 3-4) | GitHub Action + Auto PR + Docker Image | 📋 Planned |
| **v0.3.0** | 2026-08 | Issue Language Detect + Translate + Label | 📋 Planned |
| **v0.4.0** | 2026-09 | Release Notes Multi-language Generation | 📋 Planned |
| **v1.0.0** | 2026-10 | GitHub App + Dashboard + Multi-AI-Provider | 📋 Planned |

See [docs/roadmap.md](docs/roadmap.md) for detailed milestones.

---

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

```bash
git clone https://github.com/ytc301/opensource-globalizer.git
make deps
make test
make build
```

---

## License

MIT © 2026 OpenSource Globalizer AI Contributors
