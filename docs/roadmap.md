# Development Roadmap

## OpenSource Globalizer AI

---

## 版本总览

| 版本 | 时间 | 主题 | 核心交付 | 状态 |
|------|------|------|---------|------|
| **v0.1.0** | 2026-07 Week 1-2 | 🏗️ MVP CLI | CLI + HTTP API — README 翻译 | 🚧 |
| **v0.2.0** | 2026-07 Week 3-4 | 🔄 Action | GitHub Action + Docker + 自动 PR | 📋 |
| **v0.3.0** | 2026-08 | 🏷️ Issue AI | Issue 语言检测 + 翻译 + 分类 | 📋 |
| **v0.4.0** | 2026-09 | 📦 Release AI | Release Notes 多语言生成 | 📋 |
| **v1.0.0** | 2026-10 | 🤖 GitHub App | 完整 GitHub App + 多 AI Provider | 📋 |

---

## v0.1.0 — MVP: CLI README Translator + HTTP API

> **目标**: `globalizer translate README.md --lang zh-CN,ja` 一行命令出结果。

**时间**: 2026-07-14 → 2026-07-28 (2 周)

| # | 任务 | 优先级 | 状态 |
|---|------|--------|------|
| 1 | 项目初始化 + Go 模块 + 目录骨架 | P0 | ✅ |
| 2 | 完整文档体系 (SRS / 架构 / API / Roadmap) | P0 | ✅ |
| 3 | AI Provider 接口 + OpenAI 实现 | P0 | ✅ |
| 4 | goldmark AST Markdown 解析器 | P0 | ✅ |
| 5 | GORM + SQLite 翻译缓存 | P0 | ✅ |
| 6 | cobra CLI (translate / serve / languages / version) | P0 | ✅ |
| 7 | Gin HTTP Handler (POST /api/v1/translate) | P0 | ✅ |
| 8 | viper 配置管理 (env + YAML merge) | P0 | ✅ |
| 9 | zap 结构化日志 | P0 | ✅ |
| 10 | translate 命令串联 → 端到端跑通 | P0 | ⬜ |
| 11 | 单元测试 (coverage > 80%) | P0 | ⬜ |
| 12 | 添加 example 项目 + 截图 | P1 | ⬜ |

### v0.1.0 验收

- CLI `translate` 端到端生成正确的多语言 README
- goldmark 代码块/链接/Badge 全部保留
- `serve` 命令 HTTP API 正常响应
- 缓存命中后跳过 API 调用
- 测试覆盖率 > 80%

---

## v0.2.0 — GitHub Action + Docker

> **目标**: 用户配置 `.github/workflows/i18n.yml` → push README → 自动翻译 + 自动 PR。

**时间**: 2026-07-28 → 2026-08-10 (2 周)

| # | 任务 | 优先级 |
|---|------|--------|
| 1 | GitHub Action 开发 (action.yml + Dockerfile) | P0 |
| 2 | Docker 镜像构建 + 发布到 GHCR | P0 |
| 3 | Action 集成测试 (真实 GitHub 仓库) | P0 |
| 4 | peter-evans/create-pull-request 集成 | P0 |
| 5 | Action 文档 + 示例项目 | P0 |
| 6 | 翻译语言扩展到 15+ | P1 |

### v0.2.0 验收

- Action 在真实 GitHub 仓库中触发并成功创建 PR
- Docker 镜像可拉取并运行
- PR 标题/描述符合规范

---

## v0.3.0 — Issue Assistant

> **目标**: Issue 被自动检测语言 → 翻译为英文 → 添加分类标签。

**时间**: 2026-08-10 → 2026-08-31 (3 周)

| # | 任务 | 优先级 |
|---|------|--------|
| 1 | Issue Webhook 处理 | P0 |
| 2 | AI 语言检测 + 翻译服务 | P0 |
| 3 | AI Issue 自动分类 (bug/feature/question/doc) | P0 |
| 4 | 翻译结果作为 Issue Comment 自动回复 | P0 |
| 5 | 自动标签 (lang:xx / type:xx) | P0 |
| 6 | 集成测试 | P1 |

---

## v0.4.0 — Release Notes Generator

> **目标**: 输入 Changelog → 输出多语言 Release Notes。

**时间**: 2026-09 (3 周)

| # | 任务 | 优先级 |
|---|------|--------|
| 1 | Changelog 解析器 | P0 |
| 2 | Release Notes 多语言生成引擎 | P0 |
| 3 | GitHub Release API 集成 | P1 |
| 4 | AI Provider 接口扩展预留 | P2 |

---

## v1.0.0 — GitHub App + 多后端

> **目标**: 在 GitHub Marketplace 发布完整的 GitHub App。

**时间**: 2026-10 (4 周)

| # | 任务 | 优先级 |
|---|------|--------|
| 1 | GitHub App 框架 (OAuth + 安装回调) | P0 |
| 2 | 仓库配置页面 (Dashboard) | P1 |
| 3 | 自定义 Prompt 模板支持 | P2 |
| 4 | GitLab 支持调研 | P3 |
| 5 | PostgreSQL 迁移 (可选，如用户量需要) | P3 |

---

## 90 天目标 (截至 2026-10)

用于 Codex for OSS 申请的核心指标：

| 指标 | 目标 | 当前 |
|------|------|------|
| GitHub Stars | 200+ | 0 |
| Forks | 20+ | 0 |
| 真实使用项目 | 10+ | 0 |
| 社区 Issue 交流 | 有 | 无 |
| Release 次数 | 3 次以上 | 0 |
| Contributors | 5+ | 1 |

### 关键里程碑

```
Week 1-2  → v0.1.0 发布  ← 当前
Week 3-4  → v0.2.0 发布 (GitHub Action，首个真实用户)
Week 5-8  → v0.3.0 发布 (Issue Assistant，社区开始使用)
Week 9-12 → v0.4.0 发布 (Release Notes，功能完整)
Week 13-16 → v1.0.0 RC (GitHub App，提交 Codex for OSS 申请)
```
