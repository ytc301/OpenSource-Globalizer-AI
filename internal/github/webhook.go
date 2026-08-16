package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Webhook 签名相关错误。
var (
	// ErrMissingSecret webhook secret 未配置。
	ErrMissingSecret = errors.New("webhook secret is empty")
	// ErrMissingSignature 请求缺少签名头。
	ErrMissingSignature = errors.New("missing webhook signature header")
	// ErrInvalidSignature 签名校验不通过。
	ErrInvalidSignature = errors.New("webhook signature mismatch")
)

// GitHub Webhook 请求头名称。
const (
	// SHA256SignatureHeader HMAC SHA-256 签名头。
	SHA256SignatureHeader = "X-Hub-Signature-256"
)

// VerifySignature 校验 GitHub Webhook 的 HMAC SHA-256 签名。
// signature 期望格式为 "sha256=<hex digest>"。
func VerifySignature(payload []byte, signature, secret string) bool {
	if signature == "" || secret == "" {
		return false
	}

	const prefix = "sha256="
	if len(signature) <= len(prefix) || signature[:len(prefix)] != prefix {
		return false
	}

	provided, err := hex.DecodeString(signature[len(prefix):])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)

	// 常数时间比较，避免时序攻击。
	return hmac.Equal(provided, expected)
}

// VerifyRequest 从 HTTP 请求校验 Webhook 签名。
// 使用 X-Hub-Signature-256 头；secret 为空时返回 ErrMissingSecret。
func VerifyRequest(r *http.Request, body []byte, secret string) error {
	if secret == "" {
		return ErrMissingSecret
	}
	sig := r.Header.Get(SHA256SignatureHeader)
	if sig == "" {
		return ErrMissingSignature
	}
	if !VerifySignature(body, sig, secret) {
		return ErrInvalidSignature
	}
	return nil
}

// IssueEvent GitHub Issue webhook 事件 (issues.*)。
type IssueEvent struct {
	Action     string     `json:"action"` // opened / edited / closed / ...
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
	Sender     User       `json:"sender"`
}

// Issue GitHub Issue 实体。
type Issue struct {
	Number      int          `json:"number"`
	Title       string       `json:"title"`
	Body        string       `json:"body"`
	State       string       `json:"state"`
	User        User         `json:"user"`
	Labels      []Label      `json:"labels"`
	PullRequest *PullRequest `json:"pull_request,omitempty"` // 非 nil 表示该 issue 实为 PR
}

// User GitHub 用户。
type User struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Type  string `json:"type"`
}

// Label Issue 标签。
type Label struct {
	Name string `json:"name"`
}

// Repository 仓库信息。
type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
}

// ParseIssueEvent 解析 GitHub Issue webhook 事件 JSON。
func ParseIssueEvent(payload []byte) (*IssueEvent, error) {
	var ev IssueEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, fmt.Errorf("parse issue event: %w", err)
	}
	if ev.Action == "" {
		return nil, fmt.Errorf("parse issue event: missing action field")
	}
	return &ev, nil
}

// IsIssueEvent 判断事件 action 是否为需要处理的 Issue 事件 (opened / edited)。
func (e *IssueEvent) IsIssueEvent() bool {
	return e.Action == "opened" || e.Action == "edited"
}

// IsPullRequest 判断该 issue 是否实为 Pull Request。
func (e *IssueEvent) IsPullRequest() bool {
	return e.Issue.PullRequest != nil
}
