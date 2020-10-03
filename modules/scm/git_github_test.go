// +build github

package scm

import (
	"context"
	"os"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/google/go-cmp/cmp"
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

func TestGithubListCommits(t *testing.T) {
	service := &gitService{
		scm:       core.Github,
		scmClient: getGithubClient(),
	}
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	ctx := context.Background()

	commits, err := service.ListCommitsByRef(ctx, user, "octocat/Hello-World", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) <= 0 {
		t.Fatal()
	}
	if diff := cmp.Diff("b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf", commits[0].Sha); diff != "" {
		t.Fatal(diff)
	}
}
