package core

//go:generate mockgen -package mock -destination ../mock/repo_mock.go . RepoStore

// Repo defined a repository structure
type Repo struct {
	ID        uint
	URL       string
	ReportID  string
	NameSpace string
	Name      string
	Branch    string
	SCM       SCMProvider
}

// RepoSetting to customize repository
type RepoSetting struct {
	Filters FileNameFilters
}

// RepoStore repository in storage
type RepoStore interface {
	Create(repo *Repo, creator *User) error
	Update(repo *Repo) error
	Find(repo *Repo) (*Repo, error)
	Finds(urls ...string) ([]*Repo, error)
	Creator(repo *Repo) (*User, error)
	Setting(repo *Repo) (*RepoSetting, error)
	UpdateSetting(repo *Repo, setting *RepoSetting) error
}
