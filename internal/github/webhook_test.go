package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

// signPayload 计算测试用的 HMAC SHA-256 签名。
func signPayload(t *testing.T, payload []byte, secret string) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifySignature_Valid(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	secret := "test-secret"
	sig := signPayload(t, payload, secret)

	if !VerifySignature(payload, sig, secret) {
		t.Error("有效签名应通过校验")
	}
}

func TestVerifySignature_Invalid(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	secret := "test-secret"
	// 用错误的 secret 计算签名
	wrongSig := signPayload(t, payload, "wrong-secret")

	if VerifySignature(payload, wrongSig, secret) {
		t.Error("无效签名不应通过校验")
	}
}

func TestVerifySignature_TamperedPayload(t *testing.T) {
	secret := "test-secret"
	sig := signPayload(t, []byte(`{"action":"opened"}`), secret)
	// 篡改 payload
	if VerifySignature([]byte(`{"action":"edited"}`), sig, secret) {
		t.Error("payload 被篡改后签名应失效")
	}
}

func TestVerifySignature_EmptySecret(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	sig := signPayload(t, payload, "test-secret")

	if VerifySignature(payload, sig, "") {
		t.Error("空 secret 不应通过校验")
	}
}

func TestVerifySignature_WrongPrefix(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	secret := "test-secret"
	// 旧版 sha1= 前缀或错误前缀应被拒绝
	if VerifySignature(payload, "sha1=deadbeef", secret) {
		t.Error("非 sha256= 前缀应被拒绝")
	}
}

func TestVerifySignature_MalformedHex(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	if VerifySignature(payload, "sha256=not-hex-zz", "test-secret") {
		t.Error("非法 hex 签名应被拒绝")
	}
}

func TestVerifyRequest_MissingSecret(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	err := VerifyRequest(req, []byte("{}"), "")
	if err != ErrMissingSecret {
		t.Errorf("期望 ErrMissingSecret, 实际 %v", err)
	}
}

func TestVerifyRequest_MissingSignature(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	err := VerifyRequest(req, []byte("{}"), "test-secret")
	if err != ErrMissingSignature {
		t.Errorf("期望 ErrMissingSignature, 实际 %v", err)
	}
}

func TestVerifyRequest_Valid(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	secret := "test-secret"
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set(SHA256SignatureHeader, signPayload(t, payload, secret))

	if err := VerifyRequest(req, payload, secret); err != nil {
		t.Errorf("有效请求应通过校验, 实际 %v", err)
	}
}

func TestVerifyRequest_Invalid(t *testing.T) {
	payload := []byte(`{"action":"opened"}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set(SHA256SignatureHeader, "sha256="+hex.EncodeToString([]byte("fake")))

	err := VerifyRequest(req, payload, "test-secret")
	if err != ErrInvalidSignature {
		t.Errorf("期望 ErrInvalidSignature, 实际 %v", err)
	}
}

func TestParseIssueEvent_Opened(t *testing.T) {
	payload := []byte(`{
		"action": "opened",
		"issue": {
			"number": 42,
			"title": "安装失败 Ubuntu 24",
			"body": "在 Ubuntu 24.04 上 make install 崩溃",
			"state": "open",
			"user": {"login": "someuser", "id": 123, "type": "User"},
			"labels": [{"name": "bug"}]
		},
		"repository": {
			"name": "demo",
			"full_name": "ytc301/demo",
			"owner": {"login": "ytc301", "id": 9103314, "type": "User"}
		},
		"sender": {"login": "someuser", "id": 123, "type": "User"}
	}`)

	ev, err := ParseIssueEvent(payload)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if ev.Action != "opened" {
		t.Errorf("期望 action=opened, 实际 %q", ev.Action)
	}
	if ev.Issue.Number != 42 {
		t.Errorf("期望 issue number=42, 实际 %d", ev.Issue.Number)
	}
	if ev.Issue.Title != "安装失败 Ubuntu 24" {
		t.Errorf("issue 标题错误: %q", ev.Issue.Title)
	}
	if ev.Repository.FullName != "ytc301/demo" {
		t.Errorf("仓库名错误: %q", ev.Repository.FullName)
	}
	if len(ev.Issue.Labels) != 1 || ev.Issue.Labels[0].Name != "bug" {
		t.Errorf("标签解析错误: %v", ev.Issue.Labels)
	}
	if !ev.IsIssueEvent() {
		t.Error("opened 应被判定为 Issue 事件")
	}
	if ev.IsPullRequest() {
		t.Error("普通 issue 不应被判定为 PR")
	}
}

func TestParseIssueEvent_Edited(t *testing.T) {
	payload := []byte(`{"action":"edited","issue":{"number":7,"title":"标题"}}`)

	ev, err := ParseIssueEvent(payload)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if ev.Action != "edited" {
		t.Errorf("期望 action=edited, 实际 %q", ev.Action)
	}
	if !ev.IsIssueEvent() {
		t.Error("edited 应被判定为 Issue 事件")
	}
}

func TestParseIssueEvent_OtherAction(t *testing.T) {
	payload := []byte(`{"action":"closed","issue":{"number":1,"title":"x"}}`)

	ev, err := ParseIssueEvent(payload)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if ev.IsIssueEvent() {
		t.Error("closed 不应被判定为需要处理的 Issue 事件")
	}
}

func TestParseIssueEvent_IsPullRequest(t *testing.T) {
	payload := []byte(`{
		"action": "opened",
		"issue": {
			"number": 5,
			"title": "PR 标题",
			"pull_request": {"url": "https://api.github.com/repos/ytc301/demo/pulls/5"}
		}
	}`)

	ev, err := ParseIssueEvent(payload)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if !ev.IsPullRequest() {
		t.Error("带 pull_request 字段的 issue 应被判定为 PR")
	}
}

func TestParseIssueEvent_MissingAction(t *testing.T) {
	_, err := ParseIssueEvent([]byte(`{"issue":{"number":1}}`))
	if err == nil {
		t.Error("缺少 action 字段应返回错误")
	}
}

func TestParseIssueEvent_InvalidJSON(t *testing.T) {
	_, err := ParseIssueEvent([]byte(`not-json`))
	if err == nil {
		t.Error("无效 JSON 应返回错误")
	}
}
