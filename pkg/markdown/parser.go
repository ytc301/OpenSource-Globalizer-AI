package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Parser 基于 goldmark 的 Markdown 解析器。
type Parser struct {
	preserveCodeBlocks bool
	preserveLinks      bool
	preserveBadges     bool
	preserveHTML       bool
}

// ParserOption 解析器配置选项。
type ParserOption func(*Parser)

// NewParser 创建解析器实例。
func NewParser(opts ...ParserOption) *Parser {
	p := &Parser{
		preserveCodeBlocks: true,
		preserveLinks:      true,
		preserveBadges:     true,
		preserveHTML:       true,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithoutCodeBlocks 关闭代码块保留。
func WithoutCodeBlocks() ParserOption {
	return func(p *Parser) { p.preserveCodeBlocks = false }
}

// Segment Markdown 片段。
type Segment struct {
	Type    SegmentType
	Content string
	Lang    string
}

// SegmentType 片段类型。
type SegmentType int

const (
	Text      SegmentType = iota
	CodeBlock
	Heading
	Link
	Image
	HTMLBlock
	List
	Blockquote
	ThematicBreak
	Table
)

// CanTranslate 判断是否需要翻译。
func (p *Parser) CanTranslate(s Segment) bool {
	switch s.Type {
	case Text, Heading, List, Blockquote, Table:
		return true
	default:
		return false
	}
}

// ShouldPreserve 判断是否应原样保留。
func (p *Parser) ShouldPreserve(s Segment) bool {
	return !p.CanTranslate(s)
}

// Parse 解析 Markdown → 片段列表。
// 启用 GFM 扩展以支持表格、删除线等 GitHub 方言语法。
func (p *Parser) Parse(content string) []Segment {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	reader := text.NewReader([]byte(content))
	doc := md.Parser().Parse(reader)

	segments := make([]Segment, 0)
	p.walkNode(doc, []byte(content), &segments)
	return p.mergeAdjacentText(segments)
}

// walkNode 递归遍历 AST 节点。
func (p *Parser) walkNode(n ast.Node, source []byte, segments *[]Segment) {
	if n == nil {
		return
	}

	switch n.Kind() {
	case ast.KindDocument:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			p.walkNode(child, source, segments)
		}

	case ast.KindHeading:
		heading := n.(*ast.Heading)
		prefix := strings.Repeat("#", heading.Level) + " "
		// 保留标题内的链接 / 行内代码 / 粗体等 Markdown 语法
		text := p.nodeText(n, source)
		if text == "" {
			text = p.extractText(n, source)
		}
		*segments = append(*segments, Segment{
			Type:    Heading,
			Content: prefix + text,
		})

	case ast.KindFencedCodeBlock, ast.KindCodeBlock:
		if !p.preserveCodeBlocks {
			// 用户关闭代码块保留 → 作为普通文本参与翻译
			if text := strings.TrimSpace(p.nodeText(n, source)); text != "" {
				*segments = append(*segments, Segment{
					Type:    Text,
					Content: text,
				})
			}
			return
		}
		lines := n.Lines()
		var buf bytes.Buffer
		lang := ""
		if fb, ok := n.(*ast.FencedCodeBlock); ok {
			lang = string(fb.Language(source))
			buf.WriteString("```")
			buf.WriteString(lang)
			buf.WriteString("\n")
		}
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			buf.Write(line.Value(source))
		}
		if _, ok := n.(*ast.FencedCodeBlock); ok {
			buf.WriteString("```\n")
		}
		*segments = append(*segments, Segment{
			Type:    CodeBlock,
			Content: buf.String(),
			Lang:    lang,
		})

	case ast.KindCodeSpan:
		*segments = append(*segments, Segment{
			Type:    CodeBlock,
			Content: p.nodeText(n, source),
		})

	case ast.KindLink:
		if p.preserveLinks {
			*segments = append(*segments, Segment{
				Type:    Link,
				Content: p.nodeText(n, source),
			})
		} else {
			*segments = append(*segments, Segment{
				Type:    Text,
				Content: p.extractText(n, source),
			})
		}

	case ast.KindImage:
		*segments = append(*segments, Segment{
			Type:    Image,
			Content: p.nodeText(n, source),
		})

	case ast.KindHTMLBlock, ast.KindRawHTML:
		if p.preserveHTML {
			*segments = append(*segments, Segment{
				Type:    HTMLBlock,
				Content: p.nodeText(n, source),
			})
		}

	case ast.KindLinkReferenceDefinition:
		// 引用式链接定义（[id]: url）原样保留，否则翻译后引用失效
		if text := strings.TrimSpace(p.nodeText(n, source)); text != "" {
			*segments = append(*segments, Segment{
				Type:    HTMLBlock,
				Content: text,
			})
		}

	case ast.KindList:
		// 用原始源码保留列表标记（- / * / 1.）和换行结构
		text := strings.TrimSpace(blockSource(n, source))
		if text != "" {
			*segments = append(*segments, Segment{
				Type:    List,
				Content: text,
			})
		}

	case extast.KindTable:
		// 用原始源码保留表格的 | 分隔符和换行
		text := strings.TrimSpace(blockSource(n, source))
		if text != "" {
			*segments = append(*segments, Segment{
				Type:    Table,
				Content: text,
			})
		}

	case ast.KindBlockquote:
		// 保留引用块内的链接 / 粗体等 Markdown 语法
		text := strings.TrimSpace(blockSource(n, source))
		if text == "" {
			text = p.extractText(n, source)
		}
		if text != "" {
			*segments = append(*segments, Segment{
				Type:    Blockquote,
				Content: text,
			})
		}

	case ast.KindThematicBreak:
		*segments = append(*segments, Segment{
			Type:    ThematicBreak,
			Content: "---",
		})

	case ast.KindParagraph, ast.KindTextBlock:
		// 用原始 Markdown 源码：保留段落内的链接 / 粗体 / 行内代码，
		// 并据此检测 Badge（extractText 只取纯文本，会把 URL 和语法丢掉）
		rawText := strings.TrimSpace(p.nodeText(n, source))
		if rawText != "" {
			if p.preserveBadges && isBadge(rawText) {
				*segments = append(*segments, Segment{Type: Image, Content: rawText})
			} else {
				*segments = append(*segments, Segment{Type: Text, Content: rawText})
			}
		}

	default:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			p.walkNode(child, source, segments)
		}
	}
}

// extractText 提取节点内纯文本。
func (p *Parser) extractText(n ast.Node, source []byte) string {
	var parts []string
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if child.Kind() == ast.KindText || child.Kind() == ast.KindString {
			parts = append(parts, string(child.Text(source)))
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(parts, "")
}

// blockSource 提取块节点的原始 Markdown 源码。
// goldmark 容器节点（List/Table 等）没有 Lines()，需要遍历后代块节点
// 收集最小起始 / 最大结束位置，再从源码切片（保留列表标记、表格竖线、换行）。
func blockSource(n ast.Node, source []byte) string {
	start, end := -1, -1
	ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		// inline 节点（Text/Span/CodeSpan/Link/Image 等）没有 Lines()，调用会 panic
		switch c.Kind() {
		case ast.KindText, ast.KindString, ast.KindCodeSpan, ast.KindEmphasis,
			ast.KindLink, ast.KindAutoLink, ast.KindImage, ast.KindRawHTML:
			return ast.WalkContinue, nil
		}
		lines := c.Lines()
		if lines == nil || lines.Len() == 0 {
			return ast.WalkContinue, nil
		}
		first := lines.At(0).Start
		last := lines.At(lines.Len() - 1).Stop
		if start == -1 || first < start {
			start = first
		}
		if last > end {
			end = last
		}
		return ast.WalkContinue, nil
	})
	if start == -1 {
		return ""
	}
	// 回溯到行首，包含列表标记（- / * / 1.）或表格开头的 |
	for start > 0 && source[start-1] != '\n' {
		start--
	}
	return string(source[start:end])
}

// nodeText 获取节点原始文本（含标记）。
func (p *Parser) nodeText(n ast.Node, source []byte) string {
	lines := n.Lines()
	if lines == nil || lines.Len() == 0 {
		return p.extractText(n, source)
	}
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		buf.Write(seg.Value(source))
	}
	return buf.String()
}

// mergeAdjacentText 合并相邻同类型文本。
func (p *Parser) mergeAdjacentText(segments []Segment) []Segment {
	if len(segments) <= 1 {
		return segments
	}
	result := []Segment{segments[0]}
	for i := 1; i < len(segments); i++ {
		prev := &result[len(result)-1]
		curr := segments[i]
		if prev.Type == curr.Type && prev.Type == Text {
			prev.Content += "\n" + curr.Content
		} else {
			result = append(result, curr)
		}
	}
	return result
}

// Reassemble 重组 Markdown。
func Reassemble(segments []Segment, translations map[int]string) string {
	var result []string
	transIdx := 0
	for _, seg := range segments {
		if seg.Type == Text || seg.Type == Heading || seg.Type == List || seg.Type == Blockquote || seg.Type == Table {
			if t, ok := translations[transIdx]; ok {
				result = append(result, t)
			} else {
				result = append(result, seg.Content)
			}
			transIdx++
		} else {
			result = append(result, seg.Content)
		}
	}
	return strings.Join(result, "\n\n")
}

func isBadge(line string) bool {
	return strings.Contains(line, "[![") && strings.Contains(line, "](https://img.shields.io/")
}

var _ = goldmark.New
var _ = fmt.Sprintf
