package git

import (
	"github.com/covergates/covergates/core"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type repository struct {
	gitRepository *git.Repository
}

func (repo *repository) headCommit() string {
	head, err := repo.gitRepository.Head()
	if err != nil {
		return ""
	}
	return head.Hash().String()
}

func (repo *repository) ListAllFiles(commit string) ([]string, error) {
	r := repo.gitRepository
	if commit == "" {
		commit = repo.headCommit()
	}
	commitObject, err := r.CommitObject(plumbing.NewHash(commit))
	if err != nil {
		return nil, err
	}
	tree, err := commitObject.Tree()
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	tree.Files().ForEach(func(f *object.File) error {
		files = append(files, f.Name)
		return nil
	})
	return files, nil
}

func (repo *repository) Commit(commit string) (core.GitCommit, error) {
	return &commitObject{
		repo: repo,
		hash: plumbing.NewHash(commit),
	}, nil
}

func (repo *repository) HeadCommit() string {
	return repo.headCommit()
}

func (repo *repository) Branch() string {
	head, err := repo.gitRepository.Head()
	if err != nil {
		return ""
	}
	if head.Name().IsBranch() {
		return head.Name().Short()
	}
	return ""
}

func (repo *repository) Root() string {
	tree, err := repo.gitRepository.Worktree()
	if err != nil {
		return ""
	}
	return tree.Filesystem.Root()
}
