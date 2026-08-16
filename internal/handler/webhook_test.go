package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/internal/github"
	"github.com/ytc301/opensource-globalizer/internal/translator"
)

func signPayload(t *testing.T, payload []byte, secret string) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func newTestWebhookHandler(t *testing.T) (*WebhookHandler, *github.MockClient, *ai.MockProvider) {
	t.Helper()
	provider := ai.NewMockProvider()
	svc := translator.NewService(provider, nil, newTestLogger())
	client := github.NewMockClient()
	h := NewWebhookHandler(svc, client, "test-secret", newTestLogger())
	return h, client, provider
}

// issueEventJSON 构造一个 opened 事件的 payload。
const issueEventJSON = `{
	"action": "opened",
	"issue": {
		"number": 42,
		"title": "安装失败 Ubuntu 24",
		"body": "在 Ubuntu 24.04 上 make install 崩溃",
		"state": "open",
		"user": {"login": "someuser", "id": 123, "type": "User"}
	},
	"repository": {
		"name": "demo",
		"full_name": "ytc301/demo",
		"owner": {"login": "ytc301", "id": 9103314, "type": "User"}
	}
}`

func TestWebhook_EndToEnd(t *testing.T) {
	h, client, provider := newTestWebhookHandler(t)

	// 自定义分类: 中文 Issue → 语言 zh, 类型 bug, 英文摘要
	provider.ClassifyIssueFn = func(ctx context.Context, title, body string) (*ai.IssueClassifyResult, error) {
		return &ai.IssueClassifyResult{
			Language:   "zh",
			Type:       "bug",
			Summary:    "Install fails on Ubuntu 24.04",
			Confidence: 0.9,
		}, nil
	}

	var gotComment string
	var gotLabels []string
	var gotNumber int
	client.CreateIssueCommentFn = func(ctx context.Context, owner, repo string, number int, body string) error {
		gotComment = body
		gotNumber = number
		return nil
	}
	client.AddIssueLabelsFn = func(ctx context.Context, owner, repo string, number int, labels []string) error {
		gotLabels = labels
		return nil
	}

	payload := []byte(issueEventJSON)
	router := gin.New()
	router.POST("/webhook", h.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", strings.NewReader(string(payload)))
	req.Header.Set(github.SHA256SignatureHeader, signPayload(t, payload, "test-secret"))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("状态码: 期望 200, 实际 %d, body: %s", w.Code, w.Body.String())
	}

	// 评论正确
	if gotNumber != 42 {
		t.Errorf("评论 issue 号错误: %d", gotNumber)
	}
	if !strings.Contains(gotComment, "## 🌐 AI Translation") {
		t.Errorf("评论缺少标题: %q", gotComment)
	}
	if !strings.Contains(gotComment, "Install fails on Ubuntu 24.04") {
		t.Errorf("评论缺少英文摘要: %q", gotComment)
	}

	// 标签正确 (zh → zh-CN 归一化)
	if len(gotLabels) != 2 {
		t.Fatalf("期望 2 个标签, 实际 %v", gotLabels)
	}
	if gotLabels[0] != "lang:zh-CN" {
		t.Errorf("语言标签错误: %q", gotLabels[0])
	}
	if gotLabels[1] != "type:bug" {
		t.Errorf("类型标签错误: %q", gotLabels[1])
	}
}

func TestWebhook_InvalidSignature(t *testing.T) {
	h, _, _ := newTestWebhookHandler(t)

	payload := []byte(issueEventJSON)
	router := gin.New()
	router.POST("/webhook", h.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", strings.NewReader(string(payload)))
	// 用错误 secret 计算签名
	req.Header.Set(github.SHA256SignatureHeader, signPayload(t, payload, "wrong-secret"))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("无效签名应返回 401, 实际 %d", w.Code)
	}
}

func TestWebhook_InvalidJSON(t *testing.T) {
	h, _, _ := newTestWebhookHandler(t)

	payload := []byte("not-json")
	router := gin.New()
	router.POST("/webhook", h.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", strings.NewReader(string(payload)))
	req.Header.Set(github.SHA256SignatureHeader, signPayload(t, payload, "test-secret"))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("无效 JSON 应返回 400, 实际 %d", w.Code)
	}
}

func TestProcessIssue_SkipNonIssueEvent(t *testing.T) {
	h, client, provider := newTestWebhookHandler(t)

	called := false
	provider.ClassifyIssueFn = func(ctx context.Context, title, body string) (*ai.IssueClassifyResult, error) {
		called = true
		return nil, nil
	}

	ev := &github.IssueEvent{Action: "closed"}
	if err := h.processIssue(context.Background(), ev); err != nil {
		t.Fatalf("跳过事件不应报错: %v", err)
	}
	if called {
		t.Error("closed 事件不应触发 AI 分类")
	}
	_ = client
}

func TestProcessIssue_SkipPullRequest(t *testing.T) {
	h, _, provider := newTestWebhookHandler(t)

	called := false
	provider.ClassifyIssueFn = func(ctx context.Context, title, body string) (*ai.IssueClassifyResult, error) {
		called = true
		return nil, nil
	}

	ev := &github.IssueEvent{
		Action: "opened",
		Issue:  github.Issue{PullRequest: &github.PullRequest{Number: 5}},
	}
	if err := h.processIssue(context.Background(), ev); err != nil {
		t.Fatalf("PR 事件不应报错: %v", err)
	}
	if called {
		t.Error("PR 事件不应触发 AI 分类")
	}
}

func TestProcessIssue_ClassifyErrorPropagates(t *testing.T) {
	h, _, provider := newTestWebhookHandler(t)

	provider.ClassifyIssueFn = func(ctx context.Context, title, body string) (*ai.IssueClassifyResult, error) {
		return nil, context.DeadlineExceeded
	}

	ev := &github.IssueEvent{
		Action: "opened",
		Issue:  github.Issue{Number: 1, Title: "x", Body: "y"},
		Repository: github.Repository{
			Name:  "demo",
			Owner: github.User{Login: "ytc301"},
		},
	}
	if err := h.processIssue(context.Background(), ev); err == nil {
		t.Error("分类失败应向上传播")
	}
}
