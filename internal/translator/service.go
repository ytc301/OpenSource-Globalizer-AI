package translator

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/internal/store"
	"github.com/ytc301/opensource-globalizer/pkg/markdown"
	"go.uber.org/zap"
)

// Service 翻译服务。
type Service struct {
	parser   *markdown.Parser
	provider ai.Provider
	store    *store.Store
	logger   *zap.Logger
}

// NewService 创建翻译服务实例。
func NewService(provider ai.Provider, st *store.Store, logger *zap.Logger, opts ...markdown.ParserOption) *Service {
	return &Service{
		parser:   markdown.NewParser(opts...),
		provider: provider,
		store:    st,
		logger:   logger,
	}
}

// TranslateFile 翻译单个 Markdown 文件到多种目标语言。
func (s *Service) TranslateFile(ctx context.Context, content string, langs []string, model string) (map[string]string, error) {
	// 1. 解析 Markdown → 片段
	segments := s.parser.Parse(content)

	// 2. 收集需要翻译的文本
	var translatable []string
	for _, seg := range segments {
		if s.parser.CanTranslate(seg) {
			translatable = append(translatable, seg.Content)
		}
	}

	if len(translatable) == 0 {
		return nil, fmt.Errorf("文件中没有可翻译的内容")
	}

	// 3. 计算源文件的哈希（用于缓存命中）
	sourceHash := hashContent(content)

	results := make(map[string]string, len(langs))

	for _, targetLang := range langs {
		// 检查缓存
		if s.store != nil {
			if cached, _ := s.store.GetCached(sourceHash, targetLang); cached != nil {
				s.logger.Info("命中翻译缓存", zap.String("lang", targetLang))
				results[targetLang] = cached.Translated
				continue
			}
		}

		// 调用 AI 翻译（合并所有可翻译文本）
		combinedInput := joinForTranslation(translatable)
		result, err := s.provider.Translate(ctx, combinedInput, ai.TranslateOptions{
			TargetLang: targetLang,
			Model:      model,
			Preserve:   []string{"separators", "code_blocks", "links", "badges", "html"},
		})
		if err != nil {
			return nil, fmt.Errorf("翻译到 %s 失败: %w", targetLang, err)
		}

		// 拆分翻译结果
		translatedParts := splitTranslation(result.Translated, len(translatable))

		// 重组 Markdown — translationMap 的 key 与 Reassemble 内部计数器对齐
		translationMap := make(map[int]string)
		reIdx := 0
		for _, seg := range segments {
			if s.parser.CanTranslate(seg) {
				if reIdx < len(translatedParts) {
					translationMap[reIdx] = translatedParts[reIdx]
				}
				reIdx++
			}
		}

		assembled := markdown.Reassemble(segments, translationMap)
		results[targetLang] = assembled

		// 写入缓存
		if s.store != nil {
			s.store.PutCache(&store.Translation{
				SourceHash: sourceHash,
				TargetLang: targetLang,
				SourceText: content,
				Translated: assembled,
				Model:      model,
				TokensUsed: result.TokensUsed,
			})
		}

		s.logger.Info("翻译完成", zap.String("lang", targetLang), zap.Int("tokens", result.TokensUsed))
	}

	return results, nil
}

// joinForTranslation 将多个文本拼接为一个翻译请求。
// 使用特殊分隔符，在 Prompt 中明确要求 AI 保留。
const segmentSeparator = "\n\n<<<SEGMENT_SEPARATOR>>>\n\n"

func joinForTranslation(parts []string) string {
	var sb strings.Builder
	for i, part := range parts {
		if i > 0 {
			sb.WriteString(segmentSeparator)
		}
		sb.WriteString(part)
	}
	return sb.String()
}

// splitTranslation 将翻译结果按分隔符拆分回独立片段。
// 如果分隔符全部丢失，返回整个翻译作为单段（回退模式）。
func splitTranslation(translated string, count int) []string {
	if count <= 1 {
		return []string{strings.TrimSpace(translated)}
	}

	parts := strings.Split(translated, segmentSeparator)

	// 清理每段的前后空白
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimSpace(p)
	}

	// 拆分数量不匹配 → 回退：整个翻译作为第一段，后面用原文
	if len(result) != count {
		return append(result[:1], result...) // 返回实际拆分结果，让调用方处理
	}

	return result
}

// hashContent 计算内容的 SHA-256 哈希。
func hashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", h)
}

// Segments 返回解析后的片段（调试用）。
func (s *Service) Segments(content string) []markdown.Segment {
	return s.parser.Parse(content)
}
