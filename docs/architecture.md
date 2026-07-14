# Architecture Design Document

## OpenSource Globalizer AI

| 属性 | 值 |
|------|-----|
| **版本** | v0.1.0 |
| **状态** | Draft |
| **最后更新** | 2026-07-14 |

---

## 1. 架构概览

```
                              ┌──────────────────────────┐
                              │      GitHub.com           │
                              │                            │
                              │  ┌──────────────────────┐ │
                              │  │  User Repository      │ │
                              │  │  ┌──────────────────┐ │ │
                              │  │  │  README.md        │ │ │
                              │  │  │  .github/         │ │ │
                              │  │  │    workflows/     │ │ │
                              │  │  │      i18n.yml     │ │ │
                              │  │  └──────────────────┘ │ │
                              │  └──────────┬───────────┘ │
                              └─────────────┼─────────────┘
                                            │
                                    push: README.md
                                            │
                              ┌─────────────▼─────────────┐
                              │     GitHub Actions         │
                              │                             │
                              │  ┌───────────────────────┐ │
                              │  │  checkout@v4           │ │
                              │  │         ↓              │ │
                              │  │  opensource-globalizer │ │
                              │  │  /action@v1            │ │
                              │  │         ↓              │ │
                              │  │  Create Pull Request   │ │
                              │  └───────────────────────┘ │
                              └─────────────┬─────────────┘
                                            │
                                   CLI / HTTP API
                                            │
                              ┌─────────────▼─────────────┐
                              │   OpenSource Globalizer    │
                              │       Backend (Go)         │
                              │                             │
                              │  ┌───────────────────────┐ │
                              │  │   Entry Layer          │ │
                              │  │  cobra CLI / Gin HTTP  │ │
                              │  └───────────┬───────────┘ │
                              │              │              │
                              │  ┌───────────▼───────────┐ │
                              │  │   Service Layer        │ │
                              │  │  ┌───────────────────┐ │ │
                              │  │  │  Translator        │ │ │
                              │  │  │  · goldmark 解析   │ │ │
                              │  │  │  · AI 翻译         │ │ │
                              │  │  │  · 缓存命中        │ │ │
                              │  │  └────────┬──────────┘ │ │
                              │  │           │             │ │
                              │  │  ┌────────▼──────────┐ │ │
                              │  │  │   AI Provider      │ │ │
                              │  │  │   Interface        │ │ │
                              │  │  └────────┬──────────┘ │ │
                              │  └───────────┼────────────┘ │
                              │              │              │
                              │  ┌───────────▼───────────┐ │
                              │  │   Data Layer           │ │
                              │  │  ┌──────────────────┐  │ │
                              │  │  │  GORM + SQLite   │  │ │
                              │  │  │  (翻译缓存)       │  │ │
                              │  │  └──────────────────┘  │ │
                              │  └───────────────────────┘ │
                              └─────────────┬─────────────┘
                                            │
                                   HTTPS (OpenAI API)
                                            │
                              ┌─────────────▼─────────────┐
                              │      OpenAI API             │
                              │   (GPT-4o)                  │
                              └─────────────────────────────┘
```

---

## 2. 分层架构设计

采用 **Clean Architecture** 分层，由内向外：

```
┌──────────────────────────────────────┐
│         cmd/globalizer (入口)         │  ← CLI (cobra + zap) / HTTP serve (Gin)
├──────────────────────────────────────┤
│         internal/handler             │  ← Gin Handler (路由、校验、中间件)
├──────────────────────────────────────┤
│         internal/translator          │  ← 翻译业务逻辑 (Markdown解析 → AI翻译 → 缓存 → 重组)
├──────────────────────────────────────┤
│         internal/ai                  │  ← AI Provider 接口 (OpenAI / Mock)
│         internal/github              │  ← GitHub Client 接口 (go-github / Mock)
│         internal/store               │  ← GORM + SQLite 数据层
├──────────────────────────────────────┤
│         pkg/markdown                 │  ← goldmark AST 解析公共库
│         pkg/config                   │  ← viper 配置管理公共库
├──────────────────────────────────────┤
│         外部依赖                      │
│   GORM + SQLite | OpenAI API | go-github│
└──────────────────────────────────────┘
```

### 2.1 依赖方向

```
cmd → handler → translator → store/ai/github → pkg
                ↑ (依赖反转)
            interface 定义在各自包内
```

### 2.2 接口设计原则

- **AI Provider 接口** (`internal/ai`) — `internal/translator` 依赖接口而非具体实现，当前基于 OpenAI API
- **GitHub Client 接口** (`internal/github`) — 抽象 PR 创建、文件读取等操作，支持 Mock 测试
- **GORM** 作为 ORM 抽象 — 一行配置切换 SQLite ↔ PostgreSQL，v0.1 用 SQLite，v1.0 可选迁移

---

## 3. 核心模块设计

### 3.1 Translator 模块

```
internal/translator/
├── service.go         # 翻译服务 (解析 → AI → 缓存 → 重组)
└── service_test.go    # 单元测试 (待添加)
```

**核心流程：**

```
README.md
    │
    ▼
┌─────────────────┐
│ 1. goldmark 解析 │  构建 AST，识别所有节点类型
│                  │  Heading / CodeBlock / Link / Image / List / Paragraph
└────────┬─────────┘
         │
         ▼
┌─────────────────┐
│ 2. 片段分类      │  可翻译: Text, Heading, List, Blockquote
│                  │  保留:   CodeBlock, Image, Link, HTML, ThematicBreak
└────────┬─────────┘
         │
         ▼
┌─────────────────┐
│ 3. 检查缓存      │  SHA-256(源内容) → SQLite 查找
│   (GORM/SQLite)  │  命中 → 直接返回，跳过 API 调用
└────────┬─────────┘
         │
         ▼
┌─────────────────┐
│ 4. AI 翻译       │  System Prompt: 专业开源文档翻译
│   (OpenAI API)   │  规则: 保留 Markdown 结构、不翻代码块、保留 URL
└────────┬─────────┘
         │
         ▼
┌─────────────────┐
│ 5. 重组 + 写入   │  翻译结果 + 保留片段 → 按原始顺序重组
│                  │  写入缓存，输出 README.{lang}.md
└─────────────────┘
```

### 3.2 GitHub 集成模块

```
internal/github/
├── github.go          # Client 接口 + MockClient
```

**Action 执行流程：**

```
GitHub Action Trigger (push: README.md)
    │
    ▼
1. Checkout 仓库 (actions/checkout@v4)
    │
    ▼
2. 运行 Globalizer Action
    │
    ├── 读取 README.md
    ├── 调用 Translator 服务 (内嵌在 Runner 中)
    ├── 写入 docs/README.{lang}.md
    │
    ▼
3. 创建 Pull Request (peter-evans/create-pull-request@v6)
    │
    ├── Branch: i18n/translate-readme-{timestamp}
    ├── Commit: "🌍 i18n: Translate README to zh-CN, ja, ko, es"
    └── PR Body: 翻译摘要 + 文件列表
```

### 3.3 AI Provider 接口

```
internal/ai/
├── provider.go        # Provider 接口定义
├── openai.go          # OpenAI 实现 (Chat Completions API)
└── mock.go            # Mock 实现 (测试用)
```

```go
type Provider interface {
    Translate(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error)
    DetectLanguage(ctx context.Context, text string) (string, error)
    ClassifyIssue(ctx context.Context, title, body string) (*IssueClassifyResult, error)
}
```

- **当前实现**: `OpenAIProvider` (GPT-4o / GPT-4o-mini)
- **接口扩展**: Provider 接口支持未来接入其他 AI 后端
- **测试**: `MockProvider` (无需 API Key)

### 3.4 Gin HTTP Handler

```
internal/handler/
└── handler.go         # Gin 路由 + 中间件 + 处理器
```

**端点**:
- `GET /health` — 健康检查
- `POST /api/v1/translate` — 翻译 Markdown 内容
- `GET /api/v1/languages` — 支持的语言列表

---

## 4. 数据模型 (GORM + SQLite)

### 4.1 表设计

```go
type Translation struct {
    ID         uint   `gorm:"primaryKey"`
    SourceHash string `gorm:"uniqueIndex:idx_source_target;not null"`
    TargetLang string `gorm:"uniqueIndex:idx_source_target;not null"`
    SourceText string `gorm:"not null"`
    Translated string `gorm:"not null"`
    Model      string `gorm:"not null"`
    TokensUsed int
}
// GORM AutoMigrate 自动创建表和索引
```

### 4.2 设计说明

- **SQLite 单文件** — 数据库文件 `~/.globalizer/globalizer.db`，零配置零运维
- **缓存优先** — SHA-256(source) + targetLang 唯一索引，相同内容不再调 API
- **MVP 可选** — v0.1.0 CLI 模式下 Store 可以为 nil，纯文件翻译不依赖数据库
- **迁移零成本** — GORM 驱动，一行配置即可切换 SQLite → PostgreSQL

---

## 5. 部署架构

### 5.1 CLI 模式（MVP 推荐）

零依赖运行，GitHub Action Runner 即 Worker：

```
GitHub Action Runner
    │
    ├── actions/checkout@v4
    ├── go run cmd/globalizer translate ...
    └── peter-evans/create-pull-request@v6
```

### 5.2 Serve 模式

单容器部署，提供 HTTP API：

```yaml
# docker-compose.yml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - GLOBALIZER_DB_PATH=/data/globalizer.db
    volumes:
      - globalizer_data:/data

volumes:
  globalizer_data:
```

---

## 6. 架构决策记录 (ADR)

### ADR-001: 选择 Go 作为后端语言

**状态**: 已接受 (2026-07)

**决策**: 使用 Go 1.23+。

**理由**: 单二进制分发 → GitHub Action 友好；高性能；强类型适合长期维护。

### ADR-002: MVP 使用 SQLite + GORM，不引入独立数据库

**状态**: 已接受 (2026-07)

**决策**: v0.1-v0.3 使用 SQLite 本地文件存储，GORM 作为 ORM 层。

**理由**: 零配置、零运维；MVP 用户量不需要 PG 并发能力；GORM 可一行切换 SQLite → PostgreSQL，V1.0 迁移零成本。

### ADR-003: MVP 不引入 Redis / 消息队列

**状态**: 已接受 (2026-07)

**决策**: v0.1-v0.3 所有操作同步处理。CLI translate 耗时 10-30 秒，无需异步。

**理由**: 翻译是 CPU/API-bound，不是 I/O-bound；GitHub Action Runner 本身就是执行环境；减少运维依赖。

### ADR-004: Markdown 解析使用 goldmark

**状态**: 已接受 (2026-07)

**决策**: 使用 goldmark 进行 AST 级 Markdown 解析，替代手写 parser。

**理由**: 成熟的 CommonMark 兼容库；AST 操作比正则更可靠；社区维护。

### ADR-005: 不拆分 API + Worker 服务

**状态**: 已接受 (2026-07)

**决策**: v0.1-v0.4 保持单体架构。CLI 和 HTTP serve 共享同一套 Service 层。

**理由**: GitHub Action 即 Worker，不需要独立消费者；降低前期架构复杂度；代码复用度高。
