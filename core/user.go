package core

import (
	"time"

	"github.com/drone/go-scm/scm"
)

//go:generate mockgen -package mock -destination ../mock/user_mock.go . UserStore

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

// UserStore the user data to storage
type UserStore interface {
	Create(scm SCMProvider, user *scm.User, token *Token) error
	Find(scm SCMProvider, user *scm.User) (*User, error)
	// Bind a new user from another SCM to registered user
	Bind(scm SCMProvider, user *User, scmUser *scm.User, token *Token) (*User, error)
}
