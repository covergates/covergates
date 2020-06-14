package models

import (
	"errors"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/jinzhu/gorm"
)

var errEmptyRepoFiled = errors.New("repository must have SCM and URL filed")

// Repo defines a repository
type Repo struct {
	gorm.Model
	URL       string `gorm:"unique_index;not null"`
	ReportID  string
	NameSpace string `gorm:"index;not null`
	Name      string `gorm:"index;not null"`
	SCM       string `gorm:"index;not null"`
}

// RepoStore repositories in storage
type RepoStore struct {
	DB core.DatabaseService
}

func (repo *Repo) ToCoreRepo() *core.Repo {
	return &core.Repo{
		ID:        repo.ID,
		Name:      repo.Name,
		NameSpace: repo.NameSpace,
		ReportID:  repo.ReportID,
		SCM:       core.SCMProvider(repo.SCM),
		URL:       repo.URL,
	}
}

// Create a new repository
func (store *RepoStore) Create(repo *core.Repo) error {
	if repo.SCM == "" || repo.URL == "" {
		return errEmptyRepoFiled
	}
	session := store.DB.Session()
	r := &Repo{
		URL:       repo.URL,
		NameSpace: repo.NameSpace,
		Name:      repo.Name,
		SCM:       string(repo.SCM),
	}
	return session.Create(r).Error
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

func (store *RepoStore) Find(repo *core.Repo) (*core.Repo, error) {
	session := store.DB.Session()
	r := &Repo{}
	if err := session.Where(query(repo)).First(r).Error; err != nil {
		return nil, err
	}
	return r.ToCoreRepo(), nil
}

func (store *RepoStore) Finds(urls ...string) ([]*core.Repo, error) {
	session := store.DB.Session()
	var repositories []*Repo
	session = session.Where("url in (?)", urls).Find(&repositories)
	if err := session.Error; err != nil {
		return nil, err
	}
	coreRepositories := make([]*core.Repo, len(repositories))
	for i, repo := range repositories {
		coreRepositories[i] = repo.ToCoreRepo()
	}
	return coreRepositories, nil
}

func query(repo *core.Repo) *Repo {
	r := &Repo{}
	if repo.ID > 0 {
		r.ID = repo.ID
	}
	if repo.Name != "" {
		r.Name = repo.Name
	}
	if repo.NameSpace != "" {
		r.NameSpace = repo.NameSpace
	}
	if repo.SCM != "" {
		r.SCM = string(repo.SCM)
	}
	if repo.URL != "" {
		r.URL = repo.URL
	}
	return r
}

func copyRepo(dst *Repo, src *core.Repo) {
	dst.ReportID = src.ReportID
}
