package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIConfig OpenAI Provider 配置。
type OpenAIConfig struct {
	APIKey      string
	BaseURL     string // 默认: https://api.openai.com/v1
	HTTPTimeout time.Duration
}

// OpenAIProvider 实现 Provider 接口，使用 OpenAI API。
type OpenAIProvider struct {
	config     OpenAIConfig
	httpClient *http.Client
}

// NewOpenAI 创建 OpenAI Provider 实例。
func NewOpenAI(cfg OpenAIConfig) *OpenAIProvider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.HTTPTimeout == 0 {
		cfg.HTTPTimeout = 60 * time.Second
	}
	return &OpenAIProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

// chatCompletionRequest OpenAI Chat Completion 请求体。
type chatCompletionRequest struct {
	Model    string          `json:"model"`
	Messages []chatMessage   `json:"messages"`
	MaxTokens int            `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Translate 实现 Provider 接口的翻译方法。
func (p *OpenAIProvider) Translate(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error) {
	model := opts.Model
	if model == "" {
		model = "gpt-4o"
	}

	systemPrompt := buildTranslateSystemPrompt(opts)
	userPrompt := fmt.Sprintf("Translate the following content to %s:\n\n%s", opts.TargetLang, input)

	req := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	resp, err := p.chatComplete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("openai translate: %w", err)
	}

	return &TranslateResult{
		Translated: resp,
		SourceLang: opts.SourceLang,
		TokensUsed: 0, // TODO: 从 API 响应中解析
	}, nil
}

// DetectLanguage 实现 Provider 接口的语言检测方法。
func (p *OpenAIProvider) DetectLanguage(ctx context.Context, text string) (string, error) {
	// 取前 500 字符用于检测，节省 Token
	sample := text
	if len(sample) > 500 {
		sample = sample[:500]
	}

	req := chatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []chatMessage{
			{Role: "system", Content: "You are a language detector. Respond with ONLY the ISO 639-1 language code (e.g., 'en', 'zh', 'ja', 'ko'). No other text."},
			{Role: "user", Content: fmt.Sprintf("Detect the language of this text:\n\n%s", sample)},
		},
	}

	resp, err := p.chatComplete(ctx, req)
	if err != nil {
		return "", fmt.Errorf("detect language: %w", err)
	}
	return resp, nil
}

// ClassifyIssue 实现 Provider 接口的 Issue 分类方法。
func (p *OpenAIProvider) ClassifyIssue(ctx context.Context, title, body string) (*IssueClassifyResult, error) {
	req := chatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []chatMessage{
			{Role: "system", Content: issueClassifySystemPrompt},
			{Role: "user", Content: fmt.Sprintf("Title: %s\n\nBody:\n%s", title, body)},
		},
	}

	resp, err := p.chatComplete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("classify issue: %w", err)
	}

	var result IssueClassifyResult
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, fmt.Errorf("parse classify result: %w", err)
	}
	return &result, nil
}

// chatComplete 发送 Chat Completion 请求并返回响应文本。
func (p *OpenAIProvider) chatComplete(ctx context.Context, req chatCompletionRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		p.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	httpResp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api error (status %d): %s", httpResp.StatusCode, string(respBody))
	}

	// 解析响应
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

// buildTranslateSystemPrompt 构建翻译 System Prompt。
func buildTranslateSystemPrompt(opts TranslateOptions) string {
	prompt := "You are a professional translator for open-source documentation. Translate the user's content accurately while preserving all technical formatting."

	hasPreserve := func(rule string) bool {
		for _, p := range opts.Preserve {
			if p == rule {
				return true
			}
		}
		return false
	}

	if hasPreserve("code_blocks") {
		prompt += "\n- NEVER translate content inside code blocks (```)."
	}
	if hasPreserve("links") {
		prompt += "\n- Preserve all URLs and link syntax intact."
	}
	if hasPreserve("badges") {
		prompt += "\n- Keep badge markdown ([![...]](...)) unchanged."
	}
	if hasPreserve("html") {
		prompt += "\n- Do not translate HTML tag attributes."
	}

	prompt += "\n\nReturn ONLY the translated content. No explanations."
	return prompt
}

// issueClassifySystemPrompt 用于 Issue 分类的 System Prompt。
const issueClassifySystemPrompt = `You are an issue triage assistant for open-source projects.
Analyze the issue and respond with ONLY a JSON object (no markdown, no explanation):

{
  "language": "ISO 639-1 code of the issue's language",
  "type": "bug|feature|question|documentation",
  "summary": "One-sentence English summary of the issue",
  "confidence": 0.0-1.0
}`
