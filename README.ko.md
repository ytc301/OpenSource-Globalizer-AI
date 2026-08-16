# 🌍 OpenSource Globalizer AI

> 오픈소스 프로젝트를 위한 AI 기반 현지화 및 유지보수 어시스턴트

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## OpenSource Globalizer AI란?

OpenSource Globalizer AI는 **오픈소스 유지보수자가 진정한 글로벌 커뮤니티를 구축**할 수 있도록 다음을 자동화합니다:

- 📖 **README / 문서 번역** — goldmark AST를 통해 10개 이상 언어에서 Markdown 구조 보존
- 🔄 **GitHub Action 통합** — push 시 자동 번역, PR 자동 생성
- 🏷️ **이슈 트리아지 및 번역** (V2) — 언어 감지, 자동 라벨링, 비원어민 유지보수자를 위한 번역
- 📦 **릴리스 노트 생성** (V3) — 변경로그에서 다국어 릴리스 노트 생성

모두 AI로 구동됩니다. 이미 사용 중인 GitHub 워크플로우에 완전히 통합됩니다.

---

## 왜 필요한가요?

오픈소스는 본질적으로 글로벌합니다. 사용자들은 中文, 日本語, 한국어, Español, Français, Deutsch…를 사용합니다. 하지만 대부분의 유지보수자는 모든 README, 모든 Issue, 모든 릴리스 노트를 수동으로 번역할 수 없습니다.

> **OpenSource Globalizer AI는 현지화 작업 시간을 몇 시간 → 몇 초로 단축합니다.**

---

## 기능

| 기능 | 버전 | 상태 | 설명 |
|---------|-----|--------|-------------|
| 📖 **README 번역기** | v0.1 | ✅ 출시됨 | README.md를 여러 언어로 번역, goldmark AST가 모든 서식 보존 |
| 🌐 **HTTP API** | v0.1 | ✅ 출시됨 | Gin 기반 REST API, POST /api/v1/translate |
| 🔄 **GitHub Action** | v0.2 | ✅ 출시됨 | push 시 자동 번역, PR 자동 생성 |
| 🏷️ **이슈 어시스턴트** | v0.3 | ✅ 출시됨 | 이슈 언어 감지, 자동 분류, 자동 답변 + 라벨 |
| 📦 **릴리스 어시스턴트** | v0.4 | 📋 예정 | 다국어 릴리스 노트 생성 |
| 🤖 **GitHub App** | v1.0 | 📋 예정 | PR 댓글 및 리뷰를 포함한 전체 봇 통합 |

---

## 빠른 시작

> 📖 전체 설치 가이드: **[INSTALL.md](INSTALL.md)** — 5분 워크스루.

### 한 줄 설치

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker (Go 환경 불필요)

> **사전 요구사항**: `-v $(pwd):/workspace`는 현재 디렉터리를 컨테이너에 마운트하므로 **`README.md`가 있는 디렉터리에서 실행**하세요.
> macOS Docker Desktop은 기본적으로 `/Users` 등만 공유합니다 — `/tmp` 또는 공유되지 않은 다른 경로에서 실행하면 컨테이너가 파일을 볼 수 없습니다 (`no such file or directory`).
> 해결 방법: Docker Desktop → Settings → Resources → File Sharing → 디렉터리 추가, 또는 `/Users/...` 경로를 사용하세요.

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


### 한 번의 명령으로 번역

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


### HTTP API 시작

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### GitHub Action: 자동 번역 + 자동 PR (권장)

로컬 환경이 필요 없습니다 — `README.md`를 push하면 번역 PR이 자동으로 생성됩니다:

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


번역된 파일은 `README.<lang>.md` 형식으로 저장소 루트에 작성됩니다 (예: `README.zh-CN.md`, `README.ja.md`).
**설정:**

1. `OPENAI_API_KEY` 시크릿 추가: 저장소 **Settings → Secrets and variables → Actions**
2. PR 생성 활성화: 저장소 **Settings → Actions → General → Workflow permissions** → *Allow GitHub Actions to create and approve pull requests* 체크
3. `README.md`를 push → Action이 자동 번역 → PR 생성 (제목 `🌍 i18n: Auto-translate README to ...`)

> `base-url` 및 `model` 입력을 통해 OpenAI 호환 엔드포인트를 지원합니다.
> 테스트용 API 키가 없나요? `api-key`를 비워 두고 `mock: true`를 추가하면 엔드투엔드 흐름을 검증할 수 있습니다.
> 예시 워크플로우는 [.github/workflows/i18n.yml](.github/workflows/i18n.yml)에 있습니다.

### 이슈 어시스턴트: 자동 감지, 분류, 답변 및 라벨링 (v0.3)

`serve` 모드는 Issues를 엔드투엔드로 자동 처리하는 GitHub 웹훅 엔드포인트를 노출합니다:

```bash
export GITHUB_TOKEN="ghp_xxx"                 # GitHub PAT with issues:write scope
export GLOBALIZER_WEBHOOK_SECRET="secret"     # webhook HMAC SHA-256 secret
globalizer serve                              # → POST /webhook is registered
```


**설정:**

1. 웹훅 추가: 저장소 **Settings → Webhooks → Add webhook**
   - **Payload URL**: `https://your-server/webhook`
   - **Content type**: `application/json`
   - **Secret**: `GLOBALIZER_WEBHOOK_SECRET`과 동일한 값
   - **Events**: **Issues** (`opened`, `edited`)
2. `GITHUB_TOKEN` 설정 (저장소 **Settings → Secrets and variables → Actions** 또는 서버 환경 변수) — `issues:write` 범위 필요

영어가 아닌 Issue가 열리면 어시스턴트가 자동으로:

1. 언어를 **감지** → `lang:xx` 라벨 추가
2. 이슈를 **분류** → `type:bug` / `type:feature` / `type:question` / `type:documentation` 라벨 추가
3. 영어 요약을 첫 번째 댓글로 **게시**:

```
## 🌐 AI Translation

**语言:** zh-CN

**摘要:** Install fails on Ubuntu 24.04
```


> 웹훅 요청은 HMAC SHA-256 (`X-Hub-Signature-256`)으로 검증됩니다. 유효하지 않은 서명은 `401`로 거부됩니다.
> 구성은 `.globalizer.yaml`의 `github.token` / `github.webhook_secret` 또는 `GITHUB_TOKEN` / `GLOBALIZER_WEBHOOK_SECRET` 환경 변수에 있습니다.
> 감지/분류 모델의 기본값은 `gpt-4o-mini`이며, `OPENAI_ISSUE_MODEL`로 재정의할 수 있습니다.

---

## 아키텍처

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


전체 설계는 [docs/architecture.md](docs/architecture.md)를 참조하세요.

---

## 기술 스택

| 계층 | 기술 |
|-------|-----------|
| **언어** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark (AST 수준 파싱) |
| **ORM** | GORM + SQLite |
| **AI** | OpenAI API (GPT-4o / Codex) |
| **설정** | viper (env + YAML 병합) |
| **로깅** | zap (구조화) |
| **GitHub** | go-github, GitHub Actions |
| **배포** | Docker, 단일 바이너리 배포 |

---

## 프로젝트 구조

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

## 버전 로드맵

| 버전 | 일정 | 산출물 | 상태 |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07 (1-2주) | CLI README 번역기 + HTTP API | ✅ 출시됨 |
| **v0.2.0** | 2026-07 (3-4주) | GitHub Action + 자동 PR + Docker 이미지 | ✅ 출시됨 |
| **v0.3.0** | 2026-08 | 이슈 언어 감지 + 분류 + 자동 답변 + 라벨 | ✅ 출시됨 |
| **v0.4.0** | 2026-09 | 다국어 릴리스 노트 생성 | 📋 예정 |
| **v1.0.0** | 2026-10 | GitHub App + 대시보드 + 멀티 AI 제공자 | 📋 예정 |

자세한 마일스톤은 [docs/roadmap.md](docs/roadmap.md)를 참조하세요.

---

## 기여

기여를 환영합니다! 먼저 [CONTRIBUTING.md](CONTRIBUTING.md)를 읽어 주세요.

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```


---

## 라이선스

MIT © 2026 OpenSource Globalizer AI 기여자들