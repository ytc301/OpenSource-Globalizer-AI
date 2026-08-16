package github

import (
	"context"
	"fmt"
	"net/http"

	gogithub "github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

// RealClient 基于 go-github 的真实 GitHub 客户端。
type RealClient struct {
	client *gogithub.Client
}

// NewRealClient 创建真实 GitHub 客户端。
// token 为空时使用匿名访问（GitHub API 有更低的限流，仅适合只读操作）。
func NewRealClient(token string) *RealClient {
	var hc *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		hc = oauth2.NewClient(context.Background(), ts)
	}
	return &RealClient{client: gogithub.NewClient(hc)}
}

// CreatePR 通过 Git Data API 创建 Pull Request。
// 流程: 获取 base SHA → 创建 head 分支 → 上传 blob/tree/commit → 更新 ref → 创建 PR。
func (c *RealClient) CreatePR(ctx context.Context, opts CreatePROptions) (*PullRequest, error) {
	base := opts.BaseBranch
	if base == "" {
		base = "main"
	}

	// 1. 获取 base 分支的 HEAD SHA
	baseRef, _, err := c.client.Git.GetRef(ctx, opts.Owner, opts.Repo, "heads/"+base)
	if err != nil {
		return nil, fmt.Errorf("get base branch %q: %w", base, err)
	}
	baseSHA := baseRef.GetObject().GetSHA()
	if baseSHA == "" {
		return nil, fmt.Errorf("base branch %q has no commit", base)
	}

	// 2. 创建 head 分支（指向 base SHA）
	headRef := &gogithub.Reference{
		Ref:    gogithub.Ptr("refs/heads/" + opts.HeadBranch),
		Object: &gogithub.GitObject{SHA: gogithub.Ptr(baseSHA)},
	}
	if _, _, err := c.client.Git.CreateRef(ctx, opts.Owner, opts.Repo, headRef); err != nil {
		return nil, fmt.Errorf("create head branch %q: %w", opts.HeadBranch, err)
	}

	// 3. 上传每个文件为 blob
	entries := make([]*gogithub.TreeEntry, 0, len(opts.Files))
	for _, f := range opts.Files {
		blob, _, err := c.client.Git.CreateBlob(ctx, opts.Owner, opts.Repo, &gogithub.Blob{
			Content:  gogithub.Ptr(f.Content),
			Encoding: gogithub.Ptr("utf-8"),
		})
		if err != nil {
			return nil, fmt.Errorf("create blob %q: %w", f.Path, err)
		}
		entries = append(entries, &gogithub.TreeEntry{
			Path: gogithub.Ptr(f.Path),
			Mode: gogithub.Ptr("100644"),
			Type: gogithub.Ptr("blob"),
			SHA:  blob.SHA,
		})
	}

	// 4. 基于 base 创建 tree
	tree, _, err := c.client.Git.CreateTree(ctx, opts.Owner, opts.Repo, baseSHA, entries)
	if err != nil {
		return nil, fmt.Errorf("create tree: %w", err)
	}

	// 5. 创建 commit
	commit, _, err := c.client.Git.CreateCommit(ctx, opts.Owner, opts.Repo, &gogithub.Commit{
		Message: gogithub.Ptr(opts.Title),
		Tree:    tree,
		Parents: []*gogithub.Commit{{SHA: gogithub.Ptr(baseSHA)}},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("create commit: %w", err)
	}

	// 6. 更新 head 分支指向新 commit
	headRef.Object.SHA = commit.SHA
	if _, _, err := c.client.Git.UpdateRef(ctx, opts.Owner, opts.Repo, headRef, false); err != nil {
		return nil, fmt.Errorf("update head branch %q: %w", opts.HeadBranch, err)
	}

	// 7. 创建 PR
	pr, _, err := c.client.PullRequests.Create(ctx, opts.Owner, opts.Repo, &gogithub.NewPullRequest{
		Title: gogithub.Ptr(opts.Title),
		Body:  gogithub.Ptr(opts.Body),
		Head:  gogithub.Ptr(opts.HeadBranch),
		Base:  gogithub.Ptr(base),
	})
	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	return &PullRequest{Number: pr.GetNumber(), URL: pr.GetHTMLURL()}, nil
}

// GetFile 获取仓库文件内容。
func (c *RealClient) GetFile(ctx context.Context, owner, repo, path, ref string) (string, error) {
	opts := &gogithub.RepositoryContentGetOptions{Ref: ref}
	content, _, _, err := c.client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return "", fmt.Errorf("get file %q: %w", path, err)
	}
	if content == nil {
		return "", fmt.Errorf("file %q not found", path)
	}
	decoded, err := content.GetContent()
	if err != nil {
		return "", fmt.Errorf("decode file %q: %w", path, err)
	}
	return decoded, nil
}

// CreateIssueComment 在 Issue 上发布评论。
func (c *RealClient) CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) error {
	comment := &gogithub.IssueComment{Body: gogithub.Ptr(body)}
	_, _, err := c.client.Issues.CreateComment(ctx, owner, repo, number, comment)
	if err != nil {
		return fmt.Errorf("create comment on issue #%d: %w", number, err)
	}
	return nil
}

// AddIssueLabels 为 Issue 添加标签。
func (c *RealClient) AddIssueLabels(ctx context.Context, owner, repo string, number int, labels []string) error {
	if len(labels) == 0 {
		return nil
	}
	if _, _, err := c.client.Issues.AddLabelsToIssue(ctx, owner, repo, number, labels); err != nil {
		return fmt.Errorf("add labels to issue #%d: %w", number, err)
	}
	return nil
}
