# 🌍 OpenSource Globalizer AI

> オープンソースプロジェクト向けの AI 対応ローカライズ＆メンテナンスアシスタント
>
> オープンソースプロジェクト向けの AI 国際化・保守アシスタント

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## OpenSource Globalizer AI とは？

OpenSource Globalizer AI は、以下を自動化することで**オープンソースメンテナーが真にグローバルなコミュニティを構築する**のを支援します:

- 📖 **README / ドキュメント翻訳** — 10以上の言語でMarkdown構造を保持（goldmark AST経由）
- 🔄 **GitHub Action 統合** — push時に自動翻訳し、自動でPRを開く
- 🏷️ **Issue トリアージ＆翻訳** (V2) — 言語を検出し、自動ラベル付け、非ネイティブメンテナーのために翻訳
- 📦 **リリースノート生成** (V3) — チェンジログから多言語リリースノートを生成

すべてAIを搭載。すべて、すでに使っている GitHub ワークフローに統合されます。

---

## なぜ？

オープンソースは本質的にグローバルです。あなたのユーザーは中文、日本語、한국어、Español、Français、Deutsch…を話します。
しかし、ほとんどのメンテナーはすべてのREADME、すべてのIssue、すべてのリリースノートを手動で翻訳することはできません。

> **OpenSource Globalizer AI は、ローカライズ作業を数時間から数秒に短縮します。**

---

## 機能

| 機能 | Ver | ステータス | 説明 |
|---------|-----|--------|-------------|
| 📖 **README翻訳** | v0.1 | ✅ リリース済み | README.md を複数言語に翻訳。goldmark AST がすべての書式を保持 |
| 🌐 **HTTP API** | v0.1 | ✅ リリース済み | Gin による REST API、POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ リリース済み | push時に自動翻訳し、PRを自動作成 |
| 🏷️ **Issue アシスタント** | v0.3 | 📋 計画中 | Issueの言語を検出し、自動分類・翻訳 |
| 📦 **リリースアシスタント** | v0.4 | 📋 計画中 | 多言語リリースノートを生成 |
| 🤖 **GitHub App** | v1.0 | 📋 計画中 | PRコメントとレビューを含む完全なボット統合 |

---

## クイックスタート

> 📖 完全なインストールガイドは **[INSTALL.md](INSTALL.md)** を参照してください — ゼロから初回翻訳まで、5分でセットアップできます。

### ワンライナーインストール

```bash
# 下载预编译二进制（macOS/Linux/Windows）或：
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker での実行（Go 環境は不要）

> **前提**: `-v $(pwd):/workspace` により現在のディレクトリがコンテナにマウントされるため、**README.md があるディレクトリで実行する必要があります**。
> macOS の Docker Desktop はデフォルトで `/Users` などのディレクトリのみ共有します — `/tmp` などの非共有ディレクトリで実行すると、コンテナ内からファイルが見えません（`no such file or directory` エラーになります）。
> 解決策: Docker Desktop → Settings → Resources → File Sharing でそのディレクトリを追加するか、`/Users/...` 配下のパスを使用してください。

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


### ワンライナー翻訳

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


### HTTP API の起動

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


---

## アーキテクチャ

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


詳細な設計は [docs/architecture.md](docs/architecture.md) を参照してください。

---

## 技術スタック

| レイヤー | 技術 |
|-------|-----------|
| **言語** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark (ASTレベル解析) |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **設定** | viper (env + YAML マージ) |
| **ログ** | zap (構造化) |
| **GitHub** | go-github, GitHub Actions |
| **デプロイ** | Docker、単一バイナリ配布 |

---

## プロジェクト構成

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

## バージョンロードマップ

| バージョン | スケジュール | 成果物 | ステータス |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07 (第1〜2週) | CLI README翻訳 + HTTP API | ✅ リリース済み |
| **v0.2.0** | 2026-07 (第3〜4週) | GitHub Action + 自動PR + Dockerイメージ | ✅ リリース済み |
| **v0.3.0** | 2026-08 | Issue言語検出 + 翻訳 + ラベル付け | 📋 計画中 |
| **v0.4.0** | 2026-09 | 多言語リリースノート生成 | 📋 計画中 |
| **v1.0.0** | 2026-10 | GitHub App + ダッシュボード + マルチAIプロバイダー | 📋 計画中 |

詳細なマイルストーンは [docs/roadmap.md](docs/roadmap.md) を参照してください。

---

## コントリビューション

コントリビューション大歓迎です！ 最初に [CONTRIBUTING.md](CONTRIBUTING.md) をお読みください。

```bash
git clone https://github.com/ytc301/opensource-globalizer.git
make deps
make test
make build
```


---

## ライセンス

MIT © 2026 OpenSource Globalizer AI Contributors