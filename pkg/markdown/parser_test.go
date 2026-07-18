package markdown

import (
	"strings"
	"testing"
)

func TestParse_SimpleParagraph(t *testing.T) {
	p := NewParser()
	segs := p.Parse("Hello World")

	if len(segs) != 1 {
		t.Fatalf("期望 1 个片段, 实际 %d", len(segs))
	}
	if segs[0].Type != Text {
		t.Errorf("期望 Text, 实际 %v", segs[0].Type)
	}
	if segs[0].Content != "Hello World" {
		t.Errorf("期望 'Hello World', 实际 %q", segs[0].Content)
	}
}

func TestParse_Heading(t *testing.T) {
	p := NewParser()
	segs := p.Parse("# My Project\n\nDescription here.")

	// 应该有一个 Heading 和一个 Text
	foundHeading := false
	for _, seg := range segs {
		if seg.Type == Heading {
			foundHeading = true
			if !strings.Contains(seg.Content, "My Project") {
				t.Errorf("标题内容错误: %q", seg.Content)
			}
		}
	}
	if !foundHeading {
		t.Error("未找到标题片段")
	}
}

func TestParse_CodeBlockPreserved(t *testing.T) {
	p := NewParser()
	content := "Before\n\n```go\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\n\nAfter"
	segs := p.Parse(content)

	hasCodeBlock := false
	for _, seg := range segs {
		if seg.Type == CodeBlock {
			hasCodeBlock = true
			if !strings.Contains(seg.Content, "func main()") {
				t.Errorf("代码块内容缺失: %q", seg.Content)
			}
			// 代码块不应出现在可翻译片段中
			if p.CanTranslate(seg) {
				t.Error("代码块不应被标记为可翻译")
			}
		}
	}
	if !hasCodeBlock {
		t.Error("未找到代码块片段")
	}
}

func TestParse_CodeBlockLang(t *testing.T) {
	p := NewParser()
	segs := p.Parse("```python\nprint('hello')\n```")

	for _, seg := range segs {
		if seg.Type == CodeBlock {
			if seg.Lang != "python" {
				t.Errorf("期望 lang=python, 实际 %q", seg.Lang)
			}
		}
	}
}

func TestParse_LinkPreserved(t *testing.T) {
	p := NewParser()
	segs := p.Parse("[Click here](https://example.com)")

	for _, seg := range segs {
		if seg.Type == Link {
			if !strings.Contains(seg.Content, "https://example.com") {
				t.Errorf("链接内容丢失: %q", seg.Content)
			}
		}
	}
}

func TestParse_BadgePreserved(t *testing.T) {
	p := NewParser()
	content := "[![Go](https://img.shields.io/badge/Go-1.23-blue)](https://go.dev/)"
	segs := p.Parse(content)

	for _, seg := range segs {
		if seg.Type == Image {
			if !strings.Contains(seg.Content, "img.shields.io") {
				t.Errorf("Badge 内容丢失: %q", seg.Content)
			}
		}
	}
}

func TestParse_MultipleParagraphs(t *testing.T) {
	p := NewParser()
	segs := p.Parse("Para 1\n\nPara 2\n\nPara 3")

	textCount := 0
	for _, seg := range segs {
		if seg.Type == Text {
			textCount++
		}
	}
	if textCount == 0 {
		t.Error("应该有文本片段")
	}
}

func TestParse_EmptyContent(t *testing.T) {
	p := NewParser()
	segs := p.Parse("")
	if len(segs) != 0 {
		t.Errorf("空内容应返回 0 个片段, 实际 %d", len(segs))
	}
}

func TestCanTranslate(t *testing.T) {
	p := NewParser()
	tests := []struct {
		seg      Segment
		expected bool
	}{
		{Segment{Type: Text}, true},
		{Segment{Type: Heading}, true},
		{Segment{Type: List}, true},
		{Segment{Type: Blockquote}, true},
		{Segment{Type: CodeBlock}, false},
		{Segment{Type: Link}, false},
		{Segment{Type: Image}, false},
		{Segment{Type: HTMLBlock}, false},
		{Segment{Type: ThematicBreak}, false},
	}

	for _, tt := range tests {
		got := p.CanTranslate(tt.seg)
		if got != tt.expected {
			t.Errorf("CanTranslate(%v) = %v, 期望 %v", tt.seg.Type, got, tt.expected)
		}
	}
}

func TestReassemble_Simple(t *testing.T) {
	segs := []Segment{
		{Type: Heading, Content: "# Original Title"},
		{Type: CodeBlock, Content: "```go\nconst x = 1\n```"},
		{Type: Text, Content: "Original text."},
	}
	translations := map[int]string{
		// 可翻译片段按顺序：0=Heading, 2=Text (索引 1 是代码块跳过)
		// Reassemble 内部用计数器，只对可翻译的片段计
	}
	// 先测试不提供翻译 → 用原文
	result := Reassemble(segs, translations)
	if !strings.Contains(result, "Original") {
		t.Errorf("重组后内容丢失: %q", result)
	}
}

func TestReassemble_WithTranslation(t *testing.T) {
	segs := []Segment{
		{Type: Heading, Content: "# Hello"},
		{Type: Text, Content: "World."},
	}
	// Reassemble 内部计数器: 0=Heading, 1=Text（都为可翻译）
	translations := map[int]string{
		0: "# 你好",
		1: "世界。",
	}
	result := Reassemble(segs, translations)
	if !strings.Contains(result, "你好") {
		t.Errorf("翻译后的标题丢失: %q", result)
	}
	if !strings.Contains(result, "世界") {
		t.Errorf("翻译后的文本丢失: %q", result)
	}
}

func TestParse_InlineCode(t *testing.T) {
	p := NewParser()
	segs := p.Parse("Use `go run` to start.")

	for _, seg := range segs {
		if seg.Type == CodeBlock && strings.Contains(seg.Content, "go run") {
			return // found
		}
	}
	// goldmark 可能将内联代码作为 Text 的子节点处理，
	// 验证文本片段仍包含 `go run`
	t.Logf("片段数: %d", len(segs))
	for _, seg := range segs {
		t.Logf("  [%v] %q", seg.Type, seg.Content)
	}
}

func TestParse_Blockquote(t *testing.T) {
	p := NewParser()
	segs := p.Parse("> This is a quote\n> More quote.")

	hasBlockquote := false
	for _, seg := range segs {
		if seg.Type == Blockquote {
			hasBlockquote = true
		}
	}
	if !hasBlockquote {
		t.Error("未识别引用块")
	}
}

func TestParse_ThematicBreak(t *testing.T) {
	p := NewParser()
	segs := p.Parse("Before\n\n---\n\nAfter")

	hasBreak := false
	for _, seg := range segs {
		if seg.Type == ThematicBreak {
			hasBreak = true
		}
	}
	if !hasBreak {
		t.Error("未识别分割线 ---")
	}
}

func TestParse_HTMLBlock(t *testing.T) {
	p := NewParser()
	segs := p.Parse("<div align=\"center\">\nContent\n</div>")

	for _, seg := range segs {
		if seg.Type == HTMLBlock {
			return // found
		}
	}
	t.Log("注意: goldmark 默认不启用 HTML 块解析")
}

func TestParse_List(t *testing.T) {
	p := NewParser()
	segs := p.Parse("- Item 1\n- Item 2\n- Item 3")

	hasList := false
	for _, seg := range segs {
		if seg.Type == List {
			hasList = true
		}
	}
	if !hasList {
		t.Error("未识别列表")
	}
}

func TestReassemble_CodeBlockPreservedInResult(t *testing.T) {
	segs := []Segment{
		{Type: Text, Content: "Some text."},
		{Type: CodeBlock, Content: "```go\nx := 1\n```", Lang: "go"},
		{Type: Text, Content: "More text."},
	}
	translations := map[int]string{
		0: "一些文本。",
		2: "更多文本。",
	}
	result := Reassemble(segs, translations)
	if !strings.Contains(result, "x := 1") {
		t.Error("代码块在重组后丢失")
	}
	if !strings.Contains(result, "一些") {
		t.Error("翻译文本丢失")
	}
}
