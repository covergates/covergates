package git

import "github.com/go-git/go-git/v5/plumbing"

type commitObject struct {
	repo *repository
	hash plumbing.Hash
}

func (commit *commitObject) InDefaultBranch() bool {
	head := plumbing.NewHash(commit.repo.headCommit())
	memo := make(map[plumbing.Hash]bool)
	return commit.reaches(head, commit.hash, memo)
}

func (commit *commitObject) reaches(start, target plumbing.Hash, memo map[plumbing.Hash]bool) bool {
	if v, ok := memo[start]; ok {
		return v
	}
	if start == target {
		memo[start] = true
		return true
	}
	repo := commit.repo.gitRepository
	current, err := repo.CommitObject(start)
	if err != nil {
		memo[start] = false
		return false
	}
	for _, parent := range current.ParentHashes {
		ok := commit.reaches(parent, target, memo)
		if ok {
			memo[start] = true
			return true
		}
	}
	memo[start] = false
	return false
}
