package core

import (
	"context"
	"time"
)

//go:generate mockgen -package mock -destination ../mock/scm_mock.go . SCMService,Client,RepoService,UserService,ContentService,GitService

// SCMService to interact with given SCM provider
type SCMService interface {
	Client(scm SCMProvider) (Client, error)
}

// Token for OAuth authorization
type Token struct {
	Token   string
	Refresh string
	Expires time.Time
}

// Hook defines a created webhook
type Hook struct {
	ID string
}

// Client connects to a SCM provider
type Client interface {
	Repositories() RepoService
	Users() UserService
	Git() GitService
	Contents() ContentService
	Issues() IssueService
	Token(user *User) Token
}

// RepoService provides operations with SCM
type RepoService interface {
	NewReportID(repo *Repo) string
	// List repositories from SCM context
	List(ctx context.Context, user *User) ([]*Repo, error)
	Find(ctx context.Context, user *User, name string) (*Repo, error)
	CloneURL(ctx context.Context, user *User, name string) (string, error)
	CreateHook(ctx context.Context, user *User, name string) (*Hook, error)
	RemoveHook(ctx context.Context, user *User, name string, hook *Hook) error
}

// UserService defines operations with SCM
type UserService interface {
	Find(ctx context.Context, token *Token) (*User, error)
	Create(ctx context.Context, token *Token) (*User, error)
	Bind(ctx context.Context, user *User, token *Token) (*User, error)
}

// GitService provides Git operations
type GitService interface {
	GitRepository(ctx context.Context, user *User, repo string) (GitRepository, error)
	FindCommit(ctx context.Context, user *User, repo *Repo) string
}

// ContentService provides information of source codes
type ContentService interface {
	ListAllFiles(ctx context.Context, user *User, repo, ref string) ([]string, error)
	Find(ctx context.Context, user *User, repo, path, ref string) ([]byte, error)
}

// IssueService provides operation to repository issue, basically to pull request
type IssueService interface {
	CreateComment(ctx context.Context, user *User, repo string, number int, body string) (int, error)
	RemoveComment(ctx context.Context, user *User, repo string, number int, id int) error
}
