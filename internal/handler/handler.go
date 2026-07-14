package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ytc301/opensource-globalizer/internal/translator"
	"go.uber.org/zap"
)

// Handler HTTP API 处理器。
type Handler struct {
	translator *translator.Service
	logger     *zap.Logger
}

// NewHandler 创建 Handler 实例。
func NewHandler(svc *translator.Service, logger *zap.Logger) *Handler {
	return &Handler{
		translator: svc,
		logger:     logger,
	}
}

// TranslateRequest 翻译请求。
type TranslateRequest struct {
	Content    string   `json:"content" binding:"required"`
	TargetLang []string `json:"target_langs" binding:"required"`
	Model      string   `json:"model"`
}

// SetupRouter 配置 Gin 路由。
func (h *Handler) SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(h.loggerMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "0.1.0"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/translate", h.handleTranslate)
		v1.GET("/languages", h.handleLanguages)
	}

	return r
}

// handleTranslate 处理翻译请求。
func (h *Handler) handleTranslate(c *gin.Context) {
	var req TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.Model == "" {
		req.Model = "gpt-4o"
	}

	results, err := h.translator.TranslateFile(c.Request.Context(), req.Content, req.TargetLang, req.Model)
	if err != nil {
		h.logger.Error("翻译失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"translations": results,
	})
}

// handleLanguages 返回支持的语言列表。
func (h *Handler) handleLanguages(c *gin.Context) {
	langs := []gin.H{
		{"code": "en", "name": "English", "native_name": "English"},
		{"code": "zh-CN", "name": "Chinese (Simplified)", "native_name": "简体中文"},
		{"code": "ja", "name": "Japanese", "native_name": "日本語"},
		{"code": "ko", "name": "Korean", "native_name": "한국어"},
		{"code": "es", "name": "Spanish", "native_name": "Español"},
		{"code": "fr", "name": "French", "native_name": "Français"},
		{"code": "de", "name": "German", "native_name": "Deutsch"},
		{"code": "pt-BR", "name": "Portuguese (Brazil)", "native_name": "Português (Brasil)"},
		{"code": "ru", "name": "Russian", "native_name": "Русский"},
		{"code": "ar", "name": "Arabic", "native_name": "العربية"},
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "languages": langs})
}

// loggerMiddleware Gin 日志中间件。
func (h *Handler) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("HTTP 请求",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}
