package translator

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
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

		// 重组 Markdown — translationMap 的 key 与 Reassemble 内部计数器对齐。
		// 翻译为空的片段（分隔符丢失导致无法定位）跳过，保持原文。
		translationMap := make(map[int]string)
		reIdx := 0
		for _, seg := range segments {
			if s.parser.CanTranslate(seg) {
				if reIdx < len(translatedParts) && strings.TrimSpace(translatedParts[reIdx]) != "" {
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
// 使用带编号的分隔符 <<<SEGMENT_N>>>，比无名分隔符更鲁棒：
// 即使模型吞掉部分标记，仍能按编号定位各段，避免内容错位。
const segmentSeparatorPrefix = "<<<SEGMENT_"

// segmentSeparatorRe 匹配带编号的分隔符。
var segmentSeparatorRe = regexp.MustCompile(`<<<SEGMENT_(\d+)>>>`)

func joinForTranslation(parts []string) string {
	var sb strings.Builder
	for i, part := range parts {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(fmt.Sprintf("%s%d>>>\n", segmentSeparatorPrefix, i+1))
		sb.WriteString(part)
	}
	return sb.String()
}

// splitTranslation 将翻译结果按编号分隔符拆分回独立片段。
// 返回固定长度 count 的切片；某个片段的分隔符丢失或编号越界时，
// 对应位置留空字符串，由调用方回退为原文（不会错位）。
func splitTranslation(translated string, count int) []string {
	result := make([]string, count)

	if count <= 1 {
		// 单片段无需分隔符
		result[0] = strings.TrimSpace(translated)
		return result
	}

	loc := segmentSeparatorRe.FindAllStringSubmatchIndex(translated, -1)
	for i, m := range loc {
		idx, err := strconv.Atoi(translated[m[2]:m[3]])
		if err != nil || idx < 1 || idx > count {
			continue
		}
		// 片段内容 = 本标记结束 → 下一标记开始（或结尾）
		start := m[1]
		end := len(translated)
		if i+1 < len(loc) {
			end = loc[i+1][0]
		}
		result[idx-1] = strings.TrimSpace(translated[start:end])
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
