package ai

import (
	"context"
	"fmt"
	"strings"
)

// MockProvider 用于测试的 AI Provider 实现。
type MockProvider struct {
	TranslateFn       func(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error)
	DetectLanguageFn  func(ctx context.Context, text string) (string, error)
	ClassifyIssueFn   func(ctx context.Context, title, body string) (*IssueClassifyResult, error)
}

// NewMockProvider 创建一个返回固定翻译结果的 Mock Provider。
func NewMockProvider() *MockProvider {
	return &MockProvider{
		TranslateFn: func(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error) {
			return &TranslateResult{
				Translated: fmt.Sprintf("[%s translation] %s", strings.ToUpper(opts.TargetLang), input),
				SourceLang: "en",
				TokensUsed: 10,
			}, nil
		},
		DetectLanguageFn: func(ctx context.Context, text string) (string, error) {
			return "en", nil
		},
		ClassifyIssueFn: func(ctx context.Context, title, body string) (*IssueClassifyResult, error) {
			return &IssueClassifyResult{
				Language:   "en",
				Type:       "bug",
				Summary:    "Mock issue classification",
				Confidence: 0.9,
			}, nil
		},
	}
}

func (m *MockProvider) Translate(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error) {
	if m.TranslateFn != nil {
		return m.TranslateFn(ctx, input, opts)
	}
	return nil, fmt.Errorf("TranslateFn not set")
}

func (m *MockProvider) DetectLanguage(ctx context.Context, text string) (string, error) {
	if m.DetectLanguageFn != nil {
		return m.DetectLanguageFn(ctx, text)
	}
	return "", fmt.Errorf("DetectLanguageFn not set")
}

func (m *MockProvider) ClassifyIssue(ctx context.Context, title, body string) (*IssueClassifyResult, error) {
	if m.ClassifyIssueFn != nil {
		return m.ClassifyIssueFn(ctx, title, body)
	}
	return nil, fmt.Errorf("ClassifyIssueFn not set")
}
