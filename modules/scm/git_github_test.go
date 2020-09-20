// +build github

package scm

import (
	"context"
	"os"
	"testing"

	"github.com/covergates/covergates/core"
)

func TestGithubFindCommit(t *testing.T) {
	service := &gitService{
		scm:       core.Github,
		scmClient: getGithubClient(),
	}
	ctx := context.Background()
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	sha := service.FindCommit(ctx, user, &core.Repo{
		Name:      "livelogs",
		NameSpace: "blueworrybear",
		Branch:    "master",
	})
	if sha == "" {
		t.Fail()
	}
}

func TestGitHubListBranch(t *testing.T) {
	service := &gitService{
		scm:       core.Github,
		scmClient: getGithubClient(),
	}
	ctx := context.Background()
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	branches, err := service.ListBranches(ctx, user, "blueworrybear/livelogs")
	if err != nil {
		t.Fatal(err)
	}
	if len(branches) < 1 {
		t.Fatal()
	}
}
