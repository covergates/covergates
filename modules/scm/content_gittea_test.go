// +build gitea

package scm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/git"
)

func TestContentGiteaListAllFilesPerformance(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user := &core.User{
		GiteaToken: os.Getenv("GITEA_SECRET"),
	}
	service := &contentService{
		client: getGiteaClient(),
		scm:    core.Gitea,
		git:    &git.Service{},
	}
	go func() {
		defer cancel()
		files, err := service.ListAllFiles(
			ctx,
			user,
			"gitea/repo1", "master",
		)
		if err != nil {
			t.Error()
			t.Log(err)
		}
		if len(files) <= 0 {
			t.Fail()
		}
	}()
	select {
	case <-ctx.Done():
	case <-time.After(300 * time.Second):
		t.Fail()
		t.Log("timeout")
	}
}
