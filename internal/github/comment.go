package github

import (
	"fmt"
	"strings"
)

// FormatIssueComment 生成 AI 翻译评论内容。
// 格式固定为:
//
//	## 🌐 AI Translation
//
//	**语言:** {lang}
//
//	**摘要:** {summary}
func FormatIssueComment(lang, summary string) string {
	var sb strings.Builder
	sb.WriteString("## 🌐 AI Translation\n\n")
	sb.WriteString(fmt.Sprintf("**语言:** %s\n\n", lang))
	sb.WriteString(fmt.Sprintf("**摘要:** %s", summary))
	return sb.String()
}

// BuildIssueLabels 根据语言和类型生成 Issue 标签。
// 返回形如 ["lang:zh-CN", "type:bug"] 的标签列表，空字段跳过。
func BuildIssueLabels(lang, issueType string) []string {
	labels := make([]string, 0, 2)
	if lang != "" {
		labels = append(labels, "lang:"+lang)
	}
	if issueType != "" {
		labels = append(labels, "type:"+issueType)
	}
	return labels
}

// NormalizeLanguage 将语言代码归一化为项目使用的地区代码。
// AI 检测返回 ISO 639-1（如 "zh"、"pt"），标签需要地区码（如 "zh-CN"、"pt-BR"）。
func NormalizeLanguage(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "zh", "zh-cn", "zh-hans", "zh-sg":
		return "zh-CN"
	case "pt", "pt-br":
		return "pt-BR"
	default:
		return strings.ToLower(strings.TrimSpace(lang))
	}
}
