package git

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/google/go-cmp/cmp"
)

func TestRepoRoot(t *testing.T) {
	cwd, _ := os.Getwd()
	dir, err := ioutil.TempDir(cwd, "git_")
	if err != nil {
		t.Fatal(err)
	}
	storagePath := osfs.New(path.Join(dir, ".git"))
	storage := filesystem.NewStorage(storagePath, cache.NewObjectLRUDefault())
	defer func() {
		storage.Close()
		os.RemoveAll(dir)
	}()
	root := osfs.New(dir)
	_, err = git.Init(storage, root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	s := &Service{}
	repo, err := s.PlainOpen(ctx, path.Join("./", path.Base(dir)))
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(dir, repo.Root()); diff != "" {
		t.Fatal(diff)
	}
}
