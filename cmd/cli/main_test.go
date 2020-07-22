package main

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestGit(t *testing.T) {
	repo, err := git.PlainOpenWithOptions("/home/blueworrybear/projects/perl/JSON/lib", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	iter, err := repo.Branches()
	if err != nil {
		t.Error(err)
		return
	}
	iter.ForEach(func(ref *plumbing.Reference) error {
		t.Log(ref.Name())
		return nil
	})
	head, _ := repo.Head()
	t.Log(head.Hash())
	t.Log(head.Name().Short())
	t.Fail()
}
