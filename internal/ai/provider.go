package ai

import "context"

// Provider 定义了 AI 服务的抽象接口。
// 当前基于 OpenAI API 实现，接口设计支持未来扩展其他后端。
type Provider interface {
	// Translate 翻译一段文本到目标语言。
	// 输入可以是纯文本或 Markdown 片段。
	Translate(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error)

	// DetectLanguage 检测输入文本的语言代码。
	DetectLanguage(ctx context.Context, text string) (string, error)

	// ClassifyIssue 对 Issue 内容进行分类，返回类型标签。
	ClassifyIssue(ctx context.Context, title, body string) (*IssueClassifyResult, error)
}

// TranslateOptions 翻译请求的配置。
type TranslateOptions struct {
	SourceLang string   // 源语言代码 (空字符串表示自动检测)
	TargetLang string   // 目标语言代码
	Model      string   // 模型名称, e.g. "gpt-4o", "gpt-4o-mini"
	Preserve   []string // 需要保留不翻译的规则: "code_blocks", "links", "badges", "html"
}

// TranslateResult 翻译结果。
type TranslateResult struct {
	Translated string // 翻译后的文本
	SourceLang string // 检测到的源语言
	TokensUsed int    // 消耗的 Token 数
}

// IssueClassifyResult Issue 分类结果。
type IssueClassifyResult struct {
	Language    string   // 检测到的语言
	Type        string   // bug / feature / question / documentation
	Summary     string   // 一句话摘要 (英文)
	Confidence  float64  // 分类置信度 0-1
}
