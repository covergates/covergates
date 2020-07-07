// +build github

package scm

import (
	"context"
	"os"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

func TestContentGithubListAllFiles(t *testing.T) {
	user := &core.User{
		GithubToken: os.Getenv("GITHUB_SECRET"),
	}
	service := &contentService{
		client: getClient(),
		scm:    core.Github,
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
		client: getClient(),
		scm:    core.Github,
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
