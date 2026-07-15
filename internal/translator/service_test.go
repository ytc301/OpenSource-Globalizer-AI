package translator

import (
	"context"
	"strings"
	"testing"

	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/pkg/markdown"
	"go.uber.org/zap"
)

func newTestLogger() *zap.Logger {
	return zap.NewNop()
}

func TestTranslateFile_SingleLang(t *testing.T) {
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	content := "# Hello\n\nThis is a test.\n\n```go\nconst x = 1\n```\n\nMore text."
	results, err := svc.TranslateFile(context.Background(), content, []string{"zh-CN"}, "gpt-4o")
	if err != nil {
		t.Fatalf("翻译失败: %v", err)
	}

	translated, ok := results["zh-CN"]
	if !ok {
		t.Fatal("缺少 zh-CN 翻译结果")
	}

	// Mock 翻译应该保留分隔符标记
	if translated == "" {
		t.Error("翻译结果为空")
	}

	// 不应该破坏代码块
	t.Logf("翻译结果:\n%s", translated)
}

func TestTranslateFile_MultiLang(t *testing.T) {
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	content := "# Title\n\nSome content here."
	langs := []string{"zh-CN", "ja", "ko"}
	results, err := svc.TranslateFile(context.Background(), content, langs, "gpt-4o")
	if err != nil {
		t.Fatalf("翻译失败: %v", err)
	}

	for _, lang := range langs {
		if _, ok := results[lang]; !ok {
			t.Errorf("缺少 %s 翻译结果", lang)
		}
	}
}

func TestTranslateFile_EmptyContent(t *testing.T) {
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	_, err := svc.TranslateFile(context.Background(), "", []string{"zh-CN"}, "gpt-4o")
	if err == nil {
		t.Error("空内容应返回错误")
	}
}

func TestTranslateFile_CodeBlockPreserved(t *testing.T) {
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	codeSnippet := "```go\nfunc hello() { return \"world\" }\n```"
	content := "Before code.\n\n" + codeSnippet + "\n\nAfter code."

	results, err := svc.TranslateFile(context.Background(), content, []string{"zh-CN"}, "gpt-4o")
	if err != nil {
		t.Fatalf("翻译失败: %v", err)
	}

	translated := results["zh-CN"]
	// 代码块应该原样保留
	if !strings.Contains(translated, "func hello()") {
		t.Error("代码块内容被丢失或破坏")
	}
	if !strings.Contains(translated, "```") {
		t.Error("代码块标记丢失")
	}
}

func TestSegments_Debug(t *testing.T) {
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	content := "# Title\n\n## Subtitle\n\n```go\nx := 1\n```\n\nText here.\n\n[Link](https://x.com)"

	segs := svc.Segments(content)
	if len(segs) == 0 {
		t.Fatal("解析结果为空")
	}

	typeCounts := make(map[markdown.SegmentType]int)
	for _, seg := range segs {
		typeCounts[seg.Type]++
	}

	t.Logf("片段统计: Text=%d, Heading=%d, CodeBlock=%d, Link=%d",
		typeCounts[markdown.Text], typeCounts[markdown.Heading], typeCounts[markdown.CodeBlock], typeCounts[markdown.Link])

	if typeCounts[markdown.CodeBlock] == 0 {
		t.Error("未识别代码块")
	}
}

func TestTranslateFile_UsedSegmentsOrder(t *testing.T) {
	// 验证可翻译片段和保留片段交替时的重组正确性
	provider := ai.NewMockProvider()
	svc := NewService(provider, nil, newTestLogger())

	content := "# A\n\nText A.\n\n```go\ncode\n```\n\nText B."

	results, err := svc.TranslateFile(context.Background(), content, []string{"zh-CN"}, "gpt-4o")
	if err != nil {
		t.Fatalf("翻译失败: %v", err)
	}

	translated := results["zh-CN"]
	t.Logf("交替重组结果:\n%s", translated)

	// 重组后应该包含代码块
	if !strings.Contains(translated, "code") {
		t.Error("重组后代码块内容丢失")
	}
}
