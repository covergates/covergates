package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Login         string
	Email         string
	Active        bool
	Avater        string
	GiteaToken    string
	GiteaRefresh  string
	GiteaExpire   int64
	GithubToken   string
	GithubRefresh string
	GithubExpre   int64
	Repositories  []*Repo
}
