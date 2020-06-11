package core

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type User struct {
	Login         string
	Email         string
	Avater        string
	GiteaLogin    string
	GiteaEmail    string
	GiteaToken    string
	GiteaRefresh  string
	GiteaExpire   int64
	GithubLogin   string
	GithubEmail   string
	GithubToken   string
	GithubRefresh string
	GithubExpire  int64
}

type UserService interface {
	Find(ctx context.Context, scm SCMProvider) (*User, error)
}

type UserStore interface {
	Find(scm SCMProvider, user *scm.User) (*User, error)
}
