package scm

import (
	"context"
	"testing"
	"time"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

// FIXME: Change testing repository
func TestSCMClientGithub(t *testing.T) {
	config := &config.Config{
		Github: config.Github{
			Server:    "https://github.com",
			APIServer: "https://api.github.com",
		},
	}
	service := &scmClientService{
		config: config,
	}
	client, err := service.Client(core.Github)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	content, _, err := client.Contents.Find(ctx, "blueworrybear/livelogs", "README.md", "master")
	if err != nil {
		t.Error(err)
		return
	}
	if content.Path != "README.md" {
		t.Fail()
	}
	if string(content.Data) == "" {
		t.Fail()
	}
}

func TestSCMClientGitea(t *testing.T) {
	config := &config.Config{
		Gitea: config.Gitea{
			Server:     "http://localhost:3000",
			SkipVerity: true,
		},
	}
	service := &scmClientService{
		config: config,
	}
	client, err := service.Client(core.Gitea)
	if err != nil {
		t.Error(err)
		return
	}
	expire := time.Now()
	expire.AddDate(0, 0, 1)
	ctx := context.Background()
	ctx = scm.WithContext(ctx, &scm.Token{
		Token:   "",
		Refresh: "",
		Expires: expire,
	})
	content, _, err := client.Contents.Find(ctx, "gitea/repo1", "README.md", "master")
	if err != nil {
		t.Error(err)
		return
	}
	if content.Path != "README.md" {
		t.Fail()
	}
	if string(content.Data) == "" {
		t.Fail()
	}
}
