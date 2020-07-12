package scm

import (
	"context"
	"sync"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type contentService struct {
	sync.Mutex
	client *scm.Client
	git    core.Git
	scm    core.SCMProvider
}

func (service *contentService) ListAllFiles(
	ctx context.Context,
	user *core.User,
	repo, ref string,
) ([]string, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)
	commit, _, err := client.Git.FindCommit(ctx, repo, ref)
	if err != nil {
		return nil, err
	}
	r, _, err := client.Repositories.Find(ctx, repo)
	if err != nil {
		return nil, err
	}
	token := userToken(service.scm, user)
	gitRepo, err := service.git.Clone(ctx, r.Clone, token.Token)
	if err != nil {
		return nil, err
	}
	return gitRepo.ListAllFiles(commit.Sha)
}

func (service *contentService) Find(ctx context.Context, user *core.User, repo, path, ref string) ([]byte, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)
	content, _, err := client.Contents.Find(ctx, repo, path, ref)
	return content.Data, err
}
