package models

import (
	"time"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Login         string `gorm:"unique_index"`
	Name          string
	Email         string `gorm:"unique_index"`
	Active        bool
	Avater        string
	GiteaLogin    string `gorm:"index"`
	GiteaEmail    string `gorm:"index"`
	GiteaToken    string
	GiteaRefresh  string
	GiteaExpire   int64
	GithubLogin   string `gorm:"index"`
	GithubEmail   string `gorm:"index"`
	GithubToken   string
	GithubRefresh string
	GithubExpire  int64
}

type UserStore struct {
	DB core.DatabaseService
}

func (store *UserStore) Create(scm core.SCMProvider, user *scm.User, token *scm.Token) error {
	session := store.DB.Session()
	u := &User{
		Login:  user.Login,
		Email:  user.Email,
		Avater: user.Avatar,
		Active: true,
	}
	switch scm {
	case core.Github:
		u.GithubEmail = user.Email
		u.GithubLogin = user.Login
		u.GithubToken = token.Token
		u.GithubRefresh = token.Refresh
		u.GithubExpire = token.Expires.Unix()
	case core.Gitea:
		u.GiteaEmail = user.Email
		u.GiteaLogin = user.Login
		u.GiteaToken = token.Token
		u.GiteaRefresh = token.Refresh
		u.GiteaExpire = token.Expires.Unix()
	default:
		return &errNotSupportedSCM{scm}
	}
	return session.Create(u).Error
}

func (store *UserStore) Find(scm core.SCMProvider, user *scm.User) (*core.User, error) {
	session := store.DB.Session()
	var condition *User
	switch scm {
	case core.Github:
		condition = &User{
			GithubLogin: user.Login,
			GithubEmail: user.Email,
		}
	case core.Gitea:
		condition = &User{
			GiteaLogin: user.Login,
			GiteaEmail: user.Email,
		}
	default:
		return nil, &errNotSupportedSCM{scm: scm}
	}
	u := &User{}
	if err := session.Where(condition).First(u).Error; err != nil {
		return nil, err
	}
	return u.toCoreUser(), nil
}

func (u *User) toCoreUser() *core.User {
	return &core.User{
		Login:         u.Login,
		Avater:        u.Avater,
		Email:         u.Email,
		GiteaLogin:    u.GiteaLogin,
		GiteaEmail:    u.GiteaEmail,
		GiteaToken:    u.GiteaToken,
		GiteaExpire:   time.Unix(u.GiteaExpire, 0),
		GiteaRefresh:  u.GiteaRefresh,
		GithubLogin:   u.GithubLogin,
		GithubEmail:   u.GithubEmail,
		GithubToken:   u.GithubToken,
		GithubExpire:  time.Unix(u.GithubExpire, 0),
		GithubRefresh: u.GithubRefresh,
	}
}
