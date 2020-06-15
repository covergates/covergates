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

// RepoStore repository in storage
type RepoStore interface {
	Create(repo *Repo) error
	Update(repo *Repo) error
	Find(repo *Repo) (*Repo, error)
	Finds(urls ...string) ([]*Repo, error)
}
