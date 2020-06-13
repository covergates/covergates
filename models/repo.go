package models

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/jinzhu/gorm"
)

// Repo defines a repository
type Repo struct {
	gorm.Model
	URL      string `gorm:"unique_index;not null"`
	ReportID string
	Name     string `gorm:"index"`
	SCM      string `gorm:"index"`
}

// RepoStore repositories in storage
type RepoStore struct {
	DB core.DatabaseService
}

// Create a new repository
func (store *RepoStore) Create(scm core.SCMProvider, url, name string) error {
	session := store.DB.Session()
	repo := &Repo{
		URL:  url,
		Name: name,
		SCM:  string(scm),
	}
	return session.Create(repo).Error
}

// Update repository information
func (store *RepoStore) Update(repo *core.Repo) error {
	session := store.DB.Session()
	r := &Repo{}
	if err := session.Where(&Repo{URL: repo.URL}).First(r).Error; err != nil {
		return err
	}
	copyRepo(r, repo)
	return session.Save(r).Error
}

func copyRepo(dst *Repo, src *core.Repo) {
	dst.ReportID = src.ReportID
}
