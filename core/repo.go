package core

import (
	"context"
	"fmt"
)

//go:generate mockgen -package mock -destination ../mock/repo_mock.go . RepoStore,RepoService

// Repo defined a repository structure
type Repo struct {
	ID        uint
	URL       string
	ReportID  string
	NameSpace string
	Name      string
	Branch    string
	Private   bool
	SCM       SCMProvider
}

// RepoSetting to customize repository
type RepoSetting struct {
	Filters          FileNameFilters    `json:"filters"`
	MergePullRequest bool               `json:"mergePR"`
	UpdateAction     ReportUpdateAction `json:"updateAction"`
	// Protected project from unauthorized user upload report
	Protected bool `json:"protected"`
}

// RepoService provides repository opperations
type RepoService interface {
	// Synchronize repositories data from remote and store to database
	Synchronize(ctx context.Context, user *User) error
}

// RepoStore repository in storage
type RepoStore interface {
	Create(repo *Repo) error
	Update(repo *Repo) error
	UpdateOrCreate(repo *Repo) error
	BatchUpdateOrCreate(repos []*Repo) error
	Find(repo *Repo) (*Repo, error)
	Finds(urls ...string) ([]*Repo, error)
	// Creator is the user activated the repository
	Creator(repo *Repo) (*User, error)
	UpdateCreator(repo *Repo, user *User) error
	Setting(repo *Repo) (*RepoSetting, error)
	UpdateSetting(repo *Repo, setting *RepoSetting) error
	FindHook(repo *Repo) (*Hook, error)
	UpdateHook(repo *Repo, hook *Hook) error
}

// FullName is namespace+name
func (repo *Repo) FullName() string {
	return fmt.Sprintf("%s/%s", repo.NameSpace, repo.Name)
}
