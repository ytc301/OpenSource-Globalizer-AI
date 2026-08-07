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

func TestJoinForTranslation_NumberedSeparators(t *testing.T) {
	parts := []string{"Hello world", "Second paragraph", "Third"}
	joined := joinForTranslation(parts)

	want := "<<<SEGMENT_1>>>\nHello world\n\n<<<SEGMENT_2>>>\nSecond paragraph\n\n<<<SEGMENT_3>>>\nThird"
	if joined != want {
		t.Errorf("拼接结果不正确\n got: %q\nwant: %q", joined, want)
	}
}

func TestSplitTranslation_AllMarkersPreserved(t *testing.T) {
	translated := "<<<SEGMENT_1>>>\n你好世界\n\n<<<SEGMENT_2>>>\n第二段\n\n<<<SEGMENT_3>>>\n第三段"
	parts := splitTranslation(translated, 3)

	if len(parts) != 3 {
		t.Fatalf("段数应为 3, got %d", len(parts))
	}
	wants := []string{"你好世界", "第二段", "第三段"}
	for i, w := range wants {
		if parts[i] != w {
			t.Errorf("parts[%d] = %q, want %q", i, parts[i], w)
		}
	}
}

func TestSplitTranslation_MarkerLost(t *testing.T) {
	// SEGMENT_2 的标记被模型吞掉 → 第 2 段留空（回退原文）
	translated := "<<<SEGMENT_1>>>\n你好\n\n<<<SEGMENT_3>>>\n第三段"
	parts := splitTranslation(translated, 3)

	if parts[0] != "你好" {
		t.Errorf("parts[0] = %q, want 你好", parts[0])
	}
	if parts[1] != "" {
		t.Errorf("parts[1] 应留空(分隔符丢失), got %q", parts[1])
	}
	if parts[2] != "第三段" {
		t.Errorf("parts[2] = %q, want 第三段", parts[2])
	}
}

func TestSplitTranslation_AllMarkersLost(t *testing.T) {
	// 模型完全吞掉分隔符 → 全部留空，整体回退原文
	translated := "你好\n第二段\n第三段"
	parts := splitTranslation(translated, 3)

	for i, p := range parts {
		if p != "" {
			t.Errorf("parts[%d] 应为空(全丢失), got %q", i, p)
		}
	}
}

func TestSplitTranslation_OutOfRangeIndex(t *testing.T) {
	// 模型编错号（SEGMENT_9 超出范围）→ 忽略该段
	translated := "<<<SEGMENT_1>>>\n你好\n\n<<<SEGMENT_9>>>\n乱编号"
	parts := splitTranslation(translated, 2)

	if parts[0] != "你好" {
		t.Errorf("parts[0] = %q, want 你好", parts[0])
	}
	if parts[1] != "" {
		t.Errorf("parts[1] 应留空(编号越界), got %q", parts[1])
	}
}

func TestSplitTranslation_SinglePart(t *testing.T) {
	parts := splitTranslation("just one", 1)
	if len(parts) != 1 || parts[0] != "just one" {
		t.Errorf("单段拆分错误: %v", parts)
	}
}
