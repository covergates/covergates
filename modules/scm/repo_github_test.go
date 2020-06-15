// +build github

package scm

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

func getClient() *scm.Client {
	client, _ := github.New("https://api.github.com")
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Scheme: oauth2.SchemeBearer,
			Source: &oauth2.Refresher{
				Source: oauth2.ContextTokenSource(),
			},
		},
	}
	return client
}

func TestGithubList(t *testing.T) {
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	service := &repoService{
		client: getClient(),
		scm:    core.Github,
	}
	repo, err := service.List(context.Background(), user)
	if err != nil {
		t.Error(err)
		return
	}
	if len(repo) <= 0 {
		t.Fail()
	}
}

func TestGithubFind(t *testing.T) {

	service := repoService{
		client: getClient(),
		scm:    core.Github,
	}
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	repo, err := service.Find(context.Background(), user, "blueworrybear/livelogs")
	if err != nil {
		t.Error(err)
		return
	}
	if repo.NameSpace != "blueworrybear" || repo.Name != "livelogs" {
		t.Fail()
	}
	if repo.URL == "" {
		t.Fail()
	}
}
