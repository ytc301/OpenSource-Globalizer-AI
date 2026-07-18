package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OutputDir != "docs" {
		t.Errorf("OutputDir: 期望 docs, 实际 %s", cfg.OutputDir)
	}
	if cfg.Model != "gpt-4o" {
		t.Errorf("Model: 期望 gpt-4o, 实际 %s", cfg.Model)
	}
	if len(cfg.Languages) != 1 || cfg.Languages[0] != "zh-CN" {
		t.Errorf("Languages: 期望 [zh-CN], 实际 %v", cfg.Languages)
	}
	if cfg.OpenAI.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("BaseURL 不正确: %s", cfg.OpenAI.BaseURL)
	}
	if !cfg.Preserve.CodeBlocks || !cfg.Preserve.Links || !cfg.Preserve.Badges || !cfg.Preserve.HTMLTags {
		t.Error("Preserve 默认应全为 true")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port: 期望 8080, 实际 %d", cfg.Server.Port)
	}
	if cfg.Server.Mode != "debug" {
		t.Errorf("Server.Mode: 期望 debug, 实际 %s", cfg.Server.Mode)
	}
}

func TestEffectiveAPIKey(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.EffectiveAPIKey() != "" {
		t.Error("默认应无 API Key")
	}

	cfg.OpenAI.APIKey = "sk-test123"
	if cfg.EffectiveAPIKey() != "sk-test123" {
		t.Error("应返回配置中的 API Key")
	}
}

func TestLoad_NonexistentFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("不存在的文件不应报错: %v", err)
	}
	if cfg.OutputDir != "docs" {
		t.Error("应返回默认值")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, ".globalizer.yaml")
	yaml := `
languages:
  - zh-CN
  - ja
  - ko
output_dir: i18n
model: gpt-4o-mini
db_path: /tmp/test.db
openai:
  api_key: sk-test
  base_url: https://api.example.com
preserve:
  code_blocks: false
  links: true
server:
  port: 9090
  mode: release
`
	if err := os.WriteFile(cfgPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("写入测试配置: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if cfg.OutputDir != "i18n" {
		t.Errorf("OutputDir: 期望 i18n, 实际 %s", cfg.OutputDir)
	}
	if cfg.Model != "gpt-4o-mini" {
		t.Errorf("Model: 期望 gpt-4o-mini, 实际 %s", cfg.Model)
	}
	if cfg.OpenAI.APIKey != "sk-test" {
		t.Errorf("APIKey: 期望 sk-test, 实际 %s", cfg.OpenAI.APIKey)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Port: 期望 9090, 实际 %d", cfg.Server.Port)
	}
	if cfg.Server.Mode != "release" {
		t.Errorf("Mode: 期望 release, 实际 %s", cfg.Server.Mode)
	}
	if len(cfg.Languages) != 3 {
		t.Errorf("Languages: 期望 3, 实际 %d", len(cfg.Languages))
	}
	if !cfg.Preserve.Links {
		t.Error("Links 应保持 true")
	}
	if cfg.Preserve.CodeBlocks {
		t.Error("CodeBlocks 应为 false")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "bad.yaml")
	os.WriteFile(cfgPath, []byte("{{{bad yaml"), 0644)

	_, err := Load(cfgPath)
	if err == nil {
		t.Error("无效 YAML 应返回错误")
	}
}

func TestLoad_PartialYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "partial.yaml")
	yaml := `output_dir: partial_only
model: gpt-4o-mini`
	os.WriteFile(cfgPath, []byte(yaml), 0644)

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("部分配置加载失败: %v", err)
	}
	if cfg.OutputDir != "partial_only" {
		t.Error("应读取文件值")
	}
	// 未配置的字段应保持默认值
	if cfg.Server.Port != 8080 {
		t.Errorf("Port 应保持默认 8080, 实际 %d", cfg.Server.Port)
	}
	if !cfg.Preserve.CodeBlocks {
		t.Error("CodeBlocks 应保持默认 true")
	}
}
