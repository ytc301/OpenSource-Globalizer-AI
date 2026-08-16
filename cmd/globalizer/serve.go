package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/internal/github"
	"github.com/ytc301/opensource-globalizer/internal/handler"
	"github.com/ytc301/opensource-globalizer/internal/store"
	"github.com/ytc301/opensource-globalizer/internal/translator"
	"github.com/ytc301/opensource-globalizer/pkg/config"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP API 服务",
	Long:  "启动 Gin HTTP 服务器，提供翻译 API。",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建依赖
	provider := ai.NewOpenAI(ai.OpenAIConfig{
		APIKey:  cfg.EffectiveAPIKey(),
		BaseURL: cfg.OpenAI.BaseURL,
	})

	var st *store.Store
	if cfg.DBPath != "" {
		st, err = store.New(cfg.DBPath)
		if err != nil {
			logger.Warn("无法创建存储，跳过翻译缓存", zap.Error(err))
		} else {
			defer st.Close()
		}
	}

	svc := translator.NewService(provider, st, logger)
	h := handler.NewHandler(svc, logger)
	router := h.SetupRouter()

	// 注册 Issue webhook 端点（需同时配置 GitHub token + webhook secret）
	if cfg.GitHub.Token != "" && cfg.GitHub.WebhookSecret != "" {
		ghClient := github.NewRealClient(cfg.GitHub.Token)
		wh := handler.NewWebhookHandler(svc, ghClient, cfg.GitHub.WebhookSecret, logger)
		router.POST("/webhook", wh.Handle)
		logger.Info("已启用 Issue webhook 端点", zap.String("path", "/webhook"))
	} else {
		logger.Info("未配置 GITHUB_TOKEN / webhook secret，跳过 Issue webhook 端点")
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("启动 HTTP 服务", zap.String("addr", addr))
	return router.Run(addr)
}
