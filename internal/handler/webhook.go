package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ytc301/opensource-globalizer/internal/github"
	"github.com/ytc301/opensource-globalizer/internal/translator"
	"go.uber.org/zap"
)

// WebhookHandler 处理 GitHub Issue webhook 事件。
type WebhookHandler struct {
	svc    *translator.Service
	github github.Client
	secret string
	logger *zap.Logger
}

// NewWebhookHandler 创建 Webhook 处理器。
// secret 为 GitHub webhook 配置的密钥，用于 HMAC SHA-256 签名校验。
func NewWebhookHandler(svc *translator.Service, gh github.Client, secret string, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		svc:    svc,
		github: gh,
		secret: secret,
		logger: logger,
	}
}

// Handle 处理 GitHub webhook HTTP 请求。
// 流程: 校验签名 → 解析事件 → 语言检测 + 分类 → 评论 + 标签。
func (h *WebhookHandler) Handle(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "read body: " + err.Error()})
		return
	}

	// 1. 校验签名
	if err := github.VerifyRequest(c.Request, body, h.secret); err != nil {
		h.logger.Warn("webhook 签名校验失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	// 2. 解析事件
	ev, err := github.ParseIssueEvent(body)
	if err != nil {
		h.logger.Warn("webhook 事件解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// 3. 处理 Issue
	if err := h.processIssue(c.Request.Context(), ev); err != nil {
		h.logger.Error("处理 Issue 失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// processIssue 处理已解析的 Issue 事件（纯逻辑，可脱离 HTTP 测试）。
// 仅处理 opened/edited 且非 PR 的事件；其余静默跳过。
func (h *WebhookHandler) processIssue(ctx context.Context, ev *github.IssueEvent) error {
	if !ev.IsIssueEvent() || ev.IsPullRequest() {
		return nil
	}

	// 检测语言 + 分类 + 生成英文摘要
	result, err := h.svc.ClassifyIssue(ctx, ev.Issue.Title, ev.Issue.Body)
	if err != nil {
		return fmt.Errorf("classify issue #%d: %w", ev.Issue.Number, err)
	}

	lang := github.NormalizeLanguage(result.Language)
	owner := ev.Repository.Owner.Login
	repo := ev.Repository.Name

	// 发布翻译评论
	comment := github.FormatIssueComment(lang, result.Summary)
	if err := h.github.CreateIssueComment(ctx, owner, repo, ev.Issue.Number, comment); err != nil {
		return fmt.Errorf("create issue comment #%d: %w", ev.Issue.Number, err)
	}

	// 添加语言 + 类型标签
	labels := github.BuildIssueLabels(lang, result.Type)
	if err := h.github.AddIssueLabels(ctx, owner, repo, ev.Issue.Number, labels); err != nil {
		return fmt.Errorf("add issue labels #%d: %w", ev.Issue.Number, err)
	}

	h.logger.Info("Issue 处理完成",
		zap.Int("number", ev.Issue.Number),
		zap.String("lang", lang),
		zap.String("type", result.Type),
	)
	return nil
}
