package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/covergates/covergates/core"
	"gorm.io/gorm"
)

var errEmptyRepoFiled = errors.New("repository must have SCM and URL filed")

// Repo defines a repository
type Repo struct {
	gorm.Model
	URL       string `gorm:"size:256;uniqueIndex;not null"`
	ReportID  string
	NameSpace string `gorm:"index;not null"`
	Name      string `gorm:"index;not null"`
	Branch    string
	SCM       string `gorm:"index;not null"`
	Creator   string
	Private   bool
}

// RepoSetting defines user customization
type RepoSetting struct {
	gorm.Model
	RepoID uint `gorm:"uniqueIndex"`
	Config []byte
}

// RepoHook defines created hooks in a repository
type RepoHook struct {
	gorm.Model
	RepoID uint `gorm:"index"`
	Hook   string
}

// RepoStore repositories in storage
type RepoStore struct {
	DB core.DatabaseService
}

// ToCoreRepo object
func (repo *Repo) ToCoreRepo() *core.Repo {
	return &core.Repo{
		ID:        repo.ID,
		Name:      repo.Name,
		NameSpace: repo.NameSpace,
		ReportID:  repo.ReportID,
		SCM:       core.SCMProvider(repo.SCM),
		Branch:    repo.Branch,
		URL:       repo.URL,
		Private:   repo.Private,
	}
}

// Update with a new setting
func (setting *RepoSetting) Update(newSetting *core.RepoSetting) error {
	data, err := json.Marshal(newSetting)
	if err != nil {
		return err
	}
	setting.Config = data
	return nil
}

// Create a new repository
func (store *RepoStore) Create(repo *core.Repo, user *core.User) error {
	if repo.SCM == "" || repo.URL == "" {
		return errEmptyRepoFiled
	}
	session := store.DB.Session()
	r := &Repo{
		URL:       repo.URL,
		NameSpace: repo.NameSpace,
		Name:      repo.Name,
		SCM:       string(repo.SCM),
		Branch:    repo.Branch,
		Creator:   user.Login,
		Private:   repo.Private,
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

// Find Repo with seed. The non-empty filed of input will use as where condition
func (store *RepoStore) Find(repo *core.Repo) (*core.Repo, error) {
	session := store.DB.Session()
	r := &Repo{}
	if err := session.Where(repo).First(r).Error; err != nil {
		return nil, err
	}
	return r.ToCoreRepo(), nil
}

// Creator user who activated the repository
func (store *RepoStore) Creator(repo *core.Repo) (*core.User, error) {
	session := store.DB.Session()
	r := &Repo{}
	if err := session.Where(repo).First(r).Error; err != nil {
		return nil, err
	}
	user := &User{
		Login: r.Creator,
	}
	if err := session.First(user, user).Error; err != nil {
		return nil, err
	}
	return user.toCoreUser(), nil
}

// Finds all repositories with URLs
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

// Setting of the repository
func (store *RepoStore) Setting(repo *core.Repo) (*core.RepoSetting, error) {
	session := store.DB.Session()
	setting := &RepoSetting{RepoID: repo.ID}
	if err := session.Where(setting).First(setting).Error; err != nil {
		return &core.RepoSetting{}, nil
	}
	coreSetting := &core.RepoSetting{}
	if err := json.Unmarshal(setting.Config, coreSetting); err != nil {
		return nil, err
	}
	return coreSetting, nil
}

// UpdateSetting for the repository
func (store *RepoStore) UpdateSetting(repo *core.Repo, setting *core.RepoSetting) error {
	session := store.DB.Session()
	repoSetting := &RepoSetting{RepoID: repo.ID}
	if err := session.Where(repoSetting).FirstOrCreate(repoSetting).Error; err != nil {
		return err
	}
	if err := repoSetting.Update(setting); err != nil {
		return err
	}
	return session.Save(repoSetting).Error
}

// UpdateHook of the repository
func (store *RepoStore) UpdateHook(repo *core.Repo, hook *core.Hook) error {
	if repo.ID <= 0 {
		return fmt.Errorf("invalid repository")
	}
	if hook.ID == "" {
		return fmt.Errorf("invalid hook")
	}
	session := store.DB.Session()
	h := &RepoHook{}
	if err := session.FirstOrCreate(h, &RepoHook{RepoID: repo.ID}).Error; err != nil {
		return err
	}
	h.Hook = hook.ID
	return session.Save(h).Error
}

// FindHook for the repository
func (store *RepoStore) FindHook(repo *core.Repo) (*core.Hook, error) {
	session := store.DB.Session()
	h := &RepoHook{}
	if err := session.First(h, &RepoHook{RepoID: repo.ID}).Error; err != nil {
		return nil, err
	}
	return &core.Hook{
		ID: h.Hook,
	}, nil
}

func copyRepo(dst *Repo, src *core.Repo) {
	dst.ReportID = src.ReportID
	dst.Branch = src.Branch
	dst.Private = src.Private
}
