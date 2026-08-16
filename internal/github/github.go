package github

import (
	"context"
	"fmt"
)

// Client GitHub API 客户端接口。
type Client interface {
	// CreatePR 创建 Pull Request。
	CreatePR(ctx context.Context, opts CreatePROptions) (*PullRequest, error)

	// GetFile 获取仓库文件内容。
	GetFile(ctx context.Context, owner, repo, path, ref string) (string, error)

	// CreateIssueComment 在 Issue 上发布评论。
	CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) error

	// AddIssueLabels 为 Issue 添加标签。
	AddIssueLabels(ctx context.Context, owner, repo string, number int, labels []string) error
}

// CreatePROptions 创建 PR 的配置。
type CreatePROptions struct {
	Owner       string   // 仓库所有者
	Repo        string   // 仓库名
	Title       string   // PR 标题
	Body        string   // PR 描述
	HeadBranch  string   // 源分支
	BaseBranch  string   // 目标分支 (默认: main)
	Files       []PRFile // 要提交的文件
}

// PRFile PR 中的文件变更。
type PRFile struct {
	Path    string // 文件路径
	Content string // 文件内容
}

// PullRequest Pull Request 信息。
type PullRequest struct {
	Number int
	URL    string
}

// MockClient 用于测试的 GitHub 客户端。
type MockClient struct {
	CreatePRFn           func(ctx context.Context, opts CreatePROptions) (*PullRequest, error)
	GetFileFn            func(ctx context.Context, owner, repo, path, ref string) (string, error)
	CreateIssueCommentFn func(ctx context.Context, owner, repo string, number int, body string) error
	AddIssueLabelsFn     func(ctx context.Context, owner, repo string, number int, labels []string) error
}

func NewMockClient() *MockClient {
	return &MockClient{
		CreatePRFn: func(ctx context.Context, opts CreatePROptions) (*PullRequest, error) {
			return &PullRequest{Number: 1, URL: "https://github.com/owner/repo/pull/1"}, nil
		},
		GetFileFn: func(ctx context.Context, owner, repo, path, ref string) (string, error) {
			return "# Mock README\n\nThis is a mock file.", nil
		},
		CreateIssueCommentFn: func(ctx context.Context, owner, repo string, number int, body string) error {
			return nil
		},
		AddIssueLabelsFn: func(ctx context.Context, owner, repo string, number int, labels []string) error {
			return nil
		},
	}
}

func (m *MockClient) CreatePR(ctx context.Context, opts CreatePROptions) (*PullRequest, error) {
	if m.CreatePRFn != nil {
		return m.CreatePRFn(ctx, opts)
	}
	return nil, fmt.Errorf("CreatePRFn not set")
}

func (m *MockClient) GetFile(ctx context.Context, owner, repo, path, ref string) (string, error) {
	if m.GetFileFn != nil {
		return m.GetFileFn(ctx, owner, repo, path, ref)
	}
	return "", fmt.Errorf("GetFileFn not set")
}

func (m *MockClient) CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) error {
	if m.CreateIssueCommentFn != nil {
		return m.CreateIssueCommentFn(ctx, owner, repo, number, body)
	}
	return fmt.Errorf("CreateIssueCommentFn not set")
}

func (m *MockClient) AddIssueLabels(ctx context.Context, owner, repo string, number int, labels []string) error {
	if m.AddIssueLabelsFn != nil {
		return m.AddIssueLabelsFn(ctx, owner, repo, number, labels)
	}
	return fmt.Errorf("AddIssueLabelsFn not set")
}
