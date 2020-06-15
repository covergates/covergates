// +build gitea

package scm

import (
	"context"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

func TestGiteaList(t *testing.T) {

	user := &core.User{
		GiteaToken: "1749a6106454f05f689051c331680c13d78d81b7",
	}
	service := repoService{
		client: getClient(),
		scm:    core.Gitea,
	}
	repos, _ := service.List(context.Background(), user)
	if len(repos) < 1 || repos[0].Name != "repo1" || repos[0].NameSpace != "gitea" {
		t.Fail()
	}
}
