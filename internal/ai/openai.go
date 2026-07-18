package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// OpenAIConfig OpenAI Provider 配置。
type OpenAIConfig struct {
	APIKey      string
	BaseURL     string        // 默认: https://api.openai.com/v1
	HTTPTimeout time.Duration // 默认: 60s
	MaxRetries  int           // 默认: 3
}

// DefaultMaxRetries 默认重试次数。
const DefaultMaxRetries = 3

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
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = DefaultMaxRetries
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
		Translated: resp.Content,
		SourceLang: opts.SourceLang,
		TokensUsed: resp.TokensUsed,
	}, nil
}

// DetectLanguage 实现 Provider 接口的语言检测方法。
func (p *OpenAIProvider) DetectLanguage(ctx context.Context, text string) (string, error) {
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
	return resp.Content, nil
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
	if err := json.Unmarshal([]byte(resp.Content), &result); err != nil {
		return nil, fmt.Errorf("parse classify result: %w", err)
	}
	return &result, nil
}

// chatResponse OpenAI API 返回的 Chat Completion 响应。
type chatResponse struct {
	Content    string
	TokensUsed int
}

// chatComplete 发送请求，带指数退避重试。
func (p *OpenAIProvider) chatComplete(ctx context.Context, req chatCompletionRequest) (*chatResponse, error) {
	var lastErr error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		resp, retryable, err := p.doChatComplete(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if !retryable || attempt >= p.config.MaxRetries {
			break
		}
		// 指数退避: 1s, 2s, 4s
		delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, lastErr
}

// doChatComplete 执行单次 API 调用, 返回结果和是否可重试。
func (p *OpenAIProvider) doChatComplete(ctx context.Context, req chatCompletionRequest) (*chatResponse, bool, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, false, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		p.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, false, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	httpResp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, true, fmt.Errorf("http request: %w", err) // 网络错误可重试
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, true, fmt.Errorf("read response: %w", err)
	}

	// 429 (rate limit) 和 5xx (server error) → 重试
	// 401/403 (auth) 和 400 (bad request) → 不重试
	if httpResp.StatusCode >= 500 || httpResp.StatusCode == 429 {
		return nil, true, fmt.Errorf("api error (status %d, retryable): %s", httpResp.StatusCode, string(respBody))
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("api error (status %d): %s", httpResp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, false, fmt.Errorf("parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, false, fmt.Errorf("no choices in response")
	}

	return &chatResponse{
		Content:    result.Choices[0].Message.Content,
		TokensUsed: result.Usage.TotalTokens,
	}, false, nil
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

	// 分隔符保留规则 (必须最先处理)
	if hasPreserve("separators") {
		prompt += "\n- CRITICAL: Preserve the <<<SEGMENT_SEPARATOR>>> markers exactly as they appear. Do NOT translate or modify them."
		prompt += "\n- Translate each segment between separators independently."
	}
	if hasPreserve("code_blocks") {
		prompt += "\n- NEVER translate content inside code blocks (```). Keep the code fences intact."
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

	prompt += "\n\nReturn ONLY the translated content. No explanations or notes."
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
