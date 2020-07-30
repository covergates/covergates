// +build github

package scm

import (
	"context"
	"os"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/modules/git"
)

func TestContentGithubListAllFiles(t *testing.T) {
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	service := &contentService{
		client: getGithubClient(),
		scm:    core.Github,
		git:    &git.Service{},
	}
	files, err := service.ListAllFiles(
		context.Background(),
		user,
		"blueworrybear/livelogs", "master")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(files)
	if len(files) <= 20 {
		t.Fail()
	}
}

func TestContentGithubFind(t *testing.T) {
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	service := &contentService{
		client: getGithubClient(),
		scm:    core.Github,
		git:    &git.Service{},
	}
	content, err := service.Find(context.Background(), user, "blueworrybear/livelogs", "go.mod", "master")
	if err != nil {
		t.Error(err)
		return
	}
	if string(content) == "" {
		t.Fail()
	}
}
