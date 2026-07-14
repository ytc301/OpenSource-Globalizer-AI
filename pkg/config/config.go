package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用全局配置。
type Config struct {
	Languages []string       `mapstructure:"languages"`
	OutputDir string         `mapstructure:"output_dir"`
	Model     string         `mapstructure:"model"`
	DBPath    string         `mapstructure:"db_path"`
	OpenAI    OpenAIConfig   `mapstructure:"openai"`
	Preserve  PreserveConfig `mapstructure:"preserve"`
	Server    ServerConfig   `mapstructure:"server"`
}

// OpenAIConfig OpenAI API 配置。
type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// PreserveConfig 翻译保留规则。
type PreserveConfig struct {
	CodeBlocks bool `mapstructure:"code_blocks"`
	Links      bool `mapstructure:"links"`
	Badges     bool `mapstructure:"badges"`
	HTMLTags   bool `mapstructure:"html_tags"`
}

// ServerConfig HTTP 服务配置。
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DefaultConfig 返回默认配置。
func DefaultConfig() *Config {
	return &Config{
		Languages: []string{"zh-CN"},
		OutputDir: "docs",
		Model:     "gpt-4o",
		DBPath:    "~/.globalizer/globalizer.db",
		OpenAI: OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
		},
		Preserve: PreserveConfig{
			CodeBlocks: true,
			Links:      true,
			Badges:     true,
			HTMLTags:   true,
		},
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
	}
}

// Load 使用 viper 加载配置。
// 优先级: 环境变量 > 配置文件 > 默认值
func Load(cfgFile string) (*Config, error) {
	cfg := DefaultConfig()
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".globalizer")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	v.SetEnvPrefix("GLOBALIZER")
	v.BindEnv("openai.api_key", "OPENAI_API_KEY")
	v.BindEnv("db_path", "GLOBALIZER_DB_PATH")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置: %w", err)
	}

	return cfg, nil
}

// EffectiveAPIKey 返回实际使用的 API Key（环境变量优先）。
func (c *Config) EffectiveAPIKey() string {
	if c.OpenAI.APIKey != "" {
		return c.OpenAI.APIKey
	}
	return ""
}
