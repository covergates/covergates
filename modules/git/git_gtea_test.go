// +build gitea

package git

import (
	"context"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestGitRepoListAllFiles(t *testing.T) {
	s := &Service{}
	ctx := context.Background()
	repo, err := s.Clone(ctx, "http://localhost:3000/gitea/gitea.git", os.Getenv("GITEA_SECRET"))
	if err != nil {
		t.Error(err)
		return
	}
	files, err := repo.ListAllFiles("")
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) <= 0 {
		t.Fail()
	}
}

func TestCheckDefaultBranch(t *testing.T) {
	store := memory.NewStorage()
	repo, _ := git.CloneContext(context.Background(), store, nil, &git.CloneOptions{
		URL: "http://localhost:3000/gitea/JSON.git",
		Auth: &http.BasicAuth{
			Username: "1749a6106454f05f689051c331680c13d78d81b7",
			Password: "x-oauth-basic",
		},
	})
	commit := &commitObject{
		repo: &repository{
			gitRepository: repo,
		},
		hash: plumbing.NewHash("716da5bc95530f1c0400c371ad34002376b10c45"),
	}
	if ok := commit.InDefaultBranch(); ok {
		t.Fail()
	}

	master := &commitObject{
		repo: &repository{
			gitRepository: repo,
		},
		hash: plumbing.NewHash("b08f7b48f44e2df30af7cc538f4056de199338d9"),
	}
	if ok := master.InDefaultBranch(); !ok {
		t.Fail()
	}
}
