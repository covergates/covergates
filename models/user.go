package models

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Login         string `gorm:"unique_index"`
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

// TODO: Return core.User
func (store *UserStore) Find(scm core.SCMProvider, user *scm.User) (*core.User, error) {
	session := store.DB.Session()
	u := &User{}
	session.First(u)
	return nil, nil
}
