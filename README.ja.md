# 🌍 OpenSource Globalizer AI

> オープンソースプロジェクト向けの AI を活用したローカライゼーション & メンテナンスアシスタント

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## OpenSource Globalizer AI とは？

OpenSource Globalizer AI は、以下の自動化により **オープンソースメンテナーが真にグローバルなコミュニティを構築する** ことを支援します：

- 📖 **README / ドキュメント翻訳** — goldmark AST により 10 以上の言語で Markdown 構造を保持
- 🔄 **GitHub Action 統合** — プッシュ時に自動翻訳し、PR を自動作成
- 🏷️ **Issue トリアージ & 翻訳** (V2) — 言語を検出し、自動ラベル付け、非ネイティブメンテナー向けに翻訳
- 📦 **リリースノート生成** (V3) — チェンジログから多言語リリースノートを生成

すべて AI を搭載。すべて、あなたがすでに使用している GitHub ワークフローに統合されます。

---

## なぜ？

オープンソースは本質的にグローバルです。ユーザーは 中文、日本語、한국어、Español、Français、Deutsch などを話します。しかし、ほとんどのメンテナーはすべての README、すべての Issue、すべてのリリースノートを手動で翻訳することはできません。

> **OpenSource Globalizer AI はローカライゼーションの作業時間を数時間 → 数秒に短縮します。**

---

## 機能

| 機能 | バージョン | ステータス | 説明 |
|---------|-----|--------|-------------|
| 📖 **README 翻訳** | v0.1 | ✅ リリース済み | README.md を複数言語に翻訳。goldmark AST がすべての書式を保持 |
| 🌐 **HTTP API** | v0.1 | ✅ リリース済み | Gin による REST API。POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ リリース済み | プッシュ時に自動翻訳し、PR を自動作成 |
| 🏷️ **Issue アシスタント** | v0.3 | 📋 予定 | Issue の言語を検出し、自動分類、翻訳 |
| 📦 **リリースアシスタント** | v0.4 | 📋 予定 | 多言語リリースノートを生成 |
| 🤖 **GitHub App** | v1.0 | 📋 予定 | PR コメントとレビューを含む完全なボット統合 |

---

## クイックスタート

> 📖 完全なインストールガイド: **[INSTALL.md](INSTALL.md)** — 5分で完了するチュートリアル。

### 1行インストール

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker（Go 環境は不要）

> **前提条件**: `-v $(pwd):/workspace` はカレントディレクトリをコンテナにマウントします。そのため、**`README.md` を含むディレクトリから実行してください**。
> macOS の Docker Desktop はデフォルトで `/Users` などしか共有しません。`/tmp` やその他の共有されていないパスから実行すると、コンテナがファイルを認識できません（`no such file or directory`）。
> 修正方法: Docker Desktop → Settings → Resources → File Sharing でディレクトリを追加するか、`/Users/...` 配下のパスを使用してください。

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


### 1コマンドで翻訳

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


### HTTP API を起動

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### GitHub Action: 自動翻訳 + 自動 PR（推奨）

ローカル環境は不要です — `README.md` をプッシュすると、翻訳 PR が自動的に作成されます：

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


翻訳ファイルは `README.<lang>.md` としてリポジトリルートに書き込まれます（例: `README.zh-CN.md`、`README.ja.md`）。
**セットアップ:**

1. `OPENAI_API_KEY` シークレットを追加：リポジトリの **Settings → Secrets and variables → Actions**
2. PR 作成を有効化：リポジトリの **Settings → Actions → General → Workflow permissions** で *Allow GitHub Actions to create and approve pull requests* にチェックを入れる
3. `README.md` をプッシュ → Action が自動翻訳 → PR が作成されます（タイトル `🌍 i18n: Auto-translate README to ...`）

> OpenAI 互換エンドポイント（例: DeepSeek）を `base-url` と `model` 入力を介してサポートします。
> テスト用の API キーがありませんか？ `api-key` を空のままにして `mock: true` を追加すると、エンドツーエンドでフロー全体を検証できます。
> サンプルワークフローは [.github/workflows/i18n.yml](.github/workflows/i18n.yml) にあります。

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
        |  │  (serve cmd)    │   |
        |  ├─────────────────┤   |
        |  │  Translator     │───|── OpenAI API (GPT-4o)
        |  │  (goldmark AST) │   |
        |  ├─────────────────┤   |
        |  │  GORM + SQLite  │   |  ← translation cache
        |  └─────────────────┘   |
        └───────────────────────┘
```


完全な設計は [docs/architecture.md](docs/architecture.md) を参照してください。

---

## 技術スタック

| レイヤー | 技術 |
|-------|-----------|
| **言語** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark（AST レベルでのパース） |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **設定** | viper（env + YAML マージ） |
| **ロギング** | zap（構造化） |
| **GitHub** | go-github、GitHub Actions |
| **デプロイ** | Docker、シングルバイナリ配布 |

---

## プロジェクト構造

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

## バージョンロードマップ

| バージョン | 時期 | 成果物 | ステータス |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07（第1〜2週） | CLI README 翻訳 + HTTP API | ✅ リリース済み |
| **v0.2.0** | 2026-07（第3〜4週） | GitHub Action + 自動 PR + Docker イメージ | ✅ リリース済み |
| **v0.3.0** | 2026-08 | Issue 言語検出 + 翻訳 + ラベル | 📋 予定 |
| **v0.4.0** | 2026-09 | 多言語リリースノート生成 | 📋 予定 |
| **v1.0.0** | 2026-10 | GitHub App + ダッシュボード + マルチ AI プロバイダー | 📋 予定 |

詳細なマイルストーンは [docs/roadmap.md](docs/roadmap.md) を参照してください。

---

## コントリビューション

コントリビューションを歓迎します！ まず [CONTRIBUTING.md](CONTRIBUTING.md) をお読みください。

```bash
git clone https://github.com/ytc301/opensource-globalizer.git
make deps
make test
make build
```


---

## ライセンス

MIT © 2026 OpenSource Globalizer AI コントリビューター