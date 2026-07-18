package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/internal/translator"
	"go.uber.org/zap"
)

func newTestLogger() *zap.Logger {
	return zap.NewNop()
}

func newTestHandler() *Handler {
	provider := ai.NewMockProvider()
	svc := translator.NewService(provider, nil, newTestLogger())
	return NewHandler(svc, newTestLogger())
}

func TestHealthEndpoint(t *testing.T) {
	h := newTestHandler()
	router := h.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码: 期望 200, 实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Errorf("status: 期望 ok, 实际 %v", resp["status"])
	}
}

func TestLanguagesEndpoint(t *testing.T) {
	h := newTestHandler()
	router := h.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("状态码: 期望 200, 实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	success, _ := resp["success"].(bool)
	if !success {
		t.Error("success 应为 true")
	}

	langs, ok := resp["languages"].([]interface{})
	if !ok {
		t.Fatal("languages 字段缺失或格式错误")
	}
	if len(langs) != 10 {
		t.Errorf("期望 10 种语言, 实际 %d", len(langs))
	}
}

func TestTranslateEndpoint_Success(t *testing.T) {
	h := newTestHandler()
	router := h.SetupRouter()

	body := `{
		"content": "# Hello\n\nWorld.",
		"target_langs": ["zh-CN", "ja"],
		"model": "gpt-4o"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/translate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("状态码: 期望 200, 实际 %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	success, _ := resp["success"].(bool)
	if !success {
		t.Error("success 应为 true")
	}

	trans, ok := resp["translations"].(map[string]interface{})
	if !ok {
		t.Fatal("translations 字段缺失")
	}
	if _, exists := trans["zh-CN"]; !exists {
		t.Error("缺少 zh-CN 翻译")
	}
	if _, exists := trans["ja"]; !exists {
		t.Error("缺少 ja 翻译")
	}
}

func TestTranslateEndpoint_MissingContent(t *testing.T) {
	h := newTestHandler()
	router := h.SetupRouter()

	body := `{"target_langs": ["zh-CN"]}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/translate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("缺少 content 应返回 400, 实际 %d", w.Code)
	}
}

func TestTranslateEndpoint_MissingTargetLang(t *testing.T) {
	h := newTestHandler()
	router := h.SetupRouter()

	body := `{"content": "# Hello"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/translate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("缺少 target_langs 应返回 400, 实际 %d", w.Code)
	}
}
