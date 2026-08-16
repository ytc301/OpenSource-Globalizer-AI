package github

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// newTestRealClient 创建指向 mock HTTP server 的 RealClient。
func newTestRealClient(t *testing.T, token string, handler http.HandlerFunc) *RealClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := NewRealClient(token)
	baseURL, _ := url.Parse(server.URL + "/")
	c.client.BaseURL = baseURL
	c.client.UploadURL = baseURL
	return c
}

func TestRealClient_CreateIssueComment(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotBody string
	c := newTestRealClient(t, "test-token", func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	err := c.CreateIssueComment(context.Background(), "owner", "repo", 42, "hello comment")
	if err != nil {
		t.Fatalf("发评论失败: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("方法: 期望 POST, 实际 %s", gotMethod)
	}
	if gotPath != "/repos/owner/repo/issues/42/comments" {
		t.Errorf("路径错误: %s", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Errorf("认证头错误: %q", gotAuth)
	}
	if !strings.Contains(gotBody, "hello comment") {
		t.Errorf("请求体缺少评论内容: %s", gotBody)
	}
}

func TestRealClient_AddIssueLabels(t *testing.T) {
	var gotPath, gotBody string
	c := newTestRealClient(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"name":"lang:zh-CN"},{"name":"type:bug"}]`))
	})

	err := c.AddIssueLabels(context.Background(), "o", "r", 7, []string{"lang:zh-CN", "type:bug"})
	if err != nil {
		t.Fatalf("加标签失败: %v", err)
	}
	if gotPath != "/repos/o/r/issues/7/labels" {
		t.Errorf("路径错误: %s", gotPath)
	}
	if !strings.Contains(gotBody, "lang:zh-CN") || !strings.Contains(gotBody, "type:bug") {
		t.Errorf("请求体缺少标签: %s", gotBody)
	}
}

func TestRealClient_AddIssueLabels_Empty(t *testing.T) {
	called := false
	c := newTestRealClient(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	if err := c.AddIssueLabels(context.Background(), "o", "r", 7, nil); err != nil {
		t.Fatalf("空标签不应报错: %v", err)
	}
	if called {
		t.Error("空标签列表不应调用 GitHub API")
	}
}

func TestRealClient_GetFile(t *testing.T) {
	var gotPath string
	c := newTestRealClient(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		// "Hello" 的 base64
		w.Write([]byte(`{"name":"README.md","path":"README.md","content":"SGVsbG8=","encoding":"base64"}`))
	})

	content, err := c.GetFile(context.Background(), "o", "r", "README.md", "main")
	if err != nil {
		t.Fatalf("获取文件失败: %v", err)
	}
	if content != "Hello" {
		t.Errorf("文件内容: 期望 Hello, 实际 %q", content)
	}
	if !strings.Contains(gotPath, "README.md") {
		t.Errorf("路径错误: %s", gotPath)
	}
}

func TestRealClient_CreatePR(t *testing.T) {
	var calls []string
	c := newTestRealClient(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/git/ref/heads/main"):
			w.Write([]byte(`{"object":{"sha":"base-sha-123"}}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/git/refs"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"ref":"refs/heads/feature"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/git/blobs"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"sha":"blob-sha-1"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/git/trees"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"sha":"tree-sha-1"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/git/commits"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"sha":"commit-sha-1"}`))
		case r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/git/refs/heads/feature"):
			w.Write([]byte(`{"object":{"sha":"commit-sha-1"}}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/pulls"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"number":1,"html_url":"https://github.com/o/r/pull/1"}`))
		default:
			t.Errorf("未预期的请求: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})

	pr, err := c.CreatePR(context.Background(), CreatePROptions{
		Owner:      "o",
		Repo:       "r",
		Title:      "feat: test",
		Body:       "body",
		HeadBranch: "feature",
		BaseBranch: "main",
		Files:      []PRFile{{Path: "README.md", Content: "# Hi"}},
	})
	if err != nil {
		t.Fatalf("创建 PR 失败: %v", err)
	}
	if pr.Number != 1 {
		t.Errorf("PR 号错误: %d", pr.Number)
	}

	// 验证调用了全部 7 个步骤
	expected := []string{
		"GET /repos/o/r/git/ref/heads/main",
		"POST /repos/o/r/git/refs",
		"POST /repos/o/r/git/blobs",
		"POST /repos/o/r/git/trees",
		"POST /repos/o/r/git/commits",
		"PATCH /repos/o/r/git/refs/heads/feature",
		"POST /repos/o/r/pulls",
	}
	if len(calls) != len(expected) {
		t.Fatalf("API 调用数: 期望 %d, 实际 %d: %v", len(expected), len(calls), calls)
	}
	for i, want := range expected {
		if calls[i] != want {
			t.Errorf("第 %d 步调用错误: got %q, want %q", i+1, calls[i], want)
		}
	}
}
