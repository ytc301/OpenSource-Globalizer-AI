# Software Requirements Specification (SRS)

## OpenSource Globalizer AI

| 属性 | 值 |
|------|-----|
| **版本** | v0.1.0 |
| **状态** | Draft |
| **作者** | OpenSource Globalizer Team |
| **最后更新** | 2026-07-14 |

---

## 1. 概述 (Introduction)

### 1.1 项目愿景

OpenSource Globalizer AI 致力于让每一个开源项目都能轻松触达全球开发者。通过 AI 驱动的自动化工作流，帮助开源维护者消除语言障碍，降低多语言维护成本。

### 1.2 目标用户

| 用户角色 | 描述 | 核心诉求 |
|---------|------|---------|
| **开源维护者** | 拥有活跃仓库的 Maintainer | 减少文档翻译、Issue 分类的重复劳动 |
| **开源贡献者** | 非英语母语的 Contributor | 用母语参与 Issue 讨论，降低参与门槛 |
| **开源使用者** | 使用开源项目的开发者 | 以母语阅读 README、文档和 Release Notes |

### 1.3 核心价值主张

> 将开源维护者的本地化工作量从「小时级」降低到「秒级」。

### 1.4 版本策略

项目采用渐进式迭代，每个版本聚焦单一核心能力：

| 版本 | 周期 | 核心交付 |
|------|------|---------|
| **v0.1** | 2026-07 Week 1-2 | CLI + HTTP API — README 翻译 |
| **v0.2** | 2026-07 Week 3-4 | GitHub Action + Docker — 自动化 + 自动 PR |
| **v0.3** | 2026-08 | Issue 语言检测 + 翻译 + 自动标签 |
| **v0.4** | 2026-09 | Release Notes 多语言生成 |
| **v1.0** | 2026-10 | GitHub App + 多 AI Provider + Dashboard |

---

## 2. 功能需求 (Functional Requirements)

### 2.1 v0.1.0 — CLI README Translator + HTTP API

#### FR-001: README 多语言翻译 (CLI)

- **输入**: 单个 Markdown 文件路径
- **输出**: `{output_dir}/README.{lang}.md` 多语言文件
- **支持语言**: `zh-CN`, `ja`, `ko`, `es`, `fr`, `de`, `pt-BR`, `ru`, `ar`, `en`
- **Markdown 保持规则** (由 goldmark AST 保证):
  - ✅ 标题层级 (`#`, `##`, `###` …) 保持不变
  - ✅ 代码块 (` ``` `) 内容不翻译，语言标签保留
  - ✅ 内联代码 (` `` `) 不翻译
  - ✅ URL 链接不破坏
  - ✅ Badge (`[![...]](...)`) 保留原文
  - ✅ HTML 标签属性不翻译

#### FR-002: CLI 接口

```
globalizer translate <file> [flags]

Flags:
  --lang, -l    string  目标语言，逗号分隔 (默认: "zh-CN")
  --output, -o  string  输出目录 (默认: "docs")
  --model        string  OpenAI 模型 (默认: "gpt-4o")
  --config       string  配置文件路径 (默认: .globalizer.yaml)
  --dry-run      bool    预览模式，不写入文件
  --source, -s   string  源语言 (留空则自动检测)
```

#### FR-003: HTTP API (Gin)

```
POST /api/v1/translate
GET  /api/v1/languages
GET  /health
```

#### FR-004: 配置文件 (viper)

```yaml
# .globalizer.yaml
languages: [zh-CN, ja, ko, es]
output_dir: docs
model: gpt-4o
db_path: ~/.globalizer/globalizer.db

openai:
  base_url: https://api.openai.com/v1

preserve:
  code_blocks: true
  links: true
  badges: true
  html_tags: true

server:
  port: 8080
  mode: debug
```

优先级: 命令行参数 > 环境变量 > 配置文件 > 默认值

#### FR-005: 翻译缓存 (SQLite)

- 相同内容（SHA-256 哈希）不重复调用 API
- 缓存键: source_hash + target_lang
- 缓存存储: GORM + SQLite (`~/.globalizer/globalizer.db`)
- MVP 阶段缓存可选（Store 为 nil 时跳过）

### 2.2 v0.2.0 — GitHub Action

#### FR-100: GitHub Action 集成

- 监听 `README.md` 的 `push` 事件
- 自动触发翻译流程
- 自动创建包含翻译文件的 Pull Request
- PR 标题格式: `🌍 i18n: Auto-translate README to {languages}`

#### FR-101: Docker 镜像

- 发布到 GitHub Container Registry (GHCR)
- 多阶段构建，Alpine 基础镜像
- Action 引用: `uses: opensource-globalizer/action@v1`

### 2.3 v0.3.0 — Issue Assistant

#### FR-200: Issue 语言检测 + 翻译

- 自动检测 Issue 正文语言 → 添加 `lang:xx` 标签
- 将非英语 Issue 翻译为英文摘要 → 评论贴出

#### FR-201: Issue 自动分类

- 基于内容分析自动添加标签: `type:bug` / `type:feature` / `type:question` / `type:documentation`

### 2.4 v0.4.0 — Release Assistant

#### FR-300: Release Notes 生成

- 输入: Git tag + Changelog
- 输出: 多语言 Release Notes (Markdown)
- 格式: GitHub Release 兼容

---

## 3. 非功能需求 (Non-Functional Requirements)

### NFR-001: 性能

| 指标 | 目标值 |
|------|--------|
| 单文件翻译响应时间 | < 30s (10KB README) |
| API 并发处理能力 | 100 req/s |
| Markdown 解析准确率 | 99%+ (goldmark CommonMark 兼容) |

### NFR-002: 可靠性

- OpenAI API 失败时自动重试 (最多 3 次, 指数退避)
- 部分语言翻译失败不影响其他语言输出
- 所有翻译操作幂等 (同一输入多次执行结果一致)
- SQLite 缓存命中直接返回，不依赖网络

### NFR-003: 安全性

- API Key 仅通过环境变量 `OPENAI_API_KEY` 或 GitHub Secrets 传入
- zap 日志自动脱敏，不记录 API Key
- 翻译缓存仅存储哈希后的内容，不缓存用户原始仓库数据到外部

### NFR-004: 可扩展性

- AI Provider 接口化 (`internal/ai.Provider`)，当前基于 OpenAI API，接口设计保留扩展性
- GORM 驱动，一条配置切换 SQLite → PostgreSQL
- goldmark AST 扩展机制，可自定义保留规则

### NFR-005: 可维护性

- 代码覆盖率 > 80%
- 架构决策记录 (ADR) 保存在 `docs/architecture.md` 第 6 节
- 每个模块有 Mock 实现，支持不依赖外部服务的测试

---

## 4. 技术约束 (Technical Constraints)

| 约束 | 说明 |
|------|------|
| **Go 1.23+** | 最低版本要求 |
| **SQLite** | v0.1-v0.3 使用 SQLite，v1.0 可选迁移 PostgreSQL |
| **无 Redis** | v0.1-v0.4 不引入消息队列或缓存中间件 |
| **单体架构** | 不拆 API + Worker 服务，GitHub Action 即执行环境 |
| **OpenAI API** | 当前基于 OpenAI API (GPT-4o)，接口设计支持扩展 |
| **GitHub 平台** | 仅支持 GitHub (GitLab 在 v1.0 后评估) |

---

## 5. 验收标准 (Acceptance Criteria)

### v0.1.0 验收

- [ ] `globalizer translate README.md --lang zh-CN,ja` 生成正确的中文和日文 README
- [ ] goldmark AST 正确识别并保留代码块、内联代码、链接、Badge
- [ ] `globalizer serve` 启动 Gin HTTP API 并正确响应
- [ ] POST `/api/v1/translate` 返回正确的翻译 JSON
- [ ] 缓存命中后不调用 OpenAI API
- [ ] 单元测试覆盖率 > 80%
- [ ] 通过 `go vet` 和 `golangci-lint` 检查
