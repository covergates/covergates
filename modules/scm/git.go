package scm

import (
	"context"
	"fmt"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
)

type gitService struct {
	git       core.Git
	scm       core.SCMProvider
	scmClient *scm.Client
}

func (service *gitService) FindCommit(ctx context.Context, user *core.User, repo *core.Repo) string {
	client := service.scmClient
	ctx = withUser(ctx, service.scm, user)
	ref, _, err := client.Git.FindBranch(
		ctx,
		fmt.Sprintf("%s/%s", repo.NameSpace, repo.Name),
		repo.Branch,
	)
	if err != nil {
		return ""
	}
	return ref.Sha
}

// GitRepository clone
func (service *gitService) GitRepository(ctx context.Context, user *core.User, repo string) (core.GitRepository, error) {
	client := service.scmClient
	rs := &repoService{scm: service.scm, client: client}
	token := userToken(service.scm, user)
	ctx = withUser(ctx, service.scm, user)
	url, err := rs.CloneURL(ctx, user, repo)
	if err != nil {
		return nil, err
	}
	return service.git.Clone(ctx, url, token.Token)
}
