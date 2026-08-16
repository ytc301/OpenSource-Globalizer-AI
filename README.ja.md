# 🌍 OpenSource Globalizer AI

> オープンソースプロジェクト向けのAI搭載ローカライゼーション＆メンテナンスアシスタント

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## OpenSource Globalizer AIとは？

OpenSource Globalizer AIは、以下の自動化により**オープンソースメンテナーが真にグローバルなコミュニティを構築する**のを支援します：

- 📖 **README / ドキュメント翻訳** — goldmark AST により10以上の言語でMarkdown構造を維持
- 🔄 **GitHub Action連携** — push時に自動翻訳し、自動でPRを開く
- 🏷️ **Issueのトリアージ＆翻訳** (V2) — 言語を検出し、自動ラベル付け、非ネイティブメンテナー向けに翻訳
- 📦 **リリースノート生成** (V3) — チェンジログから多言語リリースノートを生成

すべてAI駆動。すべてお使いのGitHubワークフローに統合されています。

---

## なぜ？

オープンソースは本来グローバルです。ユーザーは中文、日本語、한국어、Español、Français、Deutsch…を話します。
しかし、ほとんどのメンテナーはすべてのREADME、すべてのIssue、すべてのリリースノートを手動で翻訳することはできません。

> **OpenSource Globalizer AIはローカライズ作業を数時間から数秒に短縮します。**

---

## 機能

| 機能 | バージョン | ステータス | 説明 |
|---------|-----|--------|-------------|
| 📖 **README Translator** | v0.1 | ✅ リリース済み | README.mdを複数言語へ翻訳。goldmark ASTがすべての書式を保持 |
| 🌐 **HTTP API** | v0.1 | ✅ リリース済み | GinによるREST API、POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ リリース済み | push時に自動翻訳、PRを自動作成 |
| 🏷️ **Issue Assistant** | v0.3 | ✅ リリース済み | Issueの言語を検出し、自動分類・自動返信＋ラベル付け |
| 📦 **Release Assistant** | v0.4 | 📋 予定 | 多言語リリースノートを生成 |
| 🤖 **GitHub App** | v1.0 | 📋 予定 | PRコメントとレビューを含む完全なボット統合 |

---

## クイックスタート

> 📖 完全なインストールガイド： **[INSTALL.md](INSTALL.md)** — 5分で完了するチュートリアル。

### 一行インストール

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker（Go環境不要）

> **前提条件**: `-v $(pwd):/workspace` は現在のディレクトリをコンテナにマウントするため、**`README.md`があるディレクトリから実行してください**。
> macOSのDocker Desktopはデフォルトで`/Users`などしか共有しません — `/tmp`や他の共有されていないパスから実行すると、コンテナはファイルを認識できません（`no such file or directory`）。
> 解決策：Docker Desktop → Settings → Resources → File Sharing でディレクトリを追加するか、`/Users/...`配下のパスを使用してください。

```bash
# Pull the image
docker pull ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0

# CLI mode: translate the README in the current directory (run from README.md's directory)
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 \
  translate README.md --lang zh-CN,ja

# CLI mode: custom API endpoint and model
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 \
  translate README.md --lang zh-CN --base-url https://api.example.com/v1 -m gpt-5-mini

# HTTP API mode (no mount needed)
docker run -d -p 8080:8080 -e OPENAI_API_KEY="sk-xxx" \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 serve
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


### HTTP APIを起動

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### GitHub Action：自動翻訳＋自動PR（推奨）

ローカル環境は不要です — `README.md`をpushすると翻訳PRが自動的に作成されます：

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


翻訳ファイルはリポジトリルートに`README.<lang>.md`として出力されます（例：`README.zh-CN.md`、`README.ja.md`）。
**セットアップ：**

1. `OPENAI_API_KEY`シークレットを追加：リポジトリの**Settings → Secrets and variables → Actions**
2. PR作成を有効化：リポジトリの**Settings → Actions → General → Workflow permissions** → *Allow GitHub Actions to create and approve pull requests*にチェック
3. `README.md`をpush → Actionが自動翻訳 → PRが作成されます（タイトル `🌍 i18n: Auto-translate README to ...`）

> `base-url`と`model`入力を介してOpenAI互換エンドポイント（例：DeepSeek）をサポートします。
> テスト用のAPIキーがない場合？ `api-key`を空のままにして`mock: true`を追加すると、フローをエンドツーエンドで検証できます。
> サンプルワークフローは[.github/workflows/i18n.yml](.github/workflows/i18n.yml)にあります。

### Issue Assistant：自動検出・分類・返信・ラベル付け（v0.3）

`serve`モードはGitHub webhookエンドポイントを公開し、Issueをエンドツーエンドで自動処理します：

```bash
export GITHUB_TOKEN="ghp_xxx"                 # issues:write スコープのGitHub PAT
export GLOBALIZER_WEBHOOK_SECRET="secret"     # webhook HMAC SHA-256 シークレット
globalizer serve                              # → POST /webhook が登録される
```

**セットアップ：**

1. webhookを追加：リポジトリの**Settings → Webhooks → Add webhook**
   - **Payload URL**：`https://your-server/webhook`
   - **Content type**：`application/json`
   - **Secret**：`GLOBALIZER_WEBHOOK_SECRET`と同じ値
   - **Events**：**Issues**（`opened`、`edited`）
2. `GITHUB_TOKEN`を設定（リポジトリの**Settings → Secrets and variables → Actions**、またはサーバー環境変数）、`issues:write`スコープが必要

英語以外のIssueが開かれると、アシスタントは自動的に：

1. 言語を**検出** → `lang:xx`ラベルを追加
2. **分類** → `type:bug` / `type:feature` / `type:question` / `type:documentation`ラベルを追加
3. 最初のコメントとして英語サマリーを**投稿**：

```
## 🌐 AI Translation

**言語:** zh-CN

**摘要:** Install fails on Ubuntu 24.04
```

> WebhookリクエストはHMAC SHA-256（`X-Hub-Signature-256`）で検証され、無効な署名は`401`で拒否されます。
> 設定は`.globalizer.yaml`の`github.token` / `github.webhook_secret`、または環境変数`GITHUB_TOKEN` / `GLOBALIZER_WEBHOOK_SECRET`にあります。

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


詳細な設計は[docs/architecture.md](docs/architecture.md)を参照してください。

---

## 技術スタック

| レイヤー | 技術 |
|-------|-----------|
| **言語** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark（ASTレベルでの解析） |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **設定** | viper（env + YAMLマージ） |
| **ロギング** | zap（構造化） |
| **GitHub** | go-github、GitHub Actions |
| **デプロイ** | Docker、単一バイナリ配布 |

---

## プロジェクト構成

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
| **v0.1.0** | 2026-07（第1〜2週） | CLI README翻訳 + HTTP API | ✅ リリース済み |
| **v0.2.0** | 2026-07（第3〜4週） | GitHub Action + 自動PR + Dockerイメージ | ✅ リリース済み |
| **v0.3.0** | 2026-08 | Issue言語検出 + 翻訳 + ラベル | ✅ リリース済み |
| **v0.4.0** | 2026-09 | 多言語リリースノート生成 | 📋 予定 |
| **v1.0.0** | 2026-10 | GitHub App + ダッシュボード + マルチAIプロバイダー | 📋 予定 |

詳細なマイルストーンは[docs/roadmap.md](docs/roadmap.md)を参照してください。

---

## コントリビューション

コントリビューション大歓迎です！ まず[CONTRIBUTING.md](CONTRIBUTING.md)をお読みください。

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```


---

## ライセンス

MIT © 2026 OpenSource Globalizer AI Contributors