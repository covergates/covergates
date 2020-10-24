package core

import (
	"time"

	"github.com/drone/go-scm/scm"
)

//go:generate mockgen -package mock -destination ../mock/user_mock.go . UserStore

// User is a user of the service
type User struct {
	Login  string
	Email  string
	Avatar string
	// Gitea
	GiteaLogin   string
	GiteaEmail   string
	GiteaToken   string
	GiteaRefresh string
	GiteaExpire  time.Time
	// GitLab
	GitLabLogin   string
	GitLabEmail   string
	GitLabToken   string
	GitLabRefresh string
	GitLabExpire  time.Time
	// Github
	GithubLogin   string
	GithubEmail   string
	GithubToken   string
	GithubRefresh string
	GithubExpire  time.Time
	// Bitbucket
	BitbucketLogin   string
	BitbucketEmail   string
	BitbucketToken   string
	BitbucketRefresh string
	BitbucketExpire  time.Time
}

// UserStore the user data to storage
type UserStore interface {
	Create(scm SCMProvider, user *scm.User, token *Token) error
	Find(scm SCMProvider, user *scm.User) (*User, error)
	FindByLogin(login string) (*User, error)
	Update(scm SCMProvider, user *scm.User, token *Token) error
	// Bind a new user from another SCM to registered user
	Bind(scm SCMProvider, user *User, scmUser *scm.User, token *Token) (*User, error)
	ListRepositories(user *User) ([]*Repo, error)
	UpdateRepositories(user *User, repositories []*Repo) error
}
