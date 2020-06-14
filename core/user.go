package core

import (
	"context"
	"time"

	"github.com/drone/go-scm/scm"
)

//go:generate mockgen -package mock -destination ../mock/user_mock.go . UserService,UserStore

// User is a user of the service
type User struct {
	Login         string
	Email         string
	Avater        string
	GiteaLogin    string
	GiteaEmail    string
	GiteaToken    string
	GiteaRefresh  string
	GiteaExpire   time.Time
	GithubLogin   string
	GithubEmail   string
	GithubToken   string
	GithubRefresh string
	GithubExpire  time.Time
}

// UserService defines operations with SCM
type UserService interface {
	Create(ctx context.Context, scm SCMProvider) error
	Find(ctx context.Context, scm SCMProvider) (*User, error)
}

// UserStore the user data to storage
type UserStore interface {
	Create(scm SCMProvider, user *scm.User, token *scm.Token) error
	Find(scm SCMProvider, user *scm.User) (*User, error)
}
