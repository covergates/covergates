package scm

import (
	"context"
	"testing"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
)

// FIXME: Change testing repository
func TestSCMClientGithub(t *testing.T) {
	config := &config.Config{
		Github: config.Github{
			Server:    "https://github.com",
			APIServer: "https://api.github.com",
		},
	}
	service := *&scmClientService{
		config: config,
	}
	client := service.Client(core.Github)
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
