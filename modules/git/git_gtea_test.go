// +build gitea

package git

import (
	"context"
	"os"
	"testing"
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
