package scm

import (
	"context"
	"fmt"

	"github.com/code-devel-cover/CodeCover/core"
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

func (servie *gitService) GitRepository(ctx context.Context, user *core.User, repo string) (core.GitRepository, error) {
	client := servie.scmClient
	rs := &repoService{scm: servie.scm, client: client}
	token := userToken(servie.scm, user)
	ctx = withUser(ctx, servie.scm, user)
	url, err := rs.CloneURL(ctx, user, repo)
	if err != nil {
		return nil, err
	}
	return servie.git.Clone(ctx, url, token.Token)
}
