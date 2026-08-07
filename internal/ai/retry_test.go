package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestChatComplete_RetryOnServerError 验证 5xx 错误触发重试。
func TestChatComplete_RetryOnServerError(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"boom"}`))
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}],"usage":{"total_tokens":5}}`))
	}))
	defer srv.Close()

	p := NewOpenAI(OpenAIConfig{
		APIKey:      "test",
		BaseURL:     srv.URL,
		HTTPTimeout: 30e9,
		MaxRetries:  3,
	})

	resp, err := p.chatComplete(context.Background(), chatCompletionRequest{
		Model: "test", Messages: []chatMessage{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("应重试成功: %v", err)
	}
	if resp.Content != "ok" {
		t.Errorf("响应内容错误: %q", resp.Content)
	}
	if attempts != 3 {
		t.Errorf("应重试 2 次共 3 次尝试, 实际 %d", attempts)
	}
}

// TestChatComplete_NoRetryOnAuthError 验证 401 不重试。
func TestChatComplete_NoRetryOnAuthError(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid key"}}`))
	}))
	defer srv.Close()

	p := NewOpenAI(OpenAIConfig{
		APIKey:      "bad",
		BaseURL:     srv.URL,
		HTTPTimeout: 30e9,
		MaxRetries:  3,
	})

	_, err := p.chatComplete(context.Background(), chatCompletionRequest{
		Model: "test", Messages: []chatMessage{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("应返回错误")
	}
	if attempts != 1 {
		t.Errorf("401 不应重试, 实际尝试 %d 次", attempts)
	}
	if !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("错误信息应包含 API 响应: %v", err)
	}
}

// TestChatComplete_RetryExhausted 验证重试耗尽后返回最终错误。
func TestChatComplete_RetryExhausted(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":{"message":"rate limited"}}`))
	}))
	defer srv.Close()

	p := NewOpenAI(OpenAIConfig{
		APIKey:      "test",
		BaseURL:     srv.URL,
		HTTPTimeout: 30e9,
		MaxRetries:  3,
	})

	_, err := p.chatComplete(context.Background(), chatCompletionRequest{
		Model: "test", Messages: []chatMessage{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("应返回错误")
	}
	if attempts != 4 { // 1 次初始 + 3 次重试
		t.Errorf("应尝试 4 次 (1+3 重试), 实际 %d", attempts)
	}
	if !strings.Contains(err.Error(), "rate limited") {
		t.Errorf("错误信息应包含 API 响应: %v", err)
	}
}
