# API Design Document

## OpenSource Globalizer AI

| 属性 | 值 |
|------|-----|
| **版本** | v0.1.0 |
| **状态** | Draft |
| **框架** | Gin + cobra |

---

## 1. 接口总览

| Method | Path | 描述 | 版本 | 状态 |
|--------|------|------|------|------|
| `POST` | `/api/v1/translate` | 翻译 Markdown 内容 | v0.1 | 🚧 |
| `GET` | `/api/v1/languages` | 支持的语言列表 | v0.1 | ✅ |
| `GET` | `/health` | 健康检查 | v0.1 | ✅ |

### 启动方式

```bash
# HTTP API 模式
globalizer serve
# → 默认监听 :8080，Gin debug 模式

# CLI 模式（无需启动服务）
globalizer translate README.md --lang zh-CN
```

### 响应格式

所有 API 统一使用 JSON：

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

错误：

```json
{
  "success": false,
  "data": null,
  "error": "错误描述"
}
```

---

## 2. v0.1 API 详细定义

### 2.1 POST /api/v1/translate

翻译 Markdown 文件内容。

**框架**: Gin `ShouldBindJSON` 校验

**Request:**

```json
{
  "content": "# My Project\n\nThis is a description.\n\n```go\nfunc main() {}\n```",
  "target_langs": ["zh-CN", "ja", "ko"],
  "model": "gpt-4o"
}
```

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `content` | string | ✅ | Markdown 原始内容 |
| `target_langs` | []string | ✅ | 目标语言列表 |
| `model` | string | ❌ | 模型，默认 `gpt-4o` |

**Response (200):**

```json
{
  "success": true,
  "translations": {
    "zh-CN": "# 我的项目\n\n这是一个描述。\n\n```go\nfunc main() {}\n```",
    "ja": "# 私のプロジェクト\n\nこれは説明です。\n\n```go\nfunc main() {}\n```"
  }
}
```

### 2.2 GET /api/v1/languages

**Response (200):**

```json
{
  "success": true,
  "languages": [
    {"code": "en", "name": "English", "native_name": "English"},
    {"code": "zh-CN", "name": "Chinese (Simplified)", "native_name": "简体中文"},
    {"code": "ja", "name": "Japanese", "native_name": "日本語"},
    {"code": "ko", "name": "Korean", "native_name": "한국어"},
    {"code": "es", "name": "Spanish", "native_name": "Español"},
    {"code": "fr", "name": "French", "native_name": "Français"},
    {"code": "de", "name": "German", "native_name": "Deutsch"},
    {"code": "pt-BR", "name": "Portuguese (Brazil)", "native_name": "Português (Brasil)"},
    {"code": "ru", "name": "Russian", "native_name": "Русский"},
    {"code": "ar", "name": "Arabic", "native_name": "العربية"}
  ]
}
```

### 2.3 GET /health

**Response (200):**

```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

---

## 3. CLI 接口设计

### 3.1 命令结构

```
globalizer [command] [flags]

Commands:
  translate    翻译 Markdown 文件为多语言版本
  serve        启动 HTTP API 服务 (Gin)
  languages    列出支持的语言
  version      显示版本信息
  help         帮助信息
```

### 3.2 translate 命令

```bash
globalizer translate README.md --lang zh-CN,ja,ko

# 所有参数
--lang, -l    目标语言，逗号分隔       (默认: zh-CN)
--output, -o  输出目录                (默认: docs)
--model        OpenAI 模型            (默认: gpt-4o)
--config       配置文件路径            (默认: .globalizer.yaml)
--dry-run      预览模式，不写入文件
--source, -s   源语言，留空自动检测
```

### 3.3 serve 命令

```bash
globalizer serve

# 等价于启动 Gin HTTP 服务
# 读取 .globalizer.yaml 中的 server.port 配置
# 默认监听 :8080
```

### 3.4 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 通用错误 (zap.Fatal) |
| 2 | 参数/配置错误 |
| 3 | 文件读取错误 |
| 4 | API 调用错误 |
| 5 | 文件写入错误 |

---

## 4. V2/V3 计划接口 (暂未实现)

| Method | Path | 描述 | 计划版本 |
|--------|------|------|---------|
| `POST` | `/api/v1/issues/analyze` | Issue 分析 + 分类 | v0.3 |
| `POST` | `/api/v1/issues/translate` | Issue 翻译 | v0.3 |
| `POST` | `/api/v1/releases/generate` | Release Notes 生成 | v0.4 |
