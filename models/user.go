package models

import (
	"fmt"
	"time"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/jinzhu/gorm"
)

// User data
type User struct {
	gorm.Model
	Login         string `gorm:"unique_index;not null"`
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

// UserStore user in storage
type UserStore struct {
	DB core.DatabaseService
}

// Create a new user
func (store *UserStore) Create(scm core.SCMProvider, user *scm.User, token *core.Token) error {
	session := store.DB.Session()
	u := &User{
		Login:  user.Login,
		Email:  user.Email,
		Avater: user.Avatar,
		Active: true,
	}
	if err := u.updateWithSCM(scm, user, token); err != nil {
		return err
	}
	return session.Create(u).Error
}

// Find user with SCM information
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

// Bind a new user from another SCM to registered user
func (store *UserStore) Bind(
	scm core.SCMProvider,
	user *core.User,
	scmUser *scm.User,
	token *core.Token,
) (*core.User, error) {
	if user.Login == "" {
		return user, fmt.Errorf("user login should not be empty")
	}
	if _, err := store.Find(scm, scmUser); err == nil {
		return user, errUserExist
	}
	session := store.DB.Session()
	u := &User{}
	if err := session.Where(&User{Login: user.Login}).First(u).Error; err != nil {
		return user, err
	}
	if err := u.updateWithSCM(scm, scmUser, token); err != nil {
		return user, err
	}
	if err := session.Save(u).Error; err != nil {
		return user, err
	}
	return u.toCoreUser(), nil
}

func (u *User) toCoreUser() *core.User {
	return &core.User{
		Login:         u.Login,
		Avatar:        u.Avater,
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

func (u *User) updateWithSCM(scm core.SCMProvider, user *scm.User, token *core.Token) error {
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
	return nil
}
