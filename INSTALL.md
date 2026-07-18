# 📥 安装与使用指南

## OpenSource Globalizer AI

> 面向开源项目的 AI 国际化助手 — 从零到首次翻译，5 分钟上手。

---

## 第一步：获取程序

### 方式一：下载预编译二进制（推荐，无需 Go 环境）

从 [GitHub Releases](https://github.com/ytc301/OpenSource-Globalizer-AI/releases) 下载对应平台：

| 平台 | 文件名 |
|------|--------|
| **macOS (Intel)** | `globalizer-darwin-amd64` |
| **macOS (Apple M1/M2/M3)** | `globalizer-darwin-arm64` |
| **Linux (x86_64)** | `globalizer-linux-amd64` |
| **Windows (x86_64)** | `globalizer-windows-amd64.exe` |

```bash
# macOS / Linux 下载后重命名并加执行权限
mv globalizer-darwin-arm64 globalizer
chmod +x globalizer

# 验证
./globalizer version
# → globalizer v0.1.0
```

### 方式二：Go 安装（需 Go 1.23+）

```bash
go install github.com/ytc301/OpenSource-Globalizer-AI/cmd/globalizer@latest
globalizer version
```

### 方式三：从源码编译

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
cd OpenSource-Globalizer-AI
make build
./bin/globalizer version
```

---

## 第二步：获取 OpenAI API Key

1. 访问 [platform.openai.com/api-keys](https://platform.openai.com/api-keys)
2. 登录或注册 OpenAI 账号
3. 点击「Create new secret key」创建 Key
4. 复制 Key（格式：`sk-proj-...` 或 `sk-...`）

> ⚠️ API Key 需要对应项目有 GPT-4o 模型访问权限并有足够配额。

---

## 第三步：翻译你的第一个 README

### 翻译当前目录的 README.md

```bash
# 设置 API Key（仅当前终端有效，关闭后消失）
export OPENAI_API_KEY="sk-你的密钥"

# 翻译为中文
globalizer translate README.md --lang zh-CN

# 翻译为多语言
globalizer translate README.md --lang zh-CN,ja,ko,es
```

**执行后：**

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

### 预览模式（不写文件，仅屏幕输出）

```bash
globalizer translate README.md --lang zh-CN --dry-run
```

### 使用更便宜的模型

```bash
globalizer translate README.md --lang zh-CN -m gpt-4o-mini
```

### 无需 API Key 也可测试

```bash
# Mock 模式：使用模拟数据验证全链路
globalizer translate README.md --lang zh-CN --mock --dry-run
```

---

## 第四步：校验翻译结果

翻译完成后检查：

| 检查项 | 怎么验证 |
|--------|---------|
| 代码块完整 | 打开翻译文件，确认 ` ```go ... ``` ` 原样保留 |
| 链接不破坏 | 确认 `[文字](URL)` 格式完整 |
| Badge 不丢 | 确认 `[![...]](...)` 标记完好 |
| 首次翻译耗时 | 10KB 的 README < 30 秒 |
| 二次翻译耗时 | < 1 秒（命中 SQLite 缓存，不调 API） |

---

## 第五步（可选）：启动 HTTP API 服务

```bash
# 启动服务（默认 :8080）
globalizer serve

# 另一个终端测试
curl -X POST http://localhost:8080/api/v1/translate \
  -H "Content-Type: application/json" \
  -d '{"content": "# Hello\n\nWorld.", "target_langs": ["zh-CN", "ja"]}'

# 返回
{"success":true,"translations":{"ja":"# こんにちは\n\n世界。","zh-CN":"# 你好\n\n世界。"}}
```

---

## 第六步（可选）：配置 GitHub Action 自动翻译

在你项目的 `.github/workflows/i18n.yml` 中添加：

```yaml
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
      - name: Translate README
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
        run: |
          go install github.com/ytc301/OpenSource-Globalizer-AI/cmd/globalizer@latest
          globalizer translate README.md --lang zh-CN,ja,ko,es
      - uses: peter-evans/create-pull-request@v6
        with:
          commit-message: "🌍 i18n: Auto-translate README"
          title: "🌍 i18n: Auto-translate README"
          branch: i18n/translate-readme
```

在 GitHub 仓库 Settings → Secrets → 添加 `OPENAI_API_KEY`。

完成后：每次修改 README.md 并推送，GitHub Action 自动翻译并创建 Pull Request。

---

## 命令速查

| 命令 | 说明 |
|------|------|
| `globalizer version` | 查看版本 |
| `globalizer languages` | 查看支持的语言 |
| `globalizer translate <file> --lang zh-CN` | 翻译文件 |
| `globalizer translate <file> --lang zh-CN --dry-run` | 预览模式 |
| `globalizer translate <file> --lang zh-CN --mock` | Mock 测试模式 |
| `globalizer serve` | 启动 HTTP API |
| `globalizer help` | 查看帮助 |

### 完整参数

```
translate [文件] [flags]

Flags:
  -l, --lang      目标语言，逗号分隔    (默认: zh-CN)
  -o, --output    输出目录             (默认: docs)
  -m, --model     OpenAI 模型          (默认: gpt-4o)
      --config    配置文件路径          (默认: .globalizer.yaml)
      --source    源语言，留空自动检测
      --dry-run   预览模式，不写文件
      --mock      测试模式（无需 API Key）
```

---

## 常见问题

**Q: 翻译结果显示「命中翻译缓存」但没有调用 API？**
A: SQLite 缓存中已有相同内容的翻译。删除 `~/.globalizer/globalizer.db` 清除缓存。

**Q: 提示 `insufficient_quota`？**
A: OpenAI API 配额不足，需充值或等待重置。

**Q: 提示 `model_not_found`？**
A: 你的 OpenAI 项目没有该模型的访问权限，去 platform.openai.com 开通。

**Q: 翻译后 Markdown 格式乱了？**
A: 提交 Issue 并附上原始文件和翻译结果，我们会修复 goldmark 解析规则。
