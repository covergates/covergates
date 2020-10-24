package models

import (
	"fmt"
	"time"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"gorm.io/gorm"
)

// User data
type User struct {
	gorm.Model
	Login  string `gorm:"size:256;uniqueIndex;not null"`
	Name   string
	Email  string `gorm:"index"`
	Active bool
	Avater string
	// Gitea
	GiteaLogin   string `gorm:"index"`
	GiteaEmail   string `gorm:"index"`
	GiteaToken   string
	GiteaRefresh string
	GiteaExpire  int64
	// GitLab
	GitLabLogin   string `gorm:"index"`
	GitLabEmail   string `gorm:"index"`
	GitLabToken   string
	GitLabRefresh string
	GitLabExpire  int64
	// Github
	GithubLogin   string `gorm:"index"`
	GithubEmail   string `gorm:"index"`
	GithubToken   string
	GithubRefresh string
	GithubExpire  int64
	// Bitbucket
	BitbucketLogin   string `gorm:"index"`
	BitbucketEmail   string `gorm:"index"`
	BitbucketToken   string
	BitbucketRefresh string
	BitbucketExpire  int64
	Repositories     []*Repo `gorm:"many2many:user_repositories"`
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

// Update user with new token
func (store *UserStore) Update(scm core.SCMProvider, user *scm.User, token *core.Token) error {
	session := store.DB.Session()
	u, err := store.findWithSCM(scm, user)
	if err != nil {
		return err
	}
	if err := u.updateWithSCM(scm, user, token); err != nil {
		return err
	}
	return session.Save(u).Error
}

// Find user with SCM information
func (store *UserStore) Find(scm core.SCMProvider, user *scm.User) (*core.User, error) {
	u, err := store.findWithSCM(scm, user)
	if err != nil {
		return nil, err
	}
	return u.toCoreUser(), nil
}

func (store *UserStore) findWithSCM(scm core.SCMProvider, user *scm.User) (*User, error) {
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
	case core.GitLab:
		condition = &User{
			GitLabLogin: user.Login,
			GitLabEmail: user.Email,
		}
	case core.Bitbucket:
		condition = &User{
			BitbucketLogin: user.Login,
			BitbucketEmail: user.Email,
		}
	default:
		return nil, &errNotSupportedSCM{scm: scm}
	}
	u := &User{}
	if err := session.Where(condition).First(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByLogin name
func (store *UserStore) FindByLogin(login string) (*core.User, error) {
	session := store.DB.Session()
	u := &User{}
	if err := session.Where(&User{Login: login}).First(u).Error; err != nil {
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

// ListRepositories for the user
func (store *UserStore) ListRepositories(user *core.User) ([]*core.Repo, error) {
	session := store.DB.Session()
	u := &User{}
	if err := session.Preload("Repositories").Where(
		&User{Login: user.Login}).First(u).Error; err != nil {
		return nil, err
	}
	result := make([]*core.Repo, len(u.Repositories))
	for i, repo := range u.Repositories {
		result[i] = repo.ToCoreRepo()
	}
	return result, nil
}

// UpdateRepositories for the user
func (store *UserStore) UpdateRepositories(user *core.User, repositories []*core.Repo) error {
	session := store.DB.Session()
	u := &User{}
	if err := session.Where(&User{Login: user.Login}).First(u).Error; err != nil {
		return err
	}
	userRepos := make([]*Repo, 0, len(repositories))
	for _, repo := range repositories {
		r := &Repo{}
		if err := session.Where(&Repo{URL: repo.URL}).First(r).Error; err != nil {
			return err
		}
		userRepos = append(userRepos, r)
	}
	return session.Model(u).Association("Repositories").Replace(userRepos)
}

func (u *User) toCoreUser() *core.User {
	return &core.User{
		Login:  u.Login,
		Avatar: u.Avater,
		Email:  u.Email,
		// Gitea
		GiteaLogin:   u.GiteaLogin,
		GiteaEmail:   u.GiteaEmail,
		GiteaToken:   u.GiteaToken,
		GiteaExpire:  time.Unix(u.GiteaExpire, 0),
		GiteaRefresh: u.GiteaRefresh,
		// GitLab
		GitLabLogin:   u.GitLabLogin,
		GitLabEmail:   u.GitLabEmail,
		GitLabToken:   u.GitLabToken,
		GitLabExpire:  time.Unix(u.GitLabExpire, 0),
		GitLabRefresh: u.GitLabRefresh,
		// Github
		GithubLogin:   u.GithubLogin,
		GithubEmail:   u.GithubEmail,
		GithubToken:   u.GithubToken,
		GithubExpire:  time.Unix(u.GithubExpire, 0),
		GithubRefresh: u.GithubRefresh,
		// Bitbucket
		BitbucketLogin:   u.BitbucketLogin,
		BitbucketEmail:   u.BitbucketEmail,
		BitbucketToken:   u.BitbucketToken,
		BitbucketExpire:  time.Unix(u.BitbucketExpire, 0),
		BitbucketRefresh: u.BitbucketRefresh,
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
	case core.GitLab:
		u.GitLabEmail = user.Email
		u.GitLabLogin = user.Login
		u.GitLabToken = token.Token
		u.GitLabRefresh = token.Refresh
		u.GitLabExpire = token.Expires.Unix()
	case core.Bitbucket:
		u.BitbucketEmail = user.Email
		u.BitbucketLogin = user.Login
		u.BitbucketToken = token.Token
		u.BitbucketRefresh = token.Refresh
		u.BitbucketExpire = token.Expires.Unix()
	default:
		return &errNotSupportedSCM{scm}
	}
	return nil
}
