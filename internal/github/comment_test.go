package github

import (
	"context"
	"strings"
	"testing"
)

func TestFormatIssueComment(t *testing.T) {
	got := FormatIssueComment("zh-CN", "安装失败在 Ubuntu 24.04 上")
	want := "## 🌐 AI Translation\n\n**语言:** zh-CN\n\n**摘要:** 安装失败在 Ubuntu 24.04 上"

	if got != want {
		t.Errorf("评论格式错误\n got: %q\nwant: %q", got, want)
	}

	// 必须包含固定标题和字段标签
	for _, part := range []string{"## 🌐 AI Translation", "**语言:**", "**摘要:**"} {
		if !strings.Contains(got, part) {
			t.Errorf("评论缺少必要部分 %q", part)
		}
	}
}

func TestBuildIssueLabels(t *testing.T) {
	labels := BuildIssueLabels("zh-CN", "bug")
	if len(labels) != 2 {
		t.Fatalf("期望 2 个标签, 实际 %d", len(labels))
	}
	if labels[0] != "lang:zh-CN" {
		t.Errorf("labels[0] = %q, want lang:zh-CN", labels[0])
	}
	if labels[1] != "type:bug" {
		t.Errorf("labels[1] = %q, want type:bug", labels[1])
	}
}

func TestBuildIssueLabels_OnlyLang(t *testing.T) {
	labels := BuildIssueLabels("ja", "")
	if len(labels) != 1 || labels[0] != "lang:ja" {
		t.Errorf("仅语言标签错误: %v", labels)
	}
}

func TestBuildIssueLabels_Empty(t *testing.T) {
	labels := BuildIssueLabels("", "")
	if len(labels) != 0 {
		t.Errorf("空输入应返回空标签, 实际 %v", labels)
	}
}

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"zh", "zh-CN"},
		{"zh-CN", "zh-CN"},
		{"ZH", "zh-CN"},
		{"pt", "pt-BR"},
		{"ja", "ja"},
		{"ko", "ko"},
		{"en", "en"},
	}
	for _, tt := range tests {
		if got := NormalizeLanguage(tt.in); got != tt.want {
			t.Errorf("NormalizeLanguage(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestMockClient_CommentAndLabels(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	if err := client.CreateIssueComment(ctx, "owner", "repo", 42, "comment"); err != nil {
		t.Errorf("Mock 发评论失败: %v", err)
	}
	if err := client.AddIssueLabels(ctx, "owner", "repo", 42, []string{"lang:zh-CN", "type:bug"}); err != nil {
		t.Errorf("Mock 加标签失败: %v", err)
	}
}

func TestMockClient_CommentCapturesArgs(t *testing.T) {
	client := NewMockClient()
	var gotBody string
	var gotNumber int
	client.CreateIssueCommentFn = func(ctx context.Context, owner, repo string, number int, body string) error {
		gotBody = body
		gotNumber = number
		return nil
	}

	client.CreateIssueComment(context.Background(), "o", "r", 7, "hello")
	if gotBody != "hello" || gotNumber != 7 {
		t.Errorf("参数捕获错误: body=%q number=%d", gotBody, gotNumber)
	}
}
