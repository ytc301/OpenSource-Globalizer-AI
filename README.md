# 🌍 OpenSource Globalizer AI

> AI-powered localization and maintenance assistant for open-source projects

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
| 📖 **README Translator** | v0.1 | ✅ Released | Translate README.md to multiple languages, goldmark AST preserves all formatting |
| 🌐 **HTTP API** | v0.1 | ✅ Released | REST API via Gin, POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ Released | Auto-translate on push, create PR automatically |
| 🏷️ **Issue Assistant** | v0.3 | 📋 Planned | Detect issue language, auto-classify, translate |
| 📦 **Release Assistant** | v0.4 | 📋 Planned | Generate multi-language release notes |
| 🤖 **GitHub App** | v1.0 | 📋 Planned | Full bot integration with PR comments and review |

---

## Quick Start

> 📖 Full installation guide: **[INSTALL.md](INSTALL.md)** — 5-minute walkthrough.

### One-line install

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```

### Docker (no Go environment required)

> **Prerequisite**: `-v $(pwd):/workspace` mounts the current directory into the container, so **run it from the directory containing `README.md`**.
> macOS Docker Desktop only shares `/Users` etc. by default — if you run from `/tmp` or another unshared path, the container won't see your files (`no such file or directory`).
> Fix: Docker Desktop → Settings → Resources → File Sharing → add the directory, or use a path under `/Users/...`.

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

### Translate in one command

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

### Start HTTP API

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```

### GitHub Action: auto-translate + auto-PR (recommended)

No local environment needed — push `README.md` and a translation PR is created automatically:

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

Translated files are written to the repo root as `README.<lang>.md` (e.g. `README.zh-CN.md`, `README.ja.md`).

**Setup:**

1. Add the `OPENAI_API_KEY` secret: repo **Settings → Secrets and variables → Actions**
2. Enable PR creation: repo **Settings → Actions → General → Workflow permissions** → check *Allow GitHub Actions to create and approve pull requests*
3. Push `README.md` → the Action auto-translates → a PR is created (title `🌍 i18n: Auto-translate README to ...`)

> Supports OpenAI-compatible endpoints (e.g. DeepSeek) via the `base-url` and `model` inputs.
> No API key for testing? Leave `api-key` empty and add `mock: true` to verify the flow end-to-end.
> Example workflow lives at [.github/workflows/i18n.yml](.github/workflows/i18n.yml).

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
        |  │  (serve cmd)    │   |
        |  ├─────────────────┤   |
        |  │  Translator     │───|── OpenAI API (GPT-4o)
        |  │  (goldmark AST) │   |
        |  ├─────────────────┤   |
        |  │  GORM + SQLite  │   |  ← translation cache
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
| **Deploy** | Docker, single-binary distribution |

---

## Project Structure

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

## Version Roadmap

| Version | Timeline | Deliverable | Status |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07 (Week 1-2) | CLI README Translator + HTTP API | ✅ Released |
| **v0.2.0** | 2026-07 (Week 3-4) | GitHub Action + Auto PR + Docker Image | ✅ Released |
| **v0.3.0** | 2026-08 | Issue Language Detect + Translate + Label | 📋 Planned |
| **v0.4.0** | 2026-09 | Release Notes Multi-language Generation | 📋 Planned |
| **v1.0.0** | 2026-10 | GitHub App + Dashboard + Multi-AI-Provider | 📋 Planned |

See [docs/roadmap.md](docs/roadmap.md) for detailed milestones.

---

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```

---

## License

MIT © 2026 OpenSource Globalizer AI Contributors
