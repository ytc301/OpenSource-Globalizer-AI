package ai

import (
	"context"
	"testing"
)

func TestMockProvider_Translate(t *testing.T) {
	mock := NewMockProvider()

	result, err := mock.Translate(context.Background(), "Hello World", TranslateOptions{
		TargetLang: "zh-CN",
		Model:      "gpt-4o",
	})
	if err != nil {
		t.Fatalf("Mock 翻译失败: %v", err)
	}
	if result.Translated == "" {
		t.Error("翻译结果为空")
	}
	if result.SourceLang != "en" {
		t.Errorf("期望源语言 en, 实际 %s", result.SourceLang)
	}
	t.Logf("Mock 翻译: %s → %s", "Hello World", result.Translated)
}

func TestMockProvider_DetectLanguage(t *testing.T) {
	mock := NewMockProvider()

	lang, err := mock.DetectLanguage(context.Background(), "Bonjour le monde")
	if err != nil {
		t.Fatalf("语言检测失败: %v", err)
	}
	if lang != "en" {
		t.Errorf("Mock 应返回 en, 实际 %s", lang)
	}
}

func TestMockProvider_ClassifyIssue(t *testing.T) {
	mock := NewMockProvider()

	result, err := mock.ClassifyIssue(context.Background(), "Bug: crash on startup", "When I open the app it crashes immediately.")
	if err != nil {
		t.Fatalf("Issue 分类失败: %v", err)
	}
	if result.Type != "bug" {
		t.Errorf("期望 bug, 实际 %s", result.Type)
	}
	if result.Confidence < 0 || result.Confidence > 1 {
		t.Errorf("置信度超出范围: %f", result.Confidence)
	}
}

func TestOpenAIConfig_Defaults(t *testing.T) {
	cfg := OpenAIConfig{}
	provider := NewOpenAI(cfg)

	if provider == nil {
		t.Fatal("provider 不应为 nil")
	}
}

func TestMockProvider_CustomFunc(t *testing.T) {
	mock := &MockProvider{
		TranslateFn: func(ctx context.Context, input string, opts TranslateOptions) (*TranslateResult, error) {
			return &TranslateResult{
				Translated: "自定义翻译结果",
				SourceLang: "fr",
				TokensUsed: 42,
			}, nil
		},
	}

	result, err := mock.Translate(context.Background(), "test", TranslateOptions{TargetLang: "zh-CN"})
	if err != nil {
		t.Fatalf("自定义 Mock 失败: %v", err)
	}
	if result.Translated != "自定义翻译结果" {
		t.Errorf("期望自定义结果, 实际 %s", result.Translated)
	}
	if result.TokensUsed != 42 {
		t.Errorf("期望 42 tokens, 实际 %d", result.TokensUsed)
	}
}

func TestTranslateOptions_Defaults(t *testing.T) {
	opts := TranslateOptions{
		TargetLang: "ja",
	}
	if opts.Model != "" {
		t.Error("默认应不设置模型")
	}
	if opts.SourceLang != "" {
		t.Error("默认应不设置源语言")
	}
}
